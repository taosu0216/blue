package gateway

import (
	"blue/common/config"
	"fmt"
	"github.com/panjf2000/ants"
	"net"
	"sync/atomic"
)

var wPool *ants.Pool

func initWorkPoll() {
	var err error
	// 1024个协程
	if wPool, err = ants.NewPool(config.GetGatewayWorkerPoolNum()); err != nil {
		fmt.Printf("InitWorkPoll.err :%s num:%d\n", err.Error(), config.GetGatewayWorkerPoolNum())
	}
}

func setTcpConfig(c *net.TCPConn) {
	_ = c.SetKeepAlive(true)
}

func checkTcp() bool {
	num := getTcpNum()
	maxTcpNum := config.GetGatewayMaxTcpNum()
	return num <= maxTcpNum
}

// 原子操作读取tcp连接数,保证并发安全
func getTcpNum() int32 {
	return atomic.LoadInt32(&tcpNum)
}

func addTcpNum() {
	atomic.AddInt32(&tcpNum, 1)
}

func subTcpNum() {
	atomic.AddInt32(&tcpNum, -1)
}

func getEndpoint() string {
	return fmt.Sprintf("%s:%d", config.GetGatewayServiceAddr(), config.GetGatewayRPCServerPort())
}
