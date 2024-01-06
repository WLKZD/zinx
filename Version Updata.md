## ZinxV0.6

##### 消息管理模块的属性和方法

![image-20240106165120383](C:\Users\armstrong\AppData\Roaming\Typora\typora-user-images\image-20240106165120383.png)

#### 继承到Zinx0.5版本中更新为0.6版本



![image-20240106165256557](C:\Users\armstrong\AppData\Roaming\Typora\typora-user-images\image-20240106165256557.png)‘

0.5版本Server模块和Connection模块只能绑定一个固定的路由Router

0.6版本的消息管理模块MsgHandle采用map存储MessageID和路由的对应关系，实现多路由支持

服务端根据客户端消息MessageID调用相应的路由Handle