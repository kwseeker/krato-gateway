package biz

import "github.com/google/wire"

// ProviderSet is biz providers.
// 用于依赖注入
var ProviderSet = wire.NewSet(NewEchoUsecase)
