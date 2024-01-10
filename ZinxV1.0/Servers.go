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

// Test Handle
func (this *PingRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		",data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

// hello Zinx test  自定义路由
type HelloZinxRouter struct {
	znet.BaseRouter
}

// Test Handle
func (this *HelloZinxRouter) Handle(request ziface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
		",data = ", string(request.GetData()))
	err := request.GetConnection().SendMsg(201, []byte("Hello Welcome to Zinx"))
	if err != nil {
		fmt.Println(err)
	}
}

// 创建链接之后执行钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("===> DoConnectionBegin is Called...")
	if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
		fmt.Println(err)
	}

	//给当前的链接设置一些属性
	fmt.Println("Set conn New,Hoe...")
	conn.SetProperty("Name", "YJZ")
	conn.SetProperty("Birthday", "0502")
}

// 销毁链接之前执行钩子函数
func DoConnectionLost(conn ziface.IConnection) {
	fmt.Println("===> DoConnectionLost is Called...")
	fmt.Println("connID = ", conn.GetConnID(), " is lost")

	//获取链接设置一些属性
	fmt.Println("Set conn New,Hoe...")
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name = ", name)
	}
	if birthday, err := conn.GetProperty("Birthday"); err == nil {
		fmt.Println("Birthday = ", birthday)
	}
}

func main() {
	// 1创建一个Server句柄，使用Zinx的api
	s := znet.NewServer("[zinx V0.5]")

	// 2注册链接Hook钩子函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	// 3给当前zinx框架添加i一个自定义Router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloZinxRouter{})

	// 3启动server
	s.Serve()
}
