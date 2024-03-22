package ipconf

import (
	"blue/common/config"
	"blue/ipconf/source"
	"github.com/cloudwego/hertz/pkg/app/server"
)

// RunMain 启动web服务
func RunMain(path string) {
	config.Init(path)
	source.Init()
	s := server.Default(server.WithHostPorts(":6789"))
	s.GET("/ip/list", GetIpInfoList)
	s.Spin()
}
