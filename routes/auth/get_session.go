package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"dash/plugins"
	"dash/redis"

	"github.com/danielgtaylor/huma/v2"
)

type GetSessionInput struct {
	Session http.Cookie `cookie:"session"`
}

type GetSessionOutput struct {
	Body *plugins.Session
}

func GetSession(
	ctx context.Context,
	input *GetSessionInput,
) (*GetSessionOutput, error) {
	if input.Session.Name == "" {
		return &GetSessionOutput{}, nil
	}

	raw, err := redis.Client.Get(ctx, "session:"+input.Session.Value).Result()
	if err != nil {
		return &GetSessionOutput{}, nil
	}

	var session plugins.Session
	if err := json.Unmarshal([]byte(raw), &session); err != nil {
		return nil, huma.Error500InternalServerError("internal error")
	}

	return &GetSessionOutput{
		Body: &session,
	}, nil
}
