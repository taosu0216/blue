package client

import (
	"blue/common/brpc"
	"blue/common/brpc/config"
	"blue/gateway/rpc/service"
	"context"
	"fmt"
	"time"
)

// grpc生成的代码,这里是一个client,用于调用gateway的服务
var gatewayClient service.GatewayClient

func initGatewayClient() {
	bCli, err := brpc.NewBClient(config.GetGatewayServiceName())
	if err != nil {
		panic(err)
	}
	gatewayClient = service.NewGatewayClient(bCli.Conn())
}

// state 下打印的
func Push(ctx *context.Context, fd int32, payLoad []byte) error {
	rpcCtx, _ := context.WithTimeout(*ctx, 100*time.Second)

	// TODO：同机器部署用domain socket
	// domain.SendMsg(fd, payLoad)

	resp, err := gatewayClient.Push(rpcCtx, &service.GatewayRequest{Fd: fd, Data: payLoad})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
	return nil
}
