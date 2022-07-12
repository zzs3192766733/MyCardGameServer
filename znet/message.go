package znet

type Message struct {
	Length uint32
	ID     uint32
	Data   []byte
}

func (m *Message) GetMsgID() uint32 {
	return m.ID
}

func (m *Message) GetMsgLen() uint32 {
	return m.Length
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetMsgID(id uint32) {
	m.ID = id
}

func (m *Message) SetMsgLen(len uint32) {
	m.Length = len
}

func (m *Message) SetData(bytes []byte) {
	m.Data = bytes
}

func NewMessage(ID uint32, data []byte) *Message {
	return &Message{Length: uint32(len(data)), ID: ID, Data: data}
}
