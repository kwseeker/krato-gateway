package app

import (
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type Info interface {
	ID() string
	Name() string
}

type App struct {
	opts   options //网关ID、名称、继承的上下文、监听的关闭信号数组、服务Server数组
	ctx    context.Context
	cancel func()
}

func New(opts ...Option) *App {
	o := options{
		ctx:  context.Background(),
		sigs: []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}, //支持接收关闭信号2、3、15
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}

	// opts 是Setter函数数组
	for _, opt := range opts {
		opt(&o)
	}

	ctx, cancelFunc := context.WithCancel(o.ctx)
	return &App{
		opts:   o,
		ctx:    ctx,
		cancel: cancelFunc,
	}
}

func (a *App) ID() string {
	return a.opts.id
}

func (a *App) Name() string {
	return a.opts.name
}

// Run 启动所有options.servers, 然后注册关闭钩子
func (a *App) Run() error {
	ctx := NewContext(a.ctx, a)

	eg, ctx := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	for _, srv := range a.opts.servers {
		srv := srv
		//注册关闭钩子, 应用上下文关闭后关闭所有服务节点
		eg.Go(func() error {
			<-ctx.Done()
			return srv.Stop(ctx)
		})
		//用WaitGroup等待所有服务节点启动完毕
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(ctx)
		})
	}
	wg.Wait()

	//监听关闭信号
	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		//for { //官方源码这里搞个for循环啥用？ <-ctx.Done 和 <-c 都是阻塞的, 且case中没有异步操作，应该不会循环执行
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c:
			err := a.Stop() //收到关闭信号，关闭所有servers
			if err != nil {
				log.Errorf("failed to stop app: %v", err)
				return err
			}
			return nil
		}
		//}
	})

	//阻塞等待关闭
	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		log.Errorf("gateway app exit! err: %v", err)
		return err
	}
	log.Warn("gateway app exit!")
	return nil
}

func (a *App) Stop() error {
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

type appKey struct{}

func NewContext(ctx context.Context, info Info) context.Context {
	return context.WithValue(ctx, appKey{}, info)
}

func FromContext(ctx context.Context) (info Info, ok bool) {
	info, ok = ctx.Value(appKey{}).(Info)
	return
}
