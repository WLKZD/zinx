package znet

import (
	"errors"
	"fmt"
	"net"
	"zinx/ziface"
)

// IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的IP
	IP string
	//服务器监听的端口
	Port int
}

// 定义当前客户端链接的所绑定handle api(目前这个handle是写死的，以后优化应该由用户自定义handle方法)
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显的业务
	fmt.Println("[Conn Handle] CallbackToClient")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallbackToClient error")
	}
	return nil

}

func (s *Server) Start() {
	fmt.Printf("[start]Server Linstenner at IP :%s,Port %d is starting\n", s.IP, s.Port)

	go func() {
		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addt error ：", err)
			return
		}
		//2 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		fmt.Println("start Zinx server succ, ", s.Name, " succ, Linstenning...")
		var cid uint32 = 0

		//3 阻塞的等待客服端链接，处理客户端连接业务(读写)
		for {
			//如果有客户端连接过来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			//将处理新连接的业务方法和conn进行绑定，得到我们的链接模块
			delConn := NewConnection(conn, cid, CallBackToClient)
			cid++

			//启动当前的链接业务方法
			go delConn.Start()

			//已经与客户端建立链接，做一些业务，做一个最基本的最大512字节长度的回显业务
			/*
				go func() {
					for {
						buf := make([]byte, 512)
						cnt, err := conn.Read(buf)
						if err != nil {
							fmt.Println("recv buf err", err)
							continue
						}
						fmt.Printf("read client buf:%s,cnt %d\n", buf, cnt)
						//回显功能
						if _, err := conn.Write(buf[:cnt]); err != nil {
							fmt.Println("write back buf err", err)
							continue
						}
					}
				}()
			*/
		}
	}()
}

func (s *Server) Stop() {
	//TODO 将一些服务器的资源、状态或者一下已经开辟的链接信息 进行停止或者回收
}

func (s *Server) Serve() {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {}

}

//初始化Server模块的方法

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
