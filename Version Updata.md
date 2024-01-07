## ZinxV0.6

消息管理模块的属性和方法

![image-20240106165120383](C:\Users\armstrong\AppData\Roaming\Typora\typora-user-images\image-20240106165120383.png)

##### 继承到Zinx0.5版本中更新为0.6版本



![image-20240106165256557](C:\Users\armstrong\AppData\Roaming\Typora\typora-user-images\image-20240106165256557.png)‘

0.5版本Server模块和Connection模块只能绑定一个固定的路由Router

0.6版本的消息管理模块MsgHandle采用map存储MessageID和路由的对应关系，实现多路由支持

服务端根据客户端消息MessageID调用相应的路由Handle



## ZinxV0.7

##### 实现读写协程分离

![image-20240107165321796](C:\Users\armstrong\AppData\Roaming\Typora\typora-user-images\image-20240107165321796.png)

0.6版本服务端调用Connection.SendMsg方法直接将数据写入到conn中

0.7版本在Connection对象中新增一个无缓冲管道msgChan和Connection.Writer方法，

Connection.SendMsg方法将数据写入msgChan，Connection.Writer方法将msgChan中的数据写入到conn中

如果客户端关闭连接，Connection.Reader方法读取数据失败退出并调用Conneciton.Stop方法，执行c.ExitChan <- true，

Connection.Writer读到c.ExitChan里的数据为True时也退出



## ZinxV0.8

#### 消息队列及多任务

##### 创建消息队列





![image-20240107211133194](C:\Users\armstrong\AppData\Roaming\Typora\typora-user-images\image-20240107211133194.png)

##### 创建及启动Worker工作池

![image-20240107211940818](C:\Users\armstrong\AppData\Roaming\Typora\typora-user-images\image-20240107211940818.png)

##### 发送消息给消息队列





![image-20240107211911690](C:\Users\armstrong\AppData\Roaming\Typora\typora-user-images\image-20240107211911690.png)

##### 将消息队列集成到 Zinx 框架

![image-20240107212137659](C:\Users\armstrong\AppData\Roaming\Typora\typora-user-images\image-20240107212137659.png)

0.7版本对于每个连接Conneciton都会调用go c.MsgHandle.DoMsgHandle(&req)

0.8版本服务端会创建一个工作池，规定了工作池的Woker数量，既创建的goroutine数量，每个Woker对应一个消息队列，安装分配算法（现在是根据ConnectionID%WorkerPoolSize）将消息发送给Worker工作池即可

