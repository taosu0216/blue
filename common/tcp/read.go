package tcp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

func ReadData(conn *net.TCPConn) ([]byte, error) {
	//存储从缓冲区中读取的数据
	var dataLen uint32

	//缓冲区从dataLenBuf读出数据
	dataLenBuf := make([]byte, 4)

	//第一次读的数据是数据头的长度,要从中读取数据的长度
	if err := readFixedData(conn, dataLenBuf); err != nil {
		return nil, err
	}

	//把数据头中的数据存入缓冲区,从缓冲区里面读出数据的长度并赋值给dataLen
	buffer := bytes.NewBuffer(dataLenBuf)
	if err := binary.Read(buffer, binary.BigEndian, &dataLen); err != nil {
		return nil, fmt.Errorf("read headlen error:%s", err.Error())
	}
	if dataLen <= 0 {
		return nil, fmt.Errorf("wrong headlen :%d", dataLen)
	}

	//根据数据长度读取数据
	dataBuf := make([]byte, dataLen)
	if err := readFixedData(conn, dataBuf); err != nil {
		return nil, fmt.Errorf("read headlen error:%s", err.Error())
	}
	return dataBuf, nil
}

func readFixedData(conn *net.TCPConn, buf []byte) error {
	// 设置读取操作的超时时间为120s
	_ = (*conn).SetReadDeadline(time.Now().Add(time.Duration(120) * time.Second))
	var pos int = 0
	var totalSize int = len(buf)
	for {
		index, err := (*conn).Read(buf[pos:])
		if err != nil {
			return err
		}
		pos = pos + index
		if pos == totalSize {
			break
		}
	}
	return nil
}
