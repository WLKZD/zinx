# zinx

## ZinxV0.5

### 服务端

#### 1  使用Zinx的api，创建一个Server句柄

    s := znet.NewServer("[zinx V0.5]")


**sever.go**下的**NewServer**函数创建**Server**句柄，绑定**utils.GlobalObjectd**的字段，
导入**zinx/utils**包时执行一次**init**函数，在**init**函数中调用**Reload**方法把配置文件
**conf/zinx.json**中的数据解析到**GlobalObject**中

    func NewServer(name string) ziface.IServer
    func (g *GlobalObj) Reload()
    func init()



#### 2  给当前zinx框架添加一个自定义Router

    s.AddRouter(&PingRouter{})

自定义路由**PingRouter**继承**Zinx**框架的基础路由**znet.BaseRouter**，
重写**Handle**方法，打印输出**request**的msg信息，
同时回写“ping..ping...ping.."到原生**conn**中

    type PingRouter struct {
        znet.BaseRouter
    }
    func (this *PingRouter) Handle(request ziface.IRequest){
    	fmt.Println("Call Router Handle")
    	//先读取客户端的数据，再回写ping...ping...ping
    	fmt.Println("recv from client: msgID = ", request.GetMsgID(),
    		",data = ", string(request.GetData()))
    	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
    	if err != nil {
    		fmt.Println(err)
    	}
    }



#### 3  启动server

    s.Serve()

**s.Serve**中调用**s.Start**方法


    func (s *Server) Serve() {
        //启动server的服务功能
        s.Start()
        //TODO 做一些启动服务器之后的额外业务
        //阻塞状态
        select {}
    }

**s.Start**中的**delConn**是**Zinx**框架中定义的**Connection**结构体，
调用**delConn.Start**方法执行链接绑定的业务

    func (s *Server) Start() {
        go func() {
           //1 获取一个TCP的Addr
           addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
           //2 监听服务器的地址
           listener, err := net.ListenTCP(s.IPVersion, addr)
           //3 阻塞的等待客服端链接，处理客户端连接业务(读写)
           for {
              //如果有客户端连接过来，阻塞会返回
              conn, err := listener.AcceptTCP()
              //将处理新连接的业务方法和conn进行绑定，得到我们的链接模块
              delConn := NewConnection(conn, cid, s.Router)
              cid++
              //启动当前的链接业务方法
              go delConn.Start()
           }
        }()
    }


启动链接，当前链接只有读数据的业务**c.StartReader**


    func (c Connection) Start() {
        fmt.Println("Conn Star()... ConnID=", c.ConnID)
        //启动从当前链接的读数据的业务
        go c.StartReader()
        //TODO启动从当前链接写数据的业务
    }


链接的读业务方法，将**c.GetTCPConnection**得到的原生**conn**中的二进制流拆包后，
读到**msg**消息中，初始化**req**对象，调用**c.Router**的三个Handle方法,
当前**c.Router**就是自定义的**PingRouter**，只重写了**Handle**方法


    func (c *Connection) StartReader() {
        for {
           //创建一个拆包解包对象
           dp := NewDataPack()
    
           //读取客户端的Msg Head 二进制流 8个字节
           headData := make([]byte, dp.GetHeadLen())
           io.ReadFull(c.GetTCPConnection(), headData)
    
           //拆包，得到msgID和msgDatalen 放在msg消息中
           msg, err := dp.UnPack(headData)
    
           //得到dataLen，再次读取Data，放在msg.Data
           var data []byte
           if msg.GetMsgLen() > 0 {
              data = make([]byte, msg.GetMsgLen())
           }
           msg.SetData(data)
           //得到当前conn数据的Request请求数据
           req := Request{
              conn: c,
              msg:  msg,
           }
           
           go func(request ziface.IRequest) {
              //从路由中，找到注册绑定的Conn对应的router调用
              c.Router.PreHandle(request)
              c.Router.Handle(request)
              c.Router.PostHandle(request)
           }(&req)
    
        }
    }




### 客户端

#### 1  直接链接远程服务器，得到一个conn链接

	conn, err := net.Dial("tcp", "127.0.0.1:8999")



#### 2  发送对**message**消息封包后的二进制流**binaryMsg**

**znet.NewdataPack**初始化封包对象**dp**，调用 **dp.Pack** 方法封包生成字节流切片**binaryMsg**

**znet.NewMsgPackage**初始化**message**对象，ID=0,Data:"ZinxV0.5 client Test Message"

调用**conn.Write**方法将字节流写入到**conn**中

    dp := znet.NewDataPack()
    binartMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("ZinxV0.5 client Test Message")))
    conn.Write(binaryMsg)



#### 3  对服务器返回的二进制流进行拆包，当前服务器默认发送 MsgID=1，Data:pingpingping

初始化大小为头部数据长度(**dp.GetHeaadLen()**默认返回8)的切片**binaryHead**，

调用**io.ReadFull**，先读取**conn**流中的head部分，得到ID和dataLen写入到**binaryHead**中

    binaryHead := make([]byte, dp.GetHeadLen())
     _, err := io.ReadFull(conn, binaryHead)

调用**dp.Unpack**拆包方法将二进制的**binaryHead**拆包到**msgHead**结构体中,

    msgHead, err := dp.UnPack(binaryHead)

再根据**msgHead.GetMsgLen**进行第二次读取，将**Data**读出来

类型断言将接口**msgHead**转为结构体类型**msg**

    if msgHead.GetMsgLen() > 0 {
        msg := msgHead.(*znet.Message)
        msg.Data = make([]byte, msg.GetMsgLen())
    	io.ReadFull(conn, msg.Data)
    }

