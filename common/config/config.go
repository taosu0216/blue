package config

import (
	"github.com/spf13/viper"
	"time"
)

func Init(path string) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

// GetEndpointsForDiscovery 获取服务发现的地址
func GetEndpointsForDiscovery() []string {
	return viper.GetStringSlice("discovery.endpoints")
}

// GetTimeoutForDiscovery 获取连接服务发现集群的超时时间 单位是秒
func GetTimeoutForDiscovery() time.Duration {
	return viper.GetDuration("discovery.timeout") * time.Second
}

// GetServicePathForIPConf 获取服务发现的服务路径(就是前缀prefix)
func GetServicePathForIPConf() string {
	return viper.GetString("ip_conf.service_path")
}

func GetGatewayWorkerPoolNum() int {
	return viper.GetInt("gateway.worker_pool_num")
}

func GetGatewayPort() int {
	return viper.GetInt("gateway.tcp_server_port")
}

func GetGatewayEpollerChanNum() int {
	return viper.GetInt("gateway.epoll_channel_size")
}

func GetGatewayEpollerNum() int {
	return viper.GetInt("gateway.epoll_num")
}

func GetGatewayEpollWaitQueueSize() int {
	return viper.GetInt("gateway.epoll_wait_queue_size")
}

func GetGatewayMaxTcpNum() int32 {
	return viper.GetInt32("gateway.tcp_max_num")
}

func GetGatewayServiceAddr() string {
	return viper.GetString("gateway.service_addr")
}

func GetGateWayServiceName() string {
	return viper.GetString("gateway.service_name")
}

func GetGatewayWeight() int {
	return viper.GetInt("gateway.weight")
}

func GetGatewayRPCServerPort() int {
	return viper.GetInt("gateway.rpc_server_port")
}

func GetGatewayCmdChannelNum() int {
	return viper.GetInt("gateway.cmd_channel_num")
}

func GetStateCmdChannelNum() int {
	return viper.GetInt("state.cmd_channel_num")
}

func GetStateServiceName() string {
	return viper.GetString("state.service_name")
}

func GetStateServiceAddr() string {
	return viper.GetString("state.service_addr")
}

func GetStateServerPort() int {
	return viper.GetInt("state.server_port")
}
func GetStateRPCWeight() int {
	return viper.GetInt("state.weight")
}

func GetDomainSocketPushPath() string { return viper.GetString("domain.pushpath") }

func GetDomainSocketPullPath() string { return viper.GetString("domain.pullpath") }
