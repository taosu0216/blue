package config

import "github.com/spf13/viper"

// GetDiscovName 获取discov用哪种方式实现
func GetDiscovName() string {
	// etcd
	return viper.GetString("brpc.discov.name")
}

func GetDiscovEndpoints() []string {
	// localhost:2379
	return viper.GetStringSlice("discovery.endpoints")
}

func GetGatewayServiceName() string {
	return viper.GetString("gateway.service_name")
}
