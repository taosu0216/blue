package brpc

import (
	"blue/common/brpc/discov"
	"blue/common/brpc/discov/plugin"
	"context"
	"fmt"
	"github.com/bytedance/gopkg/util/logger"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type RegisterFn func(*grpc.Server)

type BServer struct {
	serverOptions
	registers []RegisterFn
	// 拦截器,可用于链路追踪,监控,熔断等功能
	interceptors []grpc.UnaryServerInterceptor
}

type serverOptions struct {
	serviceName string
	ip          string
	port        int
	// 权重,就是撤回消息的信令权重高于发送信令的权重
	weight int
	health bool
	d      discov.Discovery
}

type ServerOption func(opts *serverOptions)

func NewBServer(opts ...ServerOption) *BServer {
	optObject := serverOptions{}
	for _, function := range opts {
		function(&optObject)
	}

	if optObject.d == nil {
		dis, err := plugin.GetDiscovInstance()
		if err != nil {
			panic(err)
		}
		optObject.d = dis
	}

	return &BServer{
		optObject,
		make([]RegisterFn, 0),
		make([]grpc.UnaryServerInterceptor, 0),
	}
}

func WithServiceName(serviceName string) ServerOption {
	return func(opts *serverOptions) {
		opts.serviceName = serviceName
	}
}

func WithIP(ip string) ServerOption {
	return func(opts *serverOptions) {
		opts.ip = ip
	}
}

// WithPort set port
func WithPort(port int) ServerOption {
	return func(opts *serverOptions) {
		opts.port = port
	}
}

// WithWeight set weight
func WithWeight(weight int) ServerOption {
	return func(opts *serverOptions) {
		opts.weight = weight
	}
}

// eg :
//
//	b.RegisterService(func(server *grpc.Server) {
//	    test.RegisterGreeterServer(server, &Server{})
//	})

func (b *BServer) RegisterService(register ...RegisterFn) {
	b.registers = append(b.registers, register...)
}

// Start 开启server
func (b *BServer) Start(ctx context.Context) {
	service := discov.Service{
		Name: b.serviceName,
		Endpoints: []*discov.Endpoint{
			{
				ServerName: b.serviceName,
				IP:         b.ip,
				Port:       b.port,
				Weight:     b.weight,
				Enable:     true,
			},
		},
	}

	// 加载中间件
	//interceptors := []grpc.UnaryServerInterceptor{
	//	serverinterceptor.RecoveryUnaryServerInterceptor(),
	//	serverinterceptor.TraceUnaryServerInterceptor(),
	//	serverinterceptor.MetricUnaryServerInterceptor(b.serviceName),
	//}
	//interceptors = append(interceptors, b.interceptors...)

	// TODO：s := grpc.NewServer(grpc.ChainUnaryInterceptor(interceptors...))

	s := grpc.NewServer()

	// 注册服务
	for _, register := range b.registers {
		register(s)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", b.ip, b.port))
	if err != nil {
		panic(err)
	}

	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
	// 服务注册
	b.d.Register(ctx, &service)

	logger.Info("start PRCP success")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-c
		switch sig {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			s.Stop()
			b.d.UnRegister(ctx, &service)
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}

}
