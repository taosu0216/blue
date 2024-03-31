package domain

import (
	"blue/common/config"
	"encoding/json"
	"log"
	"net"
)

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
		Name: config.GetDomainSocketPushPath(),
		Net:  "unix",
	}, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}

func SendMsg(endpoint string, fd int32, playLoad []byte) {
	req := &StateReq{
		Endpoint: endpoint,
		Fd:       fd,
		Data:     playLoad,
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
