package source

import (
	"blue/common/discovery"
	"fmt"
)

type EventType string

type Event struct {
	Type         EventType
	IP           string
	Port         string
	ConnectNum   float64
	MessageBytes float64
}

var eventChan chan *Event

const (
	AddNodeEvent EventType = "addNode"
	DelNodeEvent EventType = "delNode"
)

func EventChan() <-chan *Event {
	return eventChan
}

func NewEvent(ed *discovery.EndpointInfo) *Event {
	if ed == nil || ed.MetaData == nil {
		return nil
	}
	var connNum, msgBytes float64
	if data, ok := ed.MetaData["connect_num"]; ok {
		connNum = data.(float64)
	}
	if data, ok := ed.MetaData["message_bytes"]; ok {
		msgBytes = data.(float64)
	}
	return &Event{
		IP:           ed.IP,
		Port:         ed.Port,
		ConnectNum:   connNum,
		MessageBytes: msgBytes,
	}
}

func (e *Event) Key() string {
	return fmt.Sprintf("%s:%s", e.IP, e.Port)
}
