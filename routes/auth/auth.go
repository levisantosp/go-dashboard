package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"dash/db"
	"dash/ent/generated"
	"dash/redis"
	"dash/utils"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
)

type AuthOutput struct {
	Status    int           `json:"-"`
	Location  string        `         header:"Location"`
	SetCookie []http.Cookie `         header:"Set-Cookie"`
}

func Auth(ctx context.Context, _ *struct{}) (*AuthOutput, error) {
	state := uuid.New()

	params := url.Values{}
	params.Set("client_id", utils.Env.DiscordClientID)
	params.Set("redirect_uri", utils.Env.RedirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "identify")
	params.Set("state", state.String())

	return &AuthOutput{
		Status:   http.StatusFound,
		Location: utils.Env.DiscordURL + "/oauth2/authorize?" + params.Encode(),
		SetCookie: []http.Cookie{{
			Name:     "oauth_state",
			Value:    state.String(),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   60 * 5,
			Path:     "/",
		}},
	}, nil
}

type CallbackInput struct {
	Code       string      `query:"code"  required:"true"`
	State      string      `query:"state" required:"true"`
	OAuthState http.Cookie `                              cookie:"oauth_state"`
}

type CallbackOutput struct {
	Status    int           `json:"-"`
	Location  string        `         header:"Location"`
	SetCookie []http.Cookie `         header:"Set-Cookie"`
}

type DiscordTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type DiscordUserResponse struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	GlobalName string `json:"global_name"`
	Avatar     string `json:"avatar"`
}

type SessionPayload struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	GlobalName string `json:"globalName"`
	Avatar     string `json:"avatar"`
	IsAdmin    bool   `json:"isAdmin"`
}

func Callback(
	ctx context.Context,
	input *CallbackInput,
) (*CallbackOutput, error) {
	if input.OAuthState.Value != input.State {
		return nil, huma.Error400BadRequest("invalid state")
	}

	cookies := []http.Cookie{
		{
			Name:     "oauth_state",
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
			Path:     "/",
		},
	}

	params := url.Values{}
	params.Set("client_id", utils.Env.DiscordClientID)
	params.Set("client_secret", utils.Env.DiscordClientSecret)
	params.Set("redirect_uri", utils.Env.RedirectURI)
	params.Set("code", input.Code)
	params.Set("grant_type", "authorization_code")

	tokenReq, _ := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		utils.Env.DiscordURL+"/api/oauth2/token",
		strings.NewReader(params.Encode()),
	)
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenRes, err := http.DefaultClient.Do(tokenReq)
	if err != nil {
		return nil, huma.Error400BadRequest("discord error")
	}
	defer tokenRes.Body.Close()

	var tokenData DiscordTokenResponse
	json.NewDecoder(tokenRes.Body).Decode(&tokenData)

	userReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		utils.Env.DiscordURL+"/api/users/@me",
		nil,
	)
	if err != nil {
		return nil, huma.Error400BadRequest("discord error")
	}

	userReq.Header.Set("Authorization", "Bearer "+tokenData.AccessToken)

	userResp, err := http.DefaultClient.Do(userReq)
	if err != nil {
		return nil, huma.Error400BadRequest("discord error")
	}

	defer userResp.Body.Close()

	var discordUser DiscordUserResponse
	json.NewDecoder(userResp.Body).Decode(&discordUser)

	userId, err := strconv.ParseInt(discordUser.ID, 10, 64)
	if err != nil {
		return nil, huma.Error500InternalServerError("internal error")
	}

	user, err := db.Client.User.Get(ctx, userId)
	if generated.IsNotFound(err) {
		return nil, huma.Error401Unauthorized("unauthorized")
	}
	if err != nil {
		return nil, huma.Error500InternalServerError("internal error")
	}

	sessionId := make([]byte, 32)
	rand.Read(sessionId)

	payload, err := json.Marshal(SessionPayload{
		ID:         discordUser.ID,
		Username:   discordUser.Username,
		GlobalName: discordUser.GlobalName,
		Avatar:     discordUser.Avatar,
		IsAdmin:    user.IsAdmin,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("internal error")
	}

	ttl := 30 * 24 * time.Hour
	redis.Client.Set(
		ctx,
		"session:"+hex.EncodeToString(sessionId),
		payload,
		ttl,
	)

	sessionCookie := http.Cookie{
		Name:     "session",
		Value:    hex.EncodeToString(sessionId),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   int(ttl.Seconds()),
		Path:     "/",
	}

	return &CallbackOutput{
		Status:    http.StatusFound,
		Location:  utils.Env.DashboardURL,
		SetCookie: append(cookies, sessionCookie),
	}, nil
}
