# gateway-02

+ **添加基础模块**

    + **日志**
    + **pprof**
    + **配置**
    + **监控**
    
+ **创建一个后端微服务**（暂时只支持HTTP协议），实现gateway转发

+ **添加服务注册&发现**



## 日志

日志用的 kratos log 包下的实现。

代码封装的也比较简单，直接看图吧

![](../img/kratos-log-arch.png)

流程图：

![](../img/kratos-log-workflow.png)

> 注意事项：
>
> 线上不要用 Fatal等级的方法，打印日志，因为会在打印日志后中断程序运行。



## pprof

```go
_ "net/http/pprof"
go func() {
	_ = http.ListenAndServe("0.0.0.0:17070", nil)
}()
```

TODO: pprof 怎么知道是绑定的哪个端口？



## 配置

日志用的 kratos [config](https://go-kratos.dev/docs/component/config) 包下的实现。

框架内置了[本地文件file](https://github.com/go-kratos/kratos/tree/main/config/file)和[环境变量env](https://github.com/go-kratos/kratos/tree/main/config/env)的实现。

另外，在[contrib/config](https://github.com/go-kratos/kratos/tree/main/contrib/config)下面，也提供了如下的配置中心的适配供使用：

- [apollo](https://github.com/go-kratos/kratos/tree/main/contrib/config/apollo)
- [consul](https://github.com/go-kratos/kratos/tree/main/contrib/config/consul)
- [etcd](https://github.com/go-kratos/kratos/tree/main/contrib/config/etcd)
- [kubernetes](https://github.com/go-kratos/kratos/tree/main/contrib/config/kubernetes)
- [nacos](https://github.com/go-kratos/kratos/tree/main/contrib/config/nacos)

**数据结构**：

```go
type config struct {
	opts      options	//姑且叫元配置吧，用于解析配置的配置选项，比如对接多种配置源
	reader    Reader		//使用options配置的工具解析配置数据	
	cached    sync.Map		//读取配置字段的缓存，检索很耗时么？不应该吧
	observers sync.Map		//也类似缓存
	watchers  []Watcher		//监听配置源内容变化的监听器（对配置文件的监听是用epoll监听的，只讨论Linux环境）
	log       *log.Helper	//默认是log.DefaultLogger
}

type options struct {	//元配置
	sources  []Source		//配置源（环境变量、本地文件、一些配置中心中间件）
	decoder  Decoder		//解析配置内容的解码器，默认的解码器（defaultDecoder，代理）会按照文件类型解析（根据类型获取对应的Decoder）
    						// 当前kratos支持 form\json\proto\xml\yaml
	resolver Resolver		//对解析完毕后的map结构进行再次处理，默认resolver（defaultResolver）会对配置中的占位符进行填充
	logger   log.Logger
}

type KeyValue struct {	//配置加载到内存后的数据结构
	Key    string				//比如加载配置文件："config.yaml"
	Value  []byte				//比如加载配置文件,这里就是配置文件内容的byte数组
	Format string				//比如加载配置文件，这里就是文件类型，如"yaml"
}

type reader struct {	//用于读取加载后数据（KeyValue）,使用options中指定的decoder resolover解析并转成map
	opts   options					//config中的元配置
	values map[string]interface{}	//解析后的配置数据（map格式）
	lock   sync.Mutex
}
```

参考Demo: config_test.go。



## 监控



## 创建后端微服务

根据Kratos文档通过命令创建，修改，编译，启动即可。

这里创建了个很简单的Echo服务。

```shell
curl --location --request GET 'http://127.0.0.1:18001/echo/lee'
```



## [Consul](https://kingfree.gitbook.io/consul/)服务注册&发现

本地启动开发模式：

```
consul agent -dev
# 集群节点
consul members
# 重新加载配置文件
consul reload
# 优雅关闭节点
consul leave
# 查询所有注册的服务
consul catalog services
```

默认启动配置：

```verilog
==> Starting Consul agent...
              Version: '1.14.3'
           Build Date: '2022-12-13 17:13:55 +0000 UTC'
              Node ID: '85c3f184-b06e-6ae7-1d1b-38939f3f974b'
            Node name: 'Lee-Home'
           Datacenter: 'dc1' (Segment: '<all>')
               Server: true (Bootstrap: false)
          Client Addr: [127.0.0.1] (HTTP: 8500, HTTPS: -1, gRPC: 8502, gRPC-TLS: 8503, DNS: 8600)
         Cluster Addr: 127.0.0.1 (LAN: 8301, WAN: 8302)
    Gossip Encryption: false
     Auto-Encrypt-TLS: false
            HTTPS TLS: Verify Incoming: false, Verify Outgoing: false, Min Version: TLSv1_2
             gRPC TLS: Verify Incoming: false, Min Version: TLSv1_2
     Internal RPC TLS: Verify Incoming: false, Verify Outgoing: false (Verify Hostname: false), Min Version: TLSv1_2
```

web页面：http://localhost:8500/

使用配置文件将微服务注册到Consul：

```shell
#微服务注册的配置文件和consul启动配置文件放同一个目录
consul agent -dev -config-dir=/home/lee/mywork/go/src/github.com/kwseeker/kratos-gateway/deploy/consul/registray
```

在web页面Services栏中看到echo服务，说明注册成功。

另外也支持在微服务中通过请求接口注册，像Java系微服务一样由服务提供者主动注册。

[Kratos 注册服务实例](https://go-kratos.dev/docs/component/registry#%E6%B3%A8%E5%86%8C%E6%9C%8D%E5%8A%A1%E5%AE%9E%E4%BE%8B)

接下来看Gateway如果实现请求转发。

HTTP -> Proxy -> Router -> Middleware -> Client -> Selector -> Node

**客户端代理（Proxy）：**

```go
type Proxy struct {
	router            atomic.Value		//mux.NewRouter()，定义核心路由规则
	log               *log.Helper
	clientFactory     client.Factory	//用于创建连接到后台微服务客户端的工厂函数（里面包含连接consul的http客户端及连接配置）
	middlewareFactory middleware.Factory
}
```

**路由规则**：

` p.Update(bc)`会创建客户端连接、然后装饰上所有配置的middleware组件、然后再封装成http.Handler方法；最后将路由信息和http.Handler注册到mux路由。

请求处理链路：server.Handler -> http2.h2cHandler -> Proxy.router -> mux.router.Handle -> 连接指定微服务的客户端连接。

```go
type Router struct {
	// Configurable Handler to be used when no route matches.
	NotFoundHandler http.Handler
	// Configurable Handler to be used when the request method does not match the route.
	MethodNotAllowedHandler http.Handler
	// Routes to be matched, in order.
    // 路由数组，查找匹配的处理方法是 for range 循环查找，估计是考虑有些通配符形式的路径只能循环遍历匹配
	routes []*Route
	// Routes by name for URL building.
	namedRoutes map[string]*Route
	// If true, do not clear the request context after handling the request.
	//
	// Deprecated: No effect, since the context is stored on the request itself.
	KeepContext bool
	// Slice of middlewares to be called after a match is found
	middlewares []middleware
	// configuration shared with `Route`
	routeConf
}

type Gateway struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string        `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Hosts       []string      `protobuf:"bytes,2,rep,name=hosts,proto3" json:"hosts,omitempty"`
	Endpoints   []*Endpoint   `protobuf:"bytes,3,rep,name=endpoints,proto3" json:"endpoints,omitempty"`
	Middlewares []*Middleware `protobuf:"bytes,4,rep,name=middlewares,proto3" json:"middlewares,omitempty"`
}

type Endpoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Method      string               `protobuf:"bytes,1,opt,name=method,proto3" json:"method,omitempty"`
	Path        string               `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	Description string               `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Protocol    Protocol             `protobuf:"varint,4,opt,name=protocol,proto3,enum=gateway.config.v1.Protocol" json:"protocol,omitempty"`
	Timeout     *durationpb.Duration `protobuf:"bytes,5,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Middlewares []*Middleware        `protobuf:"bytes,6,rep,name=middlewares,proto3" json:"middlewares,omitempty"`
	Backends    []*Backend           `protobuf:"bytes,7,rep,name=backends,proto3" json:"backends,omitempty"`
	Retry       *Retry               `protobuf:"bytes,8,opt,name=retry,proto3" json:"retry,omitempty"`
}
```

