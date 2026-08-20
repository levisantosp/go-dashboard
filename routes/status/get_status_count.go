package status

import (
	"context"

	"dash/db"

	"github.com/danielgtaylor/huma/v2"
)

type CountStatusOutput struct {
	Body struct {
		Count int `json:"count"`
	}
}

func Count(ctx context.Context, _ *struct{}) (*CountStatusOutput, error) {
	count, err := db.Client.Status.Query().Count(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("internal error")
	}

	res := CountStatusOutput{}
	res.Body.Count = count

	return &res, nil
}
