package api

// 定义路由接口的抽象层
type IRouter interface {

	// 处理conn业务之前的钩子方法Hook
	PreHandle(request IRequest)

	// 处理conn业务的主方法
	Handle(request IRequest)

	// 处理conn业务之后的方法
	AfterHandle(request IRequest)
}
