package domain

import (
	"blue/common/config"
	"blue/gateway"
	"blue/gateway/rpc/service"
	"context"
	"encoding/json"
	"net"
)

// 需要开个携程一直监听
func ListenUnixConn() {
	listen, err := net.ListenUnix("unix", &net.UnixAddr{
		Name: config.GetDomainSocketPullPath(),
		Net:  "unix",
	})
	if err != nil {
		panic(nil)
	}

	conn, err := listen.Accept()
	if err != nil {
		panic(err)
	}

	var req service.GatewayRequest
	err = json.NewDecoder(conn).Decode(&req)

	c := context.TODO()
	gateway.CmdChan <- &service.CmdContext{
		Ctx: &c,
		Cmd: 2,
		FD:  int(req.Fd),
	}
}
