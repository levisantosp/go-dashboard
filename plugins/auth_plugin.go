package plugins

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"dash/db"
	"dash/ent/generated"
	"dash/redis"

	"github.com/danielgtaylor/huma/v2"
)

type Session struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	GlobalName string `json:"globalName"`
	Avatar     string `json:"avatar"`
	IsAdmin    bool   `json:"isAdmin"`
}

type ctxKey int

const (
	userCtxKey ctxKey = iota
	sessionCtxKey
)

func AuthPlugin(api huma.API) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		cookie, err := huma.ReadCookie(ctx, "session")
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "unauthorized")
			return
		}

		raw, err := redis.Client.Get(ctx.Context(), "session:"+cookie.Value).
			Result()
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "unauthorized")
			return
		}

		var session Session
		if err := json.Unmarshal([]byte(raw), &session); err != nil {
			huma.WriteErr(
				api,
				ctx,
				http.StatusInternalServerError,
				"internal error",
			)
			return
		}

		userId, err := strconv.ParseInt(session.ID, 10, 64)
		if err != nil {
			huma.WriteErr(
				api,
				ctx,
				http.StatusInternalServerError,
				"internal error",
			)
			return
		}

		user, err := db.Client.User.Get(ctx.Context(), userId)
		if err != nil {
			huma.WriteErr(api, ctx, http.StatusUnauthorized, "unauthorized")
			return
		}

		if !user.IsAdmin {
			huma.WriteErr(api, ctx, http.StatusForbidden, "forbidden")
			return
		}

		ctx = huma.WithValue(ctx, userCtxKey, user)
		ctx = huma.WithValue(ctx, sessionCtxKey, session)

		next(ctx)
	}
}

func GetUserFromContext(ctx context.Context) (*generated.User, bool) {
	user, ok := ctx.Value(userCtxKey).(*generated.User)
	return user, ok
}

func SessionFromContext(ctx context.Context) (Session, bool) {
	session, ok := ctx.Value(sessionCtxKey).(Session)
	return session, ok
}
