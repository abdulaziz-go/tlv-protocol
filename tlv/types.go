package tlv

import "fmt"

const (
	TypeNumber = 0x01
	TypeString = 0x02
	TypeFile   = 0x03
)

type TLVMessage struct {
	Type   uint8
	Length uint32
	Value  []byte
}

func NewNumberMessage(number int64) *TLVMessage {
	value := make([]byte, 8)

	value[0] = byte(number)
	value[1] = byte(number >> 8)
	value[2] = byte(number >> 16)
	value[3] = byte(number >> 24)
	value[4] = byte(number >> 32)
	value[5] = byte(number >> 40)
	value[6] = byte(number >> 48)
	value[7] = byte(number >> 56)

	return &TLVMessage{
		Type:   TypeNumber,
		Length: 8,
		Value:  value,
	}
}

func NewStringMessage(txt string) *TLVMessage {
	value := []byte(txt)
	return &TLVMessage{
		Type:   TypeString,
		Length: uint32(len(value)),
		Value:  value,
	}
}

func NewFileMessage(fileData []byte) *TLVMessage {
	return &TLVMessage{
		Type:   TypeFile,
		Length: uint32(len(fileData)),
		Value:  fileData,
	}
}

func (msg *TLVMessage) GetNumber() (int64, error) {
	if msg.Type != TypeNumber {
		return 0, fmt.Errorf("message type is not number %d", msg.Type)
	}

	if len(msg.Value) != 8 {
		return 0, fmt.Errorf("invalid number length expect 8 but %d", len(msg.Value))
	}

	number := int64(msg.Value[0]) |
		int64(msg.Value[1])<<8 |
		int64(msg.Value[2])<<16 |
		int64(msg.Value[3])<<24 |
		int64(msg.Value[4])<<32 |
		int64(msg.Value[5])<<40 |
		int64(msg.Value[6])<<48 |
		int64(msg.Value[7])<<56

	return number, nil
}

func (msg *TLVMessage) GetString() (string, error) {
	if msg.Type != TypeString {
		return "", fmt.Errorf("message type is not string %d", msg.Type)
	}

	return string(msg.Value), nil
}

func (msg *TLVMessage) GetFileData() ([]byte, error) {
	if msg.Type != TypeFile {
		return nil, fmt.Errorf("message type is not file %d", msg.Type)
	}

	return msg.Value, nil
}

func (msg *TLVMessage) String() string {
	var typeName string
	switch msg.Type {
	case TypeNumber:
		typeName = "NUMBER"
	case TypeString:
		typeName = "STRING"
	case TypeFile:
		typeName = "FILE"
	default:
		typeName = fmt.Sprintf("UNKNOWN(%d)", msg.Type)
	}

	return fmt.Sprintf("TLV{Type:%s, Length:%d, ValueSize:%d}",
		typeName, msg.Length, len(msg.Value))
}
