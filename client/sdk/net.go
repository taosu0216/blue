package sdk

import (
	"blue/common/tcp"
	"encoding/json"
	"fmt"
	"net"
)

type connect struct {
	sendChan, recvChan chan *Message
	conn               *net.TCPConn
	connID             uint64
	ip                 net.IP
	port               int
}

func newConnect(ip net.IP, port int) *connect {

	clientConn := &connect{
		sendChan: make(chan *Message),
		recvChan: make(chan *Message),
		ip:       ip,
		port:     port,
	}
	addr := &net.TCPAddr{IP: ip, Port: port}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Printf("DialTCP.err=%+v\n", err)
		return nil
	}
	clientConn.conn = conn
	go func() {
		for {
			data, _ := tcp.ReadData(conn)
			//if err != nil {
			//	fmt.Printf("ReadData.err=%+v\n", err)
			//}
			msg := &Message{}
			_ = json.Unmarshal(data, msg)
			clientConn.recvChan <- msg
		}
	}()
	return clientConn
}

func (c *connect) send(data *Message) {

	bytes, _ := json.Marshal(data)
	dataPack := &tcp.DataPack{Data: bytes, Len: uint32(len(bytes))}
	msg := dataPack.Marshal()
	_, _ = c.conn.Write(msg)
	//c.sendChan <- data
}

func (c *connect) recv() <-chan *Message {
	if c.recvChan == nil {
		fmt.Println("recvChan is nil")
	}
	return c.recvChan
}

func (c *connect) close() {}
