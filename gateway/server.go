package gateway

import (
	"blue/common/config"
	"blue/common/tcp"
	"errors"
	"fmt"
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
	fmt.Println("-------------im gateway stated------------")
	select {}
}

func runProc(c *connection, ep *epoller) {
	//ctx := context.Background()
	// 读取数据
	dataBuf, err := tcp.ReadData(c.conn)
	if err != nil {
		if errors.Is(err, io.EOF) {
			_ = ep.remove(c)
			//_ = client.CancelConn(&ctx, getEndpoint(), c.id, nil)
		}
		return
	}
	_ = wPool.Submit(func() {
		bytes := &tcp.DataPack{
			Len:  uint32(len(dataBuf)),
			Data: dataBuf,
		}
		_ = tcp.SendData(c.conn, bytes.Marshal())
	})
}
