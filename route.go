package socket

import (
	"io"
	"log"
	"net"
	"reflect"
	"time"
)

type Message struct {
	Cmd         uint32
	Version     uint32
	UserProgram uint32
	Param       interface{}
}


func business(conn Conn, msg *Message) {
	if _, ok := msg.Param.(*ReplyFailResult); ok {
		_, err := writeResult(conn, msg)
		if err != nil {
			Log("conn.WriteResult()", err)
		}
		return
	}

	for _, v := range Routers {
		pred := v[0]
		act := v[1]
		if pred.(func(entry Message) bool)(*msg) {
			act.(Controller).Excute(msg)
			_, err := writeResult(conn, msg)
			if err != nil {
				Log("conn.WriteResult()", err)
			}
			return
		}
	}

	_, err := writeError(conn, msg)
	if err != nil {
		Log("conn.WriteError()", err)
	}
}

func reader(conn Conn, readerChannel <-chan *Message, timeout int) {
	for {
		select {
		case msg := <-readerChannel:
			//conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
			business(conn, msg)
			conn.Close()
			return
		case <-time.After(time.Duration(timeout) * time.Second):
			conn.Close()
			Log("connection is closed.")
			return
		}
	}
}

// HandleConnection 处理长连接
func HandleConnection(conn Conn, timeout int) {
	//声明一个临时缓冲区，用来存储被截断的数据
	var tmpBuffer []byte

	//声明一个管道用于接收解包的数据
	readerChannel := make(chan *Message, 16)
	go reader(conn, readerChannel, timeout)

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				continue
			}
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				Log("exit goroutine.")
				return
			}
			Log(conn.RemoteAddr().String(), " connection error: ", err, reflect.TypeOf(err))
			return
		}
		tmpBuffer = unpack(append(tmpBuffer, buffer[:n]...), readerChannel)
	}

}

func Log(v ...interface{}) {
	log.Println(v...)
}

// 接口和该接口的参数结构
var Api map[int]func(dataLen int) interface{}

// 路由
var Routers [][2]interface{}

// Route 路由注册
func Route(rule interface{}, controller Controller) {
	switch rule.(type) {
	case func(entry Message) bool:
		{
			var arr [2]interface{}
			arr[0] = rule
			arr[1] = controller
			Routers = append(Routers, arr)
		}
	case int:
		{
			defaultJudge := func(entry Message) bool {
				cmd := int(entry.Cmd)
				if cmd > rule.(int) && cmd < rule.(int)+256 {
					if _, ok := Api[cmd]; ok {
						return true
					} else {
						return false
					}
				}
				return false
			}
			var arr [2]interface{}
			arr[0] = defaultJudge
			arr[1] = controller
			Routers = append(Routers, arr)
		}
	default:
		Log("Something is wrong in Router")
	}
}

func init() {
	Api = make(map[int]func(dataLen int) interface{}, 0)
	Routers = make([][2]interface{}, 0, 10)
}
