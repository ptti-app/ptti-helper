package helper

import (
	"context"
	"errors"

	"github.com/go-chi/chi/v5"
)

func URLParams(ctx context.Context) (map[string]string, error) {
	params := chi.RouteContext(ctx).URLParams

	mapper := map[string]string{}
	for i := range params.Keys {
		if params.Keys[i] != "" && params.Keys[i] != "*" {
			mapper[params.Keys[i]] = params.Values[i]
		}
	}

	if len(mapper) == 0 {
		return nil, errors.New("param is nil")
	}

	return mapper, nil
}
