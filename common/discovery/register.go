package discovery

import (
	"blue/common/config"
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

// 向etcd注册服务的流程

type ServiceRegister struct {
	cli           *clientv3.Client
	lease         clientv3.LeaseID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse
	key           string
	val           string //value
	ctx           *context.Context
}

// NewServiceRegister 新建服务
func NewServiceRegister(ctx *context.Context, key string, endportinfo *EndpointInfo, lease int64) (*ServiceRegister, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEndpointsForDiscovery(),
		DialTimeout: config.GetTimeoutForDiscovery(),
	})
	if err != nil {
		log.Fatalln(err)
	}
	ser := &ServiceRegister{
		cli: cli,
		val: endportinfo.Marshal(),
		key: key,
		ctx: ctx,
	}

	//申请租约设置时间keepalive
	if err := ser.putKeyWithLease(lease); err != nil {
		return nil, err
	}
	return ser, nil
}

func (s *ServiceRegister) putKeyWithLease(lease int64) error {
	//申请一个租约
	resp, err := s.cli.Grant(*s.ctx, lease)
	if err != nil {
		return err
	}
	//设置key
	_, err = s.cli.Put(*s.ctx, s.key, s.val, clientv3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	//设置续租
	keepAliveChan, err := s.cli.KeepAlive(*s.ctx, resp.ID)
	if err != nil {
		return err
	}
	s.lease = resp.ID
	s.keepAliveChan = keepAliveChan
	return nil
}
