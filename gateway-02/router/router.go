package router

import (
	"net/http"
)

// Router is 网关路由器，其中 routes 字段保存路由Map
type Router interface {
	http.Handler
	Handle(pattern, method string, handler http.Handler) error
}
