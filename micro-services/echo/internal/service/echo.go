package service

import (
	"context"

	v1 "github.com/kwseeker/kratos-gateway/echo/api/echo/v1"
	"github.com/kwseeker/kratos-gateway/echo/internal/biz"
)

// EchoService is a greeter service.
type EchoService struct {
	v1.UnimplementedEchoServer

	uc *biz.EchoUsecase
}

// NewEchoService new a greeter service.
func NewEchoService(uc *biz.EchoUsecase) *EchoService {
	return &EchoService{uc: uc}
}

// SayHello implements echo.GreeterServer.
func (s *EchoService) SayHello(ctx context.Context, in *v1.EchoRequest) (*v1.EchoReply, error) {
	g, err := s.uc.CreateEcho(ctx, &biz.Echo{Hello: in.Name})
	if err != nil {
		return nil, err
	}
	return &v1.EchoReply{Message: "Hello " + g.Hello}, nil
}
