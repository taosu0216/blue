package state

import (
	"blue/common/brpc"
	"blue/common/config"
	"blue/state/rpc/client"
	"blue/state/rpc/service"
	"context"
	"fmt"

	"google.golang.org/grpc"
)

var CmdChan chan *service.CmdContext

func RunMain(path string) {
	config.Init(path)
	CmdChan = make(chan *service.CmdContext, config.GetStateCmdChannelNum())

	s := brpc.NewBServer(
		// 统统都是赋值操作,甚至可以理解为就是多封装了两层
		brpc.WithServiceName(config.GetStateServiceName()),
		brpc.WithIP(config.GetStateServiceAddr()),
		brpc.WithPort(config.GetStateServerPort()),
		brpc.WithWeight(config.GetStateRPCWeight()))

	s.RegisterService(func(server *grpc.Server) {
		service.RegisterStateServer(server, &service.Service{CmdChannel: CmdChan})
	})

	// TODO: 同机器部署gateway和state，使用domain socket通信
	//go func(){
	//	for  {
	//		domain.ListenUnixConn()
	//	}
	//}()

	client.Init()

	go cmdHandler()

	s.Start(context.TODO())
}

func cmdHandler() {
	for cmd := range CmdChan {
		switch cmd.Cmd {
		case service.CancelConnCmd:
			fmt.Printf("cancelconn endpoint:%s, fd:%d, data:%+v", cmd.Endpoint, cmd.FD, cmd.Playload)
		case service.SendMsgCmd:
			// fmt.Println("cmdHandler", string(cmd.Playload))
			_ = client.Push(cmd.Ctx, int32(cmd.FD), cmd.Playload)
		}
	}
}
