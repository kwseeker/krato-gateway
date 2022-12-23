package app

import (
	"context"
	"github.com/kwseeker/kratos-gateway/gateway-02/app/transport"
	"os"
)

type Option func(o *options)

type options struct {
	id      string             //当前网关ID
	name    string             //当前网关名称
	ctx     context.Context    //继承的上下文
	sigs    []os.Signal        //监听的关闭信号
	servers []transport.Server //网关服务Server(HTTP\TCP等)
}

// setter

func ID(id string) Option {
	return func(o *options) { o.id = id }
}

func Name(name string) Option {
	return func(o *options) { o.name = name }
}

func Context(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func Server(srv ...transport.Server) Option {
	return func(o *options) {
		o.servers = srv
	}
}
