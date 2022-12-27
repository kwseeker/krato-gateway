package client

import (
	"context"
	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/selector/p2c"
	"github.com/hashicorp/consul/api"
	config "github.com/kwseeker/kratos-gateway/gateway-02/api/gateway/config/v1"
	"google.golang.org/protobuf/types/known/durationpb"
	"os"
	"testing"
)

func TestClient(t *testing.T) {
	logger := log.NewStdLogger(os.Stdout)
	l := log.NewHelper(logger)

	consulAddress := "localhost:8500"
	c := api.DefaultConfig()
	c.Address = consulAddress
	//核心是创建http.Client{}实例，即http客户端
	cli, err := api.NewClient(c)
	if err != nil {
		panic(err)
	}
	r := consul.New(cli)

	//endpoint *config.Endpoint
	//  - path: /echo/*
	//    method: '*'
	//    timeout: 1s
	//    protocol: HTTP
	//    backends:
	//      - target: '127.0.0.1:18001'
	endpoint := &config.Endpoint{
		Method:   "*",
		Path:     "/echo/*",
		Protocol: config.Protocol_HTTP,
		Timeout:  &durationpb.Duration{Seconds: 1},
		//Backends: [1]config.Backend{&config.Backend{Target: "127.0.0.1:18001"},},
	}

	//创建p2c选择器（两次随机选择）
	picker := p2c.New()
	applier := &nodeApplier{
		endpoint:  endpoint,
		logHelper: l,
		registry:  r,
	}
	//应用负载均衡算法，从服务节点列表中选择一个节点
	if err := applier.apply(context.Background(), picker); err != nil {
		t.Fatalf("apply()...")
	}
	//使用选中的服务节点创建连接
	client := &client{
		selector: picker,
		attempts: calcAttempts(endpoint),
		protocol: endpoint.Protocol,
	}
	retryCond, err := parseRetryCondition(endpoint)
	if err != nil {
		t.Fatalf("parseRetryCondition()...")
	}
	client.conditions = retryCond
}
