package client

import (
	"blue/common/brpc"
	"blue/common/config"
	"blue/state/rpc/service"
	"context"
	"time"
)

var stateClient service.StateClient

func initStateClient() {
	bCli, err := brpc.NewBClient(config.GetStateServiceName())
	if err != nil {
		panic(err)
	}
	stateClient = service.NewStateClient(bCli.Conn())
}

func SendMsg(ctx *context.Context, endpoint string, fd int32, playLoad []byte) error {
	rpcCtx, _ := context.WithTimeout(*ctx, 100*time.Millisecond)
	_, err := stateClient.SendMsg(rpcCtx, &service.StateRequest{
		Endpoint: endpoint,
		Fd:       fd,
		Data:     playLoad,
	})
	if err != nil {
		panic(err)
	}
	return nil
}
