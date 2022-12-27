package proxy

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	config "github.com/kwseeker/kratos-gateway/gateway-02/api/gateway/config/v1"
	"github.com/kwseeker/kratos-gateway/gateway-02/client"
	"github.com/kwseeker/kratos-gateway/gateway-02/middleware"
	"github.com/kwseeker/kratos-gateway/gateway-02/router"
	"github.com/kwseeker/kratos-gateway/gateway-02/router/mux"
	"io"
	"net"
	"net/http"
	"runtime"
	"sync/atomic"
)

const xff = "X-Forwarded-For"

// Proxy is a gateway proxy.
type Proxy struct {
	router            atomic.Value
	log               *log.Helper
	clientFactory     client.Factory
	middlewareFactory middleware.Factory
}

// New is new a gateway proxy.
func New(logger log.Logger, clientFactory client.Factory, middlewareFactory middleware.Factory) (*Proxy, error) {
	p := &Proxy{
		log:               log.NewHelper(logger),
		clientFactory:     clientFactory,
		middlewareFactory: middlewareFactory,
	}
	p.router.Store(mux.NewRouter())
	return p, nil
}

func (p *Proxy) buildMiddleware(ms []*config.Middleware, handler middleware.Handler) (middleware.Handler, error) {
	for _, c := range ms {
		m, err := p.middlewareFactory(c)
		if err != nil {
			return nil, err
		}
		handler = m(handler)
	}
	return handler, nil
}

func (p *Proxy) buildEndpoint(e *config.Endpoint, ms []*config.Middleware) (http.Handler, error) {
	caller, err := p.clientFactory(e)
	if err != nil {
		return nil, err
	}
	handler, err := p.buildMiddleware(ms, caller.Do)
	if err != nil {
		return nil, err
	}
	handler, err = p.buildMiddleware(e.Middlewares, handler)
	if err != nil {
		return nil, err
	}
	return http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {
			r.Header[xff] = append(r.Header[xff], ip)
		}
		ctx := middleware.NewRequestContext(r.Context(), &middleware.RequestOptions{})
		ctx, cancel := context.WithTimeout(ctx, e.Timeout.AsDuration())
		defer cancel()
		resp, err := handler(ctx, r)
		if err != nil {
			switch err {
			case context.Canceled:
				w.WriteHeader(499)
			case context.DeadlineExceeded:
				w.WriteHeader(504)
			default:
				w.WriteHeader(502)
			}
			return
		}
		headers := w.Header()
		for k, v := range resp.Header {
			headers[k] = v
		}
		w.WriteHeader(resp.StatusCode)
		if body := resp.Body; body != nil {
			_, _ = io.Copy(w, body)
		}
		// see https://pkg.go.dev/net/http#example-ResponseWriter-Trailers
		for k, v := range resp.Trailer {
			headers[http.TrailerPrefix+k] = v
		}
		resp.Body.Close()
	})), nil
}

// Update 刷新 Endpoints, 主要是创建客户端连接、然后装饰middleware组件、然后再封装成http.Handler方法，
// 最后将路由信息和http.Handler注册到mux路由
func (p *Proxy) Update(c *config.Gateway) error {
	router := mux.NewRouter()

	for _, e := range c.Endpoints {
		//创建客户端连接、然后装饰middleware组件、然后再封装成http.Handler方法
		handler, err := p.buildEndpoint(e, c.Middlewares)
		if err != nil {
			return err
		}
		//将路由信息和http.Handler注册到mux路由
		if err = router.Handle(e.Path, e.Method, handler); err != nil {
			return err
		}
		p.log.Infof("build endpoint: [%s] %s %s", e.Protocol, e.Method, e.Path)
	}
	p.router.Store(router)
	return nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			buf := make([]byte, 64<<10) //nolint:gomnd
			n := runtime.Stack(buf, false)
			p.log.Errorf("panic recovered: %s", buf[:n])
		}
	}()
	log.Debug("received request: ", req)
	p.router.Load().(router.Router).ServeHTTP(w, req)
}
