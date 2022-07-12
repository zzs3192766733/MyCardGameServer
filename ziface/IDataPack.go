package ziface

type IDataPack interface {
	GetHeadLength() uint
	Pack(message IMessage) ([]byte, error)
	UnPack(data []byte) (IMessage, error)
}
