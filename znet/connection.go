package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

/*
	链接模块
*/

type Connection struct {
	//当前Conn隶属于哪个Server
	TcpServer ziface.IServer

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

	//无缓冲管道，用于读、写Goroutine之间的消息通信
	msgChan chan []byte

	//消息的管理MsgID和对应的处理业务API关系
	MsgHandle ziface.IMsgHandle

	//链接属性集合
	property map[string]interface{}

	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

// 初始化链接模块的方法
func NewConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandle ziface.IMsgHandle) *Connection {
	c := &Connection{
		TcpServer: server,
		Conn:      conn,
		ConnID:    connID,
		MsgHandle: msgHandle,
		isClose:   false,
		ExitChan:  make(chan bool),
		msgChan:   make(chan []byte),
		property:  make(map[string]interface{}),
	}
	//将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().Add(c)
	return c
}

// 链接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("[Reader Groutine is running...]")
	defer fmt.Println("[Reader is exit],connId=", c.ConnID, ",remote addr is ", c.RemoteAddr().String())
	defer c.Stop()
	for {
		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端的Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		//拆包，得到msgID和msgDatalen 放在msg消息中
		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack error ", err)
			break
		}
		//得到dataLen，再次读取Data，放在msg.Data
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetData(data)
		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPooleSize > 0 {
			//已经开启了工作池机制，将消息发送给Worker工作池即可
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			go c.MsgHandle.DoMsgHandle(&req)
		}
	}
}

// 写消息Goroutine,专门发送给客户端消息的模块
func (c *Connection) StartWrite() {
	fmt.Println("[Writer Goroutine is running ...]")
	defer fmt.Println("[conn Writer exit!]", c.RemoteAddr().String())

	//不断的阻塞的等待channel的消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			//有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error ", err)
				return
			}
		case <-c.ExitChan:
			//代表Reader已经推出，此时Writer也要退出
			return
		}
	}
}

// 启动链接，让当前的链接准备开始工作
func (c *Connection) Start() {
	fmt.Println("Conn Star()... ConnID=", c.ConnID)
	//启动从当前链接的读数据的业务
	go c.StartReader()
	//TODO启动从当前链接写数据的业务
	go c.StartWrite()

	//按照开发者传递进来的 创建链接之后需要调用的处理业务，执行对应Hook函数
	c.TcpServer.CallOnConnStart(c)
}

// 停止链接，结束当前链接的工作
func (c *Connection) Stop() {
	fmt.Println("Conn Stop().. ConnID=", c.ConnID)

	//如果当前链接已经关闭
	if c.isClose == true {
		return
	}
	c.isClose = true

	//按照开发者传递进来的 销毁链接之前需要调用的处理业务，执行对应Hook函数
	c.TcpServer.CallOnConnStop(c)

	//关闭socket链接
	c.Conn.Close()

	//告知Writer关闭
	c.ExitChan <- true

	//将当前链接从ConnMar中摘除掉
	c.TcpServer.GetConnMgr().Remove(c)

	//回收资源
	close(c.ExitChan)
	close(c.msgChan)
}

// 获取当前链接的绑定scoket conn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前连接模块的链接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的TCP状态 IP port
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()

}

// 提供一个SendMsg方法
func (c *Connection) SendMsg(msgid uint32, data []byte) error {
	if c.isClose == true {
		return errors.New("Connection closed when send msg")
	}

	//data进行封包 MsgDataLen/MsgId/Data
	dp := NewDataPack()
	binaryMsg, err := dp.Pack(NewMsgPackage(msgid, data))
	if err != nil {
		fmt.Println("Pack error msg id = ", msgid)
		return errors.New("Pack error msg")
	}
	//将数据发送给客户端
	c.msgChan <- binaryMsg

	return nil
}

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

// 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

// 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
}
