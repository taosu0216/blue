package brpc

import (
	"blue/common/brpc/discov"
	"blue/common/brpc/discov/plugin"
	bresolver "blue/common/brpc/resolver"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"time"
)

type BClient struct {
	serviceName  string
	d            discov.Discovery
	interceptors []grpc.UnaryClientInterceptor
	conn         *grpc.ClientConn
}

const (
	dialTimeout = 5 * time.Second
)

func NewBClient(serviceName string, interceptors ...grpc.UnaryClientInterceptor) (*BClient, error) {
	b := &BClient{
		serviceName:  serviceName,
		interceptors: interceptors,
	}

	if b.d == nil {
		dis, err := plugin.GetDiscovInstance()
		if err != nil {
			panic(err)
		}

		b.d = dis
	}

	resolver.Register(bresolver.NewDiscovBuilder(b.d))

	conn, err := b.dial()
	b.conn = conn

	return b, err
}

func (b *BClient) Conn() *grpc.ClientConn {
	return b.conn
}

func (b *BClient) dial() (*grpc.ClientConn, error) {
	svcCfg := fmt.Sprintf(`{"loadBalancingPolicy":"%s"}`, roundrobin.Name)
	balancerOpt := grpc.WithDefaultServiceConfig(svcCfg)

	//interceptors := []grpc.UnaryClientInterceptor{
	//	clientinterceptor.TraceUnaryClientInterceptor(),
	//	clientinterceptor.MetricUnaryClientInterceptor(),
	//}
	//interceptors = append(interceptors, b.interceptors...)

	options := []grpc.DialOption{
		balancerOpt,
		//grpc.WithChainUnaryInterceptor(interceptors...),
		grpc.WithInsecure(),
	}

	ctx, _ := context.WithTimeout(context.Background(), dialTimeout)

	return grpc.DialContext(ctx, fmt.Sprintf("discov:///%v", b.serviceName), options...)
}
