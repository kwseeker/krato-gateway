# kratos-gateway

揣摩分析 kratos gateway 从0到1设计与实现。

> 大部分开源项目初始提交都是一个很大的提交，可能项目本来并不是开源项目，只是中间开源了，代码直接copy到github上，导致原始提交丢失，因此无法看到项目怎么从零开始设计优化拓展的。
>
> 这里结合自己对源码的理解，复现kratos gateway从0到1的过程。

## Modules：

+ **gateway-01**

  本质上还是一个web应用。

  包括：

  + context 上下文生命周期管理

  + 关闭信号监听

  + Web server 配置&生命周期&请求处理

