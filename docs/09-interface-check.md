# 接口实现检查

为了保障 Rum 和 RumGroup 实现的方法一致， 就需要对这两个类型进行 接口实现检查。 而在 go 中， 类与接口并无直接关系， 一个类 只要实现了 某个接口 的所有方法， 那这个类就是该接口的实现。 为了保障在编码阶段就能检查， go 提供了一个特殊的语法糖。

使用 golang 检查是否实现接口的语法糖

```go
// var _ interfaceName = structInstance
var _ IRumRoutes = &RumGroup{}
```

在 gin 中， 也有相关接口检查的代码实践。