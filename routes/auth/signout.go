package auth

import (
	"context"
	"net/http"

	"dash/redis"

	"github.com/danielgtaylor/huma/v2"
)

type SignOutInput struct {
	Session http.Cookie `cookie:"session"`
}

type SignOutOutput struct {
	SetCookie http.Cookie `header:"Set-Cookie"`
}

func SignOut(ctx context.Context, input *SignOutInput) (*SignOutOutput, error) {
	res := &SignOutOutput{
		SetCookie: http.Cookie{
			Name:     "session",
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteLaxMode,
			MaxAge:   -1,
			Path:     "/",
		},
	}

	if input.Session.Name == "" {
		if err := redis.Client.Unlink(ctx, "session:"+input.Session.Value).
			Err(); err != nil {
			return nil, huma.Error500InternalServerError("internal error")
		}
	}

	return res, nil
}
