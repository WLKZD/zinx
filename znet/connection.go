package znet

import (
	"fmt"
	"net"
	"zinx/utils"
	"zinx/ziface"
)

/*
	链接模块
*/

type Connection struct {
	//当前链接的socket TCP 套接字
	Conn *net.TCPConn

	//链接的ID
	ConnID uint32

	//当前的链接状态
	isClose bool

	//当前链接锁绑定的处理业务方法API
	//handleAPI ziface.HandleFunc

	//告知当前链接已经推出的/停止 channel
	ExitChan chan bool

	//该链接处理的方法Router
	Router ziface.IRouter
}

// 初始化链接模块的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:     conn,
		ConnID:   connID,
		Router:   router,
		isClose:  false,
		ExitChan: make(chan bool),
	}
	return c
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Groutine is running...")
	defer fmt.Println("connId=", c.ConnID, "Reader is exit,remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//读取客户端的数据到buf中
		buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			continue
		}

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			data: buf,
		}

		go func(request ziface.IRequest) {
			//从路由中，找到注册绑定的Conn对应的router调用
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

// 启动链接，让当前的链接准备开始工作
func (c Connection) Start() {
	fmt.Println("Conn Star()... ConnID=", c.ConnID)
	//启动从当前链接的读数据的业务
	go c.StartReader()
	//TODO启动从当前链接写数据的业务

}

// 停止链接，结束当前链接的工作
func (c Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID=", c.ConnID)

	//如果当前链接已经关闭
	if c.isClose == true {
		return
	}
	c.isClose = true

	//关闭socket链接
	c.Conn.Close()

	//回收资源
	close(c.ExitChan)
}

// 获取当前链接的绑定scoket conn
func (c Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的链接ID
func (c Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态 IP port
func (c Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 发生数据，将数据发送给远程的客户端
func (c Connection) Send(bytes []byte) error {
	return nil
}
