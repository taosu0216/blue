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

func GetGatewayRPCServerPort() int {
	return viper.GetInt("gateway.rpc_server_port")
}
