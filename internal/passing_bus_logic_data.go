package internal

import (
	"context"
	"net/http"
)

type Logic interface {
	BusinessLogic(ctx context.Context, authToken string) (string, error)
}

type Controller struct {
	Logic Logic
}

func (c Controller) DoLogic(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authToken, ok := authFromContext(ctx)

	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	result, err := c.Logic.BusinessLogic(ctx, authToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(result))
}
