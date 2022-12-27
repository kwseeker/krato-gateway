package mux

import (
	"github.com/gorilla/mux"
	"github.com/kwseeker/kratos-gateway/gateway-02/router"
	"net/http"
	"strings"
)

var _ = new(router.Router)

type muxRouter struct {
	*mux.Router
}

// NewRouter new a mux router.
func NewRouter() router.Router {
	return &muxRouter{
		Router: mux.NewRouter().StrictSlash(true),
	}
}

func (r *muxRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}

// Handle 注册一对新的路由和Handler,
// NewRoute()中实现了 r.routes = append(r.routes, route)
func (r *muxRouter) Handle(pattern, method string, handler http.Handler) error {
	next := r.Router.NewRoute().Handler(handler)
	if strings.HasSuffix(pattern, "*") {
		// /api/echo/*
		next = next.PathPrefix(strings.TrimRight(pattern, "*"))
	} else {
		// /api/echo/hello
		// /api/echo/[a-z]+
		// /api/echo/{name}
		next = next.Path(pattern)
	}
	if method != "" && method != "*" {
		next = next.Methods(method)
	}
	return next.GetError()
}
