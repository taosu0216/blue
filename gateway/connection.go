package gateway

import (
	"errors"
	"net"
	"reflect"
	"sync"
	"time"
)

var node = &ConnIDGenerater{}

const (
	version      = uint64(0) // 版本控制
	sequenceBits = uint64(16)

	maxSequence = int64(-1) ^ (int64(-1) << sequenceBits)

	timeLeft    = uint8(16) // timeLeft = sequenceBits // 时间戳向左偏移量
	versionLeft = uint8(63) // 左移动到最高位
	// 2020-05-20 08:00:00 +0800 CST
	twepoch = int64(1589923200000) // 常量时间戳(毫秒)
)

type connection struct {
	id   uint64 // 进程级别的生命周期
	fd   int
	e    *epoller
	conn *net.TCPConn
}

type ConnIDGenerater struct {
	mu        sync.Mutex
	LastStamp int64 // 记录上一次ID的时间戳
	Sequence  int64 // 当前毫秒已经生成的ID序列号(从0 开始累加) 1毫秒内最多生成2^16个ID
}

func newConnection(conn *net.TCPConn) *connection {
	var id uint64
	var err error
	if id, err = node.NextID(); err != nil {
		panic(err) // 在线服务需要解决这个问题 ，报错而不能panic
	}
	return &connection{
		id:   id,
		fd:   socketFD(conn),
		conn: conn,
	}
}

func (gen *ConnIDGenerater) NextID() (uint64, error) {
	gen.mu.Lock()
	defer gen.mu.Unlock()
	return gen.nextID()
}

func (gen *ConnIDGenerater) nextID() (uint64, error) {
	timeStamp := gen.getMilliSeconds()
	// 避免时钟回拨
	if timeStamp < gen.LastStamp {
		return 0, errors.New("time is moving backwards,waiting until")
	}
	// 如果是同一时间生成的，则进行毫秒内序列
	if gen.LastStamp == timeStamp {
		gen.Sequence = (gen.Sequence + 1) & maxSequence
		if gen.Sequence == 0 { // 如果这里发生溢出，就等到下一个毫秒时再分配，这样就一定出现重复
			for timeStamp <= gen.LastStamp {
				timeStamp = gen.getMilliSeconds()
			}
		}
	} else { // 如果与上次分配的时间戳不等，则为了防止可能的时钟飘移现象，就必须重新计数
		gen.Sequence = 0
	}
	gen.LastStamp = timeStamp
	// 减法可以压缩一下时间戳
	id := ((timeStamp - twepoch) << timeLeft) | gen.Sequence
	connID := uint64(id) | (version << versionLeft)
	return connID, nil
}

func (gen *ConnIDGenerater) getMilliSeconds() int64 {
	return time.Now().UnixNano() / 1e6
}

// 获取socket fd
/*
	type TCPConn struct {
		conn
	}

	type conn struct {
		fd *netFD
	}

	type netFD struct {
		pfd poll.FD

		// immutable until Close
		family      int
		sotype      int
		isConnected bool // handshake completed or use of association with peer
		net         string
		laddr       Addr
		raddr       Addr
	}

	type FD struct {
		// Lock sysfd and serialize access to Read and Write methods.
		fdmu fdMutex

		// System file descriptor. Immutable until Close.
		Sysfd int

		// I/O poller.
		pd pollDesc

		// Writev cache.
		iovecs *[]syscall.Iovec

		// Semaphore signaled when file is closed.
		csema uint32

		// Non-zero if this file has been set to blocking mode.
		isBlocking uint32

		// Whether this is a streaming descriptor, as opposed to a
		// packet-based descriptor like a UDP socket. Immutable.
		IsStream bool

		// Whether a zero byte read indicates EOF. This is false for a
		// message based socket connection.
		ZeroReadIsEOF bool

		// Whether this is a file rather than a network socket.
		isFile bool
	}
*/
func socketFD(conn *net.TCPConn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(*conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	//Sysfd才是linux系统常说的fd
	return int(pfdVal.FieldByName("Sysfd").Int())
}

func (c *connection) Close() {
	ep.tables.Delete(c.id)
	if c.e != nil {
		c.e.fdToConnTable.Delete(c.fd)
	}
	err := c.conn.Close()
	panic(err)
}
