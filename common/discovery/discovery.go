package discovery

import (
	"blue/common/config"
	"context"
	"github.com/bytedance/gopkg/util/logger"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"sync"
)

// 从etcd获得服务信息

type ServiceDiscovery struct {
	cli  *clientv3.Client //etcd client
	lock sync.Mutex
	ctx  *context.Context
}

func NewServiceDiscovery(ctx *context.Context) *ServiceDiscovery {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.GetEndpointsForDiscovery(),
		DialTimeout: config.GetTimeoutForDiscovery(),
	})
	if err != nil {
		logger.Fatal(err)
	}
	return &ServiceDiscovery{
		cli: cli,
		ctx: ctx,
	}
}

func (s *ServiceDiscovery) Close() error {
	return s.cli.Close()
}

// WatchService 初始化服务列表和监视
func (s *ServiceDiscovery) WatchService(prefix string, set, del func(key, value string)) error {
	//根据前缀获取现有的key
	resp, err := s.cli.Get(*s.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}
	for _, ev := range resp.Kvs {
		set(string(ev.Key), string(ev.Value))
	}
	s.watcher(prefix, resp.Header.Revision+1, set, del)
	return nil
}

// watcher 监听前缀
// 这里有一个for循环,监听channel进行阻塞,所以初始化的时候给开了一个携程,每当数据有更新,这里就会一直监听
func (s *ServiceDiscovery) watcher(prefix string, revision int64, set, del func(key, value string)) {
	//rch是一个channel
	rch := s.cli.Watch(*s.ctx, prefix, clientv3.WithPrefix(), clientv3.WithRev(revision))
	logger.CtxInfof(*s.ctx, "watching prefix:%s now...", prefix)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			switch ev.Type {
			case mvccpb.PUT: //修改或者新增
				set(string(ev.Kv.Key), string(ev.Kv.Value))
			case mvccpb.DELETE: //删除
				del(string(ev.Kv.Key), string(ev.Kv.Value))
			}
		}
	}
}
