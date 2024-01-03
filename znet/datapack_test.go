package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 只是负责测试datapack拆包、封包的单元测试
func TestDataPack(t *testing.T) {
	/*
		模拟的服务器
	*/
	//1 创建socketTCP
	linster, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen error ", err)
		return
	}
	//创建一个go承载 复制从客户端处理业务
	go func() {
		//2从客户端读取数据，拆包处理
		for {
			conn, err := linster.Accept()
			if err != nil {
				fmt.Println("server accepct error ", err)
			}
			go func(conn net.Conn) {
				//处理客户端请求
				//---->拆包过程<----
				//定义一个拆包的对象dp
				dp := NewDataPack()
				for {
					//1第一次从conn读，把包的head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error ", err)
						return
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack error ", err)
					}
					if msgHead.GetMsgLen() > 0 {
						//msg是有数据的，需要进行第二次读取
						//2第二次从conn读，根据head中的dataLen再读取data内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())

						//根据datalen的长度再次从io流中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data error ", err)
							return
						}

						//完整的一个消息已经读取完毕
						fmt.Println("-->Recv MsgID:", msg.Id, "datalen=", msg.DataLen, "data:", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	/*
		模拟的客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err ", err)
		return
	}

	//创建一个封包对象 dp
	dp := NewDataPack()

	//模拟粘包过程，封装两个msg一同发送
	//封装第一个msg1包
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	senData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("clien pack msg1 error ", err)
		return
	}

	//封装第二个msg2包
	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', '0', '!', '!'},
	}
	senData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("clien pack msg1 error ", err)
		return
	}
	//使两个包黏在一起
	senData1 = append(senData1, senData2...)

	//一次性发送给服务器
	conn.Write(senData1)
}
