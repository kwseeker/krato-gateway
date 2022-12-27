package main

import (
	"context"
	"flag"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/hashicorp/consul/api"
	configv1 "github.com/kwseeker/kratos-gateway/gateway-02/api/gateway/config/v1"
	"github.com/kwseeker/kratos-gateway/gateway-02/app"
	"github.com/kwseeker/kratos-gateway/gateway-02/client"
	"github.com/kwseeker/kratos-gateway/gateway-02/middleware"
	"github.com/kwseeker/kratos-gateway/gateway-02/proxy"
	"github.com/kwseeker/kratos-gateway/gateway-02/server"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

var (
	conf string
	//consul
	consulAddress    string
	consulToken      string
	consulDatacenter string
)

func init() {
	flag.StringVar(&conf, "conf", "config.yaml", "config path, eg: -conf config.yaml")
	//consul
	flag.StringVar(&consulAddress, "consul.address", "127.0.0.1:8500", "consul address, eg: 127.0.0.1:8500")
	flag.StringVar(&consulToken, "consul.token", "", "consul token, eg: xxx")
	flag.StringVar(&consulDatacenter, "consul.datacenter", "", "consul datacenter, eg: dc1")
}

// consul http客户端
func registry() *consul.Registry {
	if consulAddress != "" {
		c := api.DefaultConfig()
		c.Address = consulAddress
		c.Token = consulToken
		c.Datacenter = consulDatacenter
		//核心是创建http.Client{}实例，即http客户端
		client, err := api.NewClient(c)
		if err != nil {
			panic(err)
		}
		return consul.New(client)
	}
	return nil
}

/*
-conf config.yaml
*/
func main() {
	//http://localhost:17070/debug/pprof/
	adminAddr := "0.0.0.0:17070"
	bindAddr := ":18080"
	timeout := 5 * time.Second
	idleTimeout := 30 * time.Second

	logger := log.NewStdLogger(os.Stdout)
	l := log.NewHelper(logger)

	//pprof
	go func() {
		_ = http.ListenAndServe(adminAddr, nil)
	}()

	//配置
	//1)指定配置源
	c := config.New(
		config.WithSource(
			file.NewSource(conf),
		),
	)
	//2)加载配置
	if err := c.Load(); err != nil {
		l.Fatalf("failed to load config: %v", err)
	}
	//3)解析到配置类实例
	bc := new(configv1.Gateway)
	if err := c.Scan(bc); err != nil {
		log.Fatalf("failed to scan config: %v", err)
	}

	//从注册中心使用负载均衡策略选择服务节点，创建连接到微服务节点的客户端工厂（工厂函数）
	//gateway请求 -> 负载均衡 -> 服务节点
	//registry()创建连接consul的客户端，clientFactory创建连接到后台微服务的客户端
	clientFactory := client.NewFactory(logger, registry())
	//客户端连接代理
	p, err := proxy.New(logger, clientFactory, middleware.Create) //这里中间件是指额外装饰的模块不是单独的服务
	if err != nil {
		l.Fatalf("failed to new proxy: %v", err)
	}
	//p := new(handler.Echo)

	//更新配置, 如果有配置 Endpoints，创建连接后台微服务的客户端caller
	if err := p.Update(bc); err != nil {
		log.Fatalf("failed to update service config: %v", err)
	}
	//监听配置变更

	//启动Http Server, 接收外来HTTP请求
	ctx := context.Background()
	srv := server.New(p, bindAddr, timeout, idleTimeout)
	a := app.New(
		app.Name("srv1"),
		app.Context(ctx),
		app.Server(srv),
	)
	if err := a.Run(); err != nil {
		l.Errorf("failed to run servers: %v", err)
	}
}
