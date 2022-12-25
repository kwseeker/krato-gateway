package main

import (
	"context"
	"flag"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kwseeker/kratos-gateway/gateway-02/app"
	"github.com/kwseeker/kratos-gateway/gateway-02/server"
	"github.com/kwseeker/kratos-gateway/gateway-02/server/handler"
	"net/http"
	_ "net/http/pprof"
	"time"
)

var conf string

func init() {
	flag.StringVar(&conf, "conf", "config.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	//http://localhost:17070/debug/pprof/
	adminAddr := "0.0.0.0:17070"
	bindAddr := ":18080"
	timeout := 5 * time.Second
	idleTimeout := 30 * time.Second

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
		log.Fatalf("failed to load config: %v", err)
	}

	p := new(handler.Echo)

	ctx := context.Background()
	srv := server.New(*p, bindAddr, timeout, idleTimeout)

	a := app.New(
		app.Name("srv1"),
		app.Context(ctx),
		app.Server(srv),
	)
	if err := a.Run(); err != nil {
		log.Errorf("failed to run servers: %v", err)
	}
}
