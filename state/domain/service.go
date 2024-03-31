package domain

import (
	"blue/common/config"
	"blue/state"
	"blue/state/rpc/service"
	"context"
	"encoding/json"
	"net"
)

// 需要开个携程一直监听
func ListenUnixConn() {
	listen, err := net.ListenUnix("unix", &net.UnixAddr{
		Name: config.GetDomainSocketPushPath(),
		Net:  "unix",
	})
	if err != nil {
		panic(nil)
	}

	conn, err := listen.Accept()
	if err != nil {
		panic(err)
	}

	var req service.StateRequest
	err = json.NewDecoder(conn).Decode(&req)

	c := context.TODO()
	state.CmdChan <- &service.CmdContext{
		Ctx:      &c,
		Cmd:      2,
		FD:       int(req.Fd),
		Endpoint: req.Endpoint,
	}
}
