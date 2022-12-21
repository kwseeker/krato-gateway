package server

import (
	"context"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net"
	"net/http"
	"time"
)

// Server gateway server 继承 http.Server
type Server struct {
	*http.Server
}

func New(handler http.Handler, addr string, timeout time.Duration, idleTimeout time.Duration) *Server {
	srv := &Server{
		Server: &http.Server{
			Addr: addr, //host:port
			Handler: h2c.NewHandler(handler, &http2.Server{
				IdleTimeout: idleTimeout,
			}), //路由处理, 套娃了个http2.Server处理http的请求
			//TLSConfig:         nil,
			ReadTimeout:       timeout,     //读请求超时时间
			ReadHeaderTimeout: timeout,     //读请求头超时时间
			WriteTimeout:      timeout,     //写响应超时时间
			IdleTimeout:       idleTimeout, //keep-alive 空闲超时时间
			//MaxHeaderBytes:    0,
			//TLSNextProto:      nil,
			//ConnState:         nil,
			//ErrorLog:    nil,
			//BaseContext: nil,				//???
			//ConnContext: nil,
		},
	}
	return srv
}

// Start the server.
func (s *Server) Start(ctx context.Context) error {
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}
	return s.ListenAndServe()
}

// Stop the server.
func (s *Server) Stop(ctx context.Context) error {
	return s.Shutdown(ctx)
}
