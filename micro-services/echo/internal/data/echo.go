package data

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/kwseeker/kratos-gateway/echo/internal/biz"
)

type echoRepo struct {
	data *Data
	log  *log.Helper
}

// NewEchoRepo .
func NewEchoRepo(data *Data, logger log.Logger) biz.EchoRepo {
	return &echoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *echoRepo) Save(ctx context.Context, g *biz.Echo) (*biz.Echo, error) {
	return g, nil
}

func (r *echoRepo) Update(ctx context.Context, g *biz.Echo) (*biz.Echo, error) {
	return g, nil
}

func (r *echoRepo) FindByID(context.Context, int64) (*biz.Echo, error) {
	return nil, nil
}

func (r *echoRepo) ListByHello(context.Context, string) ([]*biz.Echo, error) {
	return nil, nil
}

func (r *echoRepo) ListAll(context.Context) ([]*biz.Echo, error) {
	return nil, nil
}
