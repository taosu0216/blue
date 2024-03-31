package tcp

import (
	"net"
)

// 将指定数据包发送到网络中
func SendData(conn *net.TCPConn, data []byte) error {
	totalLen := len(data)
	writeLen := 0
	for {
		l, err := conn.Write(data[writeLen:])
		if err != nil {
			return err
		}
		writeLen = writeLen + l
		if writeLen >= totalLen {
			break
		}
	}
	return nil
}
