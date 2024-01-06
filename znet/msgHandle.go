package znet

import (
	"fmt"
	"strconv"
	"zinx/ziface"
)

/*
消息处理模块的实现
*/
type MsgHandle struct {
	Apis map[uint32]ziface.IRouter
}

// 初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

func (mh *MsgHandle) DoMsgHandle(request ziface.IRequest) {
	router, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " id NOT FOUNT! Need Register")
		return
	}
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router ziface.IRouter) {
	//1 判断当前msg绑定的Api方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeat api, msgId = " + strconv.Itoa(int(msgId)))
	}
	//2添加msg与API的绑定是关系
	mh.Apis[msgId] = router
	fmt.Println("Add api MsgID = ", msgId, " succ! ")
}
