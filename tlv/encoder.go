package tlv

import (
	"bytes"
	"fmt"
)

func (msg *TLVMessage) Encode() ([]byte, error) {
	if len(msg.Value) != int(msg.Length) {
		return nil, fmt.Errorf("length mismatch: length %d , value %d", msg.Length, len(msg.Value))
	}

	bufferSize := msg.GetEncodedSize()
	buffer := make([]byte, bufferSize)

	buffer[0] = msg.Type

	buffer[1] = byte(msg.Length)
	buffer[2] = byte(msg.Length >> 8)
	buffer[3] = byte(msg.Length >> 16)
	buffer[4] = byte(msg.Length >> 24) // maximum 4 gb value o'tkazish uchun uint32 dan foydalanganda lengthni maximum byte 4 byte bo'ladi

	copy(buffer[5:], msg.Value)
	return buffer, nil
}

func EncodeMultiple(messages []*TLVMessage) ([]byte, error) {
	if len(messages) == 0 {
		return []byte{}, nil
	}

	var buffer bytes.Buffer

	for i, msg := range messages {
		if msg == nil {
			return nil, fmt.Errorf("message at index %d is nil", i)
		}

		encoded, err := msg.Encode()
		if err != nil {
			return nil, fmt.Errorf("failed to encode message at index %d: %v", i, err)
		}

		buffer.Write(encoded)
	}

	return buffer.Bytes(), nil
}

func (msg *TLVMessage) GetEncodedSize() int {
	return 1 + 4 + int(msg.Length)
}
