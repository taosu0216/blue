package gateway

import (
	"blue/common/config"
	"blue/common/tcp"
	"context"
	"errors"
	"fmt"
	"github.com/hardcore-os/plato/gateway/rpc/client"
	"io"
	"net"
)

func RunMain(configPath string) {
	// Run the gateway
	config.Init(configPath)
	initWorkPoll()
	ln, err := net.ListenTCP("tcp", &net.TCPAddr{
		Port: config.GetGatewayPort(),
	})
	if err != nil {
		panic(err)
	}
	initEPoll(ln, runProc)
	select {}
}
func runProc(c *connection, ep *epoller) {
	ctx := context.Background()
	// 读取数据
	dataBuf, err := tcp.ReadData(c.conn)
	if err != nil {
		if errors.Is(err, io.EOF) {
			_ = ep.remove(c)
			_ = client.CancelConn(&ctx, getEndpoint(), c.id, nil)
		}
		return
	}
	err = wPool.Submit(func() {
		// step2:交给 state server rpc 处理
		_ = client.SendMsg(&ctx, getEndpoint(), c.id, dataBuf)
	})
	if err != nil {
		fmt.Errorf("runProc:err:%+v\n", err.Error())
	}
}
