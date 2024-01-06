package ziface

/*
	消息管理抽象层
*/

type IMsgHandle interface {
	//调度/执行对应的Router消息处理方法
	DoMsgHandle(request IRequest)
	//添加路由
	AddRouter(msgId uint32, router IRouter)
}
