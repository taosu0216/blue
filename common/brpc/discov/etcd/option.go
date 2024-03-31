package etcd

import "time"

var (
	defaultOption = Options{
		endpoints:              []string{"127.0.0.1:2379"},
		dialTimeout:            10 * time.Second,
		syncFlushCacheInterval: 10 * time.Second,
		keepAliveInterval:      10,
	}
)

type Options struct {
	endpoints                          []string
	keepAliveInterval                  int64
	syncFlushCacheInterval             time.Duration
	dialTimeout                        time.Duration
	registerServiceOrKeepAliveInterval time.Duration
}

type Option func(o *Options)

func WithEndpoints(endpoints []string) Option {
	return func(o *Options) {
		o.endpoints = endpoints
	}
}
