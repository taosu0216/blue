package domain

import (
	"blue/common/config"
	"encoding/json"
	"log"
	"net"
)

type GatewayReq struct {
	Fd   int32  `json:"fd"`
	Data []byte `json:"data"`
}

var (
	UnixConn *net.UnixConn
)

type StateReq struct {
	Endpoint string `json:"endpoint"`
	Fd       int32  `json:"fd"`
	Data     []byte `json:"data"`
}

//func init() {
//	get()
//}

func get() {
	UnixConn = NewSocketFile()
}

func NewSocketFile() *net.UnixConn {
	conn, err := net.DialUnix("unix", &net.UnixAddr{
		Name: config.GetDomainSocketPullPath(),
		Net:  "unix",
	}, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}

func SendMsg(fd int32, payLoad []byte) {
	req := &GatewayReq{
		Fd:   fd,
		Data: payLoad,
	}
	encoded, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	_, err = UnixConn.Write(encoded)
	if err != nil {
		panic(err)
	}
}
