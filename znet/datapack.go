package znet

import (
	"MyGameServer/ziface"
	"bytes"
	"encoding/binary"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLength() uint {
	return 8
}

func (d *DataPack) Pack(message ziface.IMessage) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	var err error
	if err = binary.Write(buffer, binary.LittleEndian, message.GetMsgLen()); err != nil {
		return nil, err
	}
	if err = binary.Write(buffer, binary.LittleEndian, message.GetMsgID()); err != nil {
		return nil, err
	}
	if err = binary.Write(buffer, binary.LittleEndian, message.GetData()); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (d *DataPack) UnPack(data []byte) (ziface.IMessage, error) {
	buffer := bytes.NewReader(data)
	message := &Message{}
	var err error
	if err = binary.Read(buffer, binary.LittleEndian, &message.Length); err != nil {
		return nil, err
	}
	if err = binary.Read(buffer, binary.LittleEndian, &message.ID); err != nil {
		return nil, err
	}
	if err = binary.Read(buffer, binary.LittleEndian, &message.Data); err != nil {
		return nil, err
	}
	return message, nil
}
