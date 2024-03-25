package tcp

import (
	"bytes"
	"encoding/binary"
)

type DataPack struct {
	Len  uint32
	Data []byte
}

func (d *DataPack) Marshal() []byte {
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, d.Len)
	return append(bytesBuffer.Bytes(), d.Data...)
}
