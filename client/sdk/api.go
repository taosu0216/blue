package sdk

const (
	MsgTypeText = "text"
	MsgTypeAck  = "ack"
)

type Chat struct {
	Nick      string
	UserID    string
	SessionID string
	Conn      *conncet
}

type Message struct {
	Type       string
	Name       string
	FromUserID string
	ToUserID   string
	Content    string
	Session    string
}

func NewChat(addr, nick, userID, sessionID string) *Chat {
	return &Chat{
		Nick:      nick,
		UserID:    userID,
		SessionID: sessionID,
		Conn:      newConnect(addr),
	}
}

func (chat *Chat) Send(msg *Message) {
	chat.Conn.send(msg)
}

func (chat *Chat) Close() {
	chat.Conn.close()
}

func (chat *Chat) Recv() <-chan *Message {
	return chat.Conn.recv()
}
