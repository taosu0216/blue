package perf

import (
	"blue/client/sdk"
	"net"
)

var (
	TcpConnNum int32
)

func RunMain() {
	for i := 0; i < int(TcpConnNum); i++ {
		sdk.NewChat(net.ParseIP("172.31.182.205"), 8900, "logic", "1223", "123")
	}
}
