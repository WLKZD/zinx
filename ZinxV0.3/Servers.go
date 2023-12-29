package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

/*
	基于Zinx框架来开发的 服务器应用程序
*/

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test PreHandle
func (this *PingRouter) PreHandle(request ziface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before error")
	}
}

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call Router Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("call back ping...ping...ping... error")
	}
}

// Test PostHandle
func (this *PingRouter) PostHandle(request ziface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after error")
	}
}

func main() {
	// 1创建一个Server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.2]")
	// 2给当前zinx框架添加i一个自定义Router
	s.AddRouter(&PingRouter{})
	// 3启动server
	s.Serve()
}
