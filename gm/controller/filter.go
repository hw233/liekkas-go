package controller

import (
	"gm/token"
	"net/http"
	"shared/utility/errors"
)

type HttpRequestFilter interface {
	filter(r *http.Request) error
}

type BearerAuthenticationFilter struct {
}

func (b *BearerAuthenticationFilter) filter(r *http.Request) error {
	Authentication := r.Header.Get("Authorization")
	if Authentication != token.Bearer {
		return errors.New("Bearer Authentication")
	}
	return nil
}
