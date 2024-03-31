package plugin

import (
	"blue/common/brpc/config"
	"blue/common/brpc/discov"
	"blue/common/brpc/discov/etcd"
	"errors"
	"fmt"
)

func GetDiscovInstance() (discov.Discovery, error) {
	// 目前只有etcd
	name := config.GetDiscovName()
	switch name {
	case "etcd":
		// localhost:2379
		return etcd.NewETCDRegister(etcd.WithEndpoints(config.GetDiscovEndpoints()))
	}
	return nil, errors.New(fmt.Sprintf("not exist plugin:%s", name))
}
