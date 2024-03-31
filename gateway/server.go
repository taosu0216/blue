package gateway

import (
	"blue/common/brpc"
	"blue/common/config"
	"blue/common/tcp"
	"blue/gateway/rpc/client"
	"blue/gateway/rpc/service"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"net"
)

var cmdChannel chan *service.CmdContext

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

	cmdChannel = make(chan *service.CmdContext, config.GetGatewayCmdChannelNum())

	s := brpc.NewBServer(
		brpc.WithServiceName(config.GetGateWayServiceName()),
		brpc.WithIP(config.GetGatewayServiceAddr()),
		brpc.WithPort(config.GetGatewayRPCServerPort()),
		brpc.WithWeight(config.GetGatewayWeight()),
	)

	s.RegisterService(func(server *grpc.Server) {
		service.RegisterGatewayServer(server, &service.Service{CmdChannel: cmdChannel})
	})

	fmt.Println("-------------im gateway stated------------")
	// 启动rpc 客户端
	client.Init()
	// 启动 命令处理写协程
	go cmdHandler()
	// 启动 rpc server
	s.Start(context.TODO())
}

func runProc(c *connection, ep *epoller) {
	ctx := context.Background()
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
		_ = client.SendMsg(&ctx, getEndpoint(), int32(c.fd), dataBuf)
	})
}

func cmdHandler() {
	for cmd := range cmdChannel {
		// 异步提交到协池中完成发送任务
		switch cmd.Cmd {
		case service.DelConnCmd:
			_ = wPool.Submit(func() { closeConn(cmd) })
		case service.PushCmd:
			_ = wPool.Submit(func() { sendMsgByCmd(cmd) })
		default:
			panic("command undefined")
		}
	}
}
func closeConn(cmd *service.CmdContext) {
	if connPtr, ok := ep.tables.Load(cmd.FD); ok {
		conn, _ := connPtr.(*connection)
		conn.Close()
		ep.tables.Delete(cmd.FD)
	}
}
func sendMsgByCmd(cmd *service.CmdContext) {
	if connPtr, ok := ep.tables.Load(cmd.FD); ok {
		conn, _ := connPtr.(*connection)
		dp := tcp.DataPack{
			Len:  uint32(len(cmd.Payload)),
			Data: cmd.Payload,
		}
		_ = tcp.SendData(conn.conn, dp.Marshal())
	} else {
		fmt.Println("sendMsgByCmd failed")
	}
}
