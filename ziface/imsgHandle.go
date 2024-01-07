package ziface

/*
	消息管理抽象层
*/

type IMsgHandle interface {
	//调度/执行对应的Router消息处理方法
	DoMsgHandle(request IRequest)
	//为消息添加具体的处理逻辑
	AddRouter(msgId uint32, router IRouter)
	//启动Worker工作池
	StartWorkerPool()
	//将消息发送给消息队列处理
	SendMsgToTaskQueue(request IRequest)
}
