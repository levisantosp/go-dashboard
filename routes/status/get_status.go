package status

import (
	"context"

	"dash/db"
	"dash/ent/generated"

	"github.com/danielgtaylor/huma/v2"
)

type GetStatusInput struct {
	ID int64 `path:"id"`
}

type StatusOutputBody struct {
	ID   int64  `json:"id"`
	Type string `json:"type"`
}

type StatusOutput struct {
	Body StatusOutputBody
}

func Get(ctx context.Context, input *GetStatusInput) (*StatusOutput, error) {
	status, err := db.Client.Status.Get(ctx, input.ID)
	if err != nil {
		if generated.IsNotFound(err) {
			return nil, huma.Error404NotFound("not found")
		}

		return nil, huma.Error500InternalServerError("internal error")
	}

	res := StatusOutput{
		Body: StatusOutputBody{
			ID:   status.ID,
			Type: string(status.Type),
		},
	}

	return &res, nil
}
