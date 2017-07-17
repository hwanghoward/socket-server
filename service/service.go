package main

import (
	"log"
	"net"
	"os"

	"socket-server"
	. "socket-server/service/controller"
)

func CheckError(err error) {
	if err != nil {
		log.Printf("Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func Log(v ...interface{}) {
	log.Println(v...)
}

func init() {
	socket.Api[CMD_REQUEST_CREATE_RAID]     = NewParams(NewCreateRaidParam, 32)
	socket.Api[CMD_REQUEST_DETORY_RAID]     = NewParams(NewRaidIdentity, 16)
	socket.Api[CMD_REQUEST_CREATE_HOTSPARE] = NewParams(NewHotSpareParam, 20)
	socket.Api[CMD_REQUEST_DETORY_HOTSPARE] = NewParams(NewHotSpareParam, 20)
	socket.Api[CMD_REQUEST_FORMAT_DISK]     = NewParams(NewDiskIdentity, 4)

	socket.Api[CMD_REQUEST_QUERY_ENCLORSURE_INFO]   = NewParams(NewNoParam, 0)
	socket.Api[CMD_REQUEST_QUERY_RAID_INFO]         = NewParams(NewRaidIdentity, 0, 16)
	socket.Api[CMD_REQUEST_QUERY_DISK_INFO]         = NewParams(NewNoParam, 0)
	socket.Api[CMD_REQUEST_QUERY_RAID_REBUILD_INFO] = NewParams(NewRaidIdentity, 16)

	var op Operator
	var qy Query
	socket.Route(0x200, &op)
	socket.Route(0x300, &qy)
}

func main() {
	netListen, err := net.Listen("tcp", "0.0.0.0:6060")
	CheckError(err)
	defer netListen.Close()
	Log("Waiting for clients")
	for {
		conn, err := netListen.Accept()
		if err != nil {
			continue
		}
		Log(conn.RemoteAddr().String(), " tcp connect success")
		// 如果此链接超过6秒没有发送新的数据，将被关闭
		go socket.HandleConnection(socket.Conn{conn}, 10)
	}
}
