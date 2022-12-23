package main

import (
	"context"
	"fmt"
	"github.com/kwseeker/kratos-gateway/gateway-01/app"
	"github.com/kwseeker/kratos-gateway/gateway-01/server"
	"github.com/kwseeker/kratos-gateway/gateway-01/server/handler"
	"time"
)

func main() {
	//adminAddr := "0.0.0.0:7072"
	adminAddr := ":7072"
	timeout := 5 * time.Second
	idleTimeout := 30 * time.Second

	//go func() {
	//	_ = http.ListenAndServe(adminAddr, nil)
	//}()

	p := new(handler.Echo)

	ctx := context.Background()
	srv := server.New(*p, adminAddr, timeout, idleTimeout)

	a := app.New(
		app.Name("srv1"),
		app.Context(ctx),
		app.Server(srv),
	)
	if err := a.Run(); err != nil {
		_ = fmt.Errorf("failed to run servers: %v", err)
	}
}
