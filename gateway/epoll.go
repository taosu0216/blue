package gateway

import (
	"blue/common/config"
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"log"
	"net"
	"runtime"
	"sync"
	"syscall"
)

var (
	ep     *ePoll
	tcpNum int32
)

type ePoll struct {
	eChan  chan *connection
	tables sync.Map
	eSize  int
	done   chan struct{}

	//长连接
	ln *net.TCPListener
	//回调函数
	f func(c *connection, ep *epoller)
}

// epoller 对象 轮询器
type epoller struct {
	fd            int
	fdToConnTable sync.Map
}

func initEPoll(ln *net.TCPListener, f func(c *connection, ep *epoller)) {
	setLimit()
	ep = newEPoll(ln, f)
	ep.createAcceptProcess()
	ep.startEPoll()
}

func newEPoll(ln *net.TCPListener, f func(c *connection, ep *epoller)) *ePoll {
	return &ePoll{
		eChan:  make(chan *connection, config.GetGatewayEpollerChanNum()),
		ln:     ln,
		f:      f,
		done:   make(chan struct{}),
		tables: sync.Map{},
		//4个epoll对象
		eSize: config.GetGatewayEpollerNum(),
	}
}

// 设置go 进程打开文件数的限制
// 默认fd最大只有1024,修改成最大值
func setLimit() {
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	log.Println("set cur limit: ", rLimit.Cur)
}

func (e *ePoll) createAcceptProcess() {
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				conn, err := e.ln.AcceptTCP()
				//手动设置tcp为长连接模式
				setTcpConfig(conn)
				//获取到连接并且进行熔断限流
				if !checkTcp() {
					_ = conn.Close()
					continue
				}
				if err != nil {
					var netErr net.Error
					//如果是网络抽风的临时错误直接跳过就可以
					if errors.As(err, &netErr) && netErr.Temporary() {
						_ = fmt.Errorf("accept temp err: %v", netErr)
						continue
					}
					_ = fmt.Errorf("accept err: %v", err)
				}
				c := newConnection(conn)
				ep.addTask(c)
			}
		}()
	}
}

func (e *ePoll) startEPoll() {
	for i := 0; i < e.eSize; i++ {
		go e.startEProc()
	}
}

func (e *ePoll) addTask(c *connection) {
	fmt.Println("new tcp connection created,fd: ", c.fd, " tcp num: ", tcpNum)
	e.eChan <- c
}

func newEPoller() (*epoller, error) {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}
	return &epoller{
		fd: fd,
	}, nil
}

func (e *ePoll) startEProc() {
	ep, err := newEPoller()
	if err != nil {
		panic(err)
	}
	// 监听连接创建事件
	go func() {
		for {
			select {
			case <-e.done:
				return
			case conn := <-e.eChan:
				addTcpNum()
				fmt.Println("tcpNum: ", tcpNum)
				if err := ep.add(conn); err != nil {
					fmt.Printf("failed to add connection %v\n", err)
					conn.Close() //登录未成功直接关闭连接
					continue
				}
			}
		}
	}()
	for {
		select {
		case <-e.done:
			return
		default:
			connections, err := ep.wait(200) // 200ms 一次轮询避免 忙轮询
			if err != nil && !errors.Is(err, syscall.EINTR) {
				fmt.Printf("failed to epoll wait %v\n", err)
				continue
			}
			for _, conn := range connections {
				if conn == nil {
					break
				}
				e.f(conn, ep)
			}
		}
	}
}

func (e *epoller) add(c *connection) error {
	fd := c.fd
	unixEpoll := &unix.EpollEvent{
		Events: unix.EPOLLIN | unix.EPOLLHUP,
		Fd:     int32(fd),
	}
	//e.fd的fd就是调用系统底层的epoll_create1创建的epoll对象的fd,这里是把当前链接放入epoll中进行监听
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, unixEpoll)
	if err != nil {
		return err
	}
	e.fdToConnTable.Store(c.fd, c)
	// TODO: id问题
	ep.tables.Store(c.id, c)
	ep.tables.Store(fd, c)
	c.e = e
	return nil
}

func (e *epoller) remove(c *connection) error {
	subTcpNum()
	fd := c.fd
	err := unix.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	ep.tables.Delete(c.id)
	e.fdToConnTable.Delete(c.fd)
	return nil
}

func (e *epoller) wait(msec int) ([]*connection, error) {
	// 最多100
	events := make([]unix.EpollEvent, config.GetGatewayEpollWaitQueueSize())
	// 监听系统级epoll对象的fd,等待事件发生
	n, err := unix.EpollWait(e.fd, events, msec)
	if err != nil {
		return nil, err
	}
	connections := make([]*connection, 0, n)
	// var connections []*connection
	for i := 0; i < n; i++ {
		if conn, ok := e.fdToConnTable.Load(int(events[i].Fd)); ok {
			connections = append(connections, conn.(*connection))
		}
	}
	return connections, nil
}
