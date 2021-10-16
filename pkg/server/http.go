package server

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Params = httprouter.Params

type Router interface {
	HandlerFunc(method, path string, handler http.HandlerFunc)
}

func ParamsFromContext(ctx context.Context) Params {
	return httprouter.ParamsFromContext(ctx)
}
