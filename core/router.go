package core

import "PhalGo-TCP/core/api"

type BaseRouter struct {

}

// 让用户继承重写
func(br *BaseRouter) PreHandle(request api.IRequest){

}

func (br *BaseRouter)Handle(request api.IRequest){

}

func (br *BaseRouter) AfterHandle(request api.IRequest){

}
