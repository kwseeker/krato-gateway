package biz

import (
	"context"

	v1 "github.com/kwseeker/kratos-gateway/echo/api/echo/v1"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

type Echo struct {
	Hello string
}

type EchoRepo interface {
	Save(context.Context, *Echo) (*Echo, error)
	Update(context.Context, *Echo) (*Echo, error)
	FindByID(context.Context, int64) (*Echo, error)
	ListByHello(context.Context, string) ([]*Echo, error)
	ListAll(context.Context) ([]*Echo, error)
}

type EchoUsecase struct {
	repo EchoRepo
	log  *log.Helper
}

func NewEchoUsecase(repo EchoRepo, logger log.Logger) *EchoUsecase {
	return &EchoUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *EchoUsecase) CreateEcho(ctx context.Context, g *Echo) (*Echo, error) {
	uc.log.WithContext(ctx).Infof("CreateEcho: %v", g.Hello)
	return uc.repo.Save(ctx, g)
}
