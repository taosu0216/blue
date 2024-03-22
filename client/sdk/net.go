package sdk

type conncet struct {
	serverAddr         string
	SendChan, RecvChan chan *Message
}

func newConnect(addr string) *conncet {
	return &conncet{
		serverAddr: addr,
		//TODO: 为什么这里要有缓冲??????
		SendChan: make(chan *Message, 10),
		RecvChan: make(chan *Message, 10),
	}
}

func (c *conncet) send(data *Message) {
	c.SendChan <- data
}

func (c *conncet) recv() <-chan *Message {
	return c.RecvChan
}

func (c *conncet) close() {}
