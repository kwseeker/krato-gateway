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
	opts      options		//姑且叫元配置吧，用于解析配置的配置选项，因为可以对接多种配置源
	reader    Reader
	cached    sync.Map
	observers sync.Map
	watchers  []Watcher		//监听配置源内容变化的监听器
	log       *log.Helper	//默认是log.DefaultLogger
}

type options struct {
	sources  []Source		//配置源
	decoder  Decoder		//
	resolver Resolver
	logger   log.Logger
}

type KeyValue struct {
	Key    string
	Value  []byte
	Format string
}
```



## 监控



