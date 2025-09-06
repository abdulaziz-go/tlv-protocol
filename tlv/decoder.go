package tlv

import "fmt"

func Decode(data []byte) (*TLVMessage, error) {
	if len(data) < 5 {
		return nil, fmt.Errorf("data too short expected 5 bytes not %d", len(data))
	}

	msgType := data[0]
	length := uint32(data[1]) | uint32(data[2])<<8 | uint32(data[3])<<16 | uint32(data[4])<<24

	totalNeeded := 5 + int(length)

	if len(data) < totalNeeded {
		return nil, fmt.Errorf("insufficient data need %d bytes not %d", totalNeeded, len(data))
	}

	value := make([]byte, length)
	copy(value, data[5:5+length])

	msg := &TLVMessage{
		Type:   msgType,
		Length: length,
		Value:  value,
	}

	if err := validateDecodedType(msg); err != nil {
		return nil, fmt.Errorf("invalid decoded message %v", err)
	}

	return msg, nil
}

func validateDecodedType(msg *TLVMessage) error {
	switch msg.Type {
	case TypeNumber:
		if msg.Length != 8 {
			return fmt.Errorf("number type must have length 8 not %d", msg.Length)
		}
	case TypeString:
		if msg.Length == 0 {
			return fmt.Errorf("string content cannot be empty")
		}
	case TypeFile:
		if msg.Length == 0 {
			return fmt.Errorf("file type cannot be empty")
		}
	default:
		return fmt.Errorf("unknown message type: %d", msg.Type)
	}
	return nil
}

func DecodeMultiple(data []byte) ([]*TLVMessage, error) {
	var messages []*TLVMessage
	offset := 0

	for offset < len(data) {
		if offset+5 > len(data) {
			return nil, fmt.Errorf("incomplete TLV at offset %d need 4 bytes for header not %d",
				offset, len(data)-offset)
		}
		length := uint32(data[offset+1]) | uint32(data[offset+2])<<8 | uint32(data[offset+3])<<16 | uint32(data[offset+4])<<24
		totalSize := 5 + int(length)

		if offset+totalSize > len(data) {
			return nil, fmt.Errorf("incomplete TLV at offset %d: need %d bytes not %d",
				offset, totalSize, len(data)-offset)
		}

		msg, err := Decode(data[offset : offset+totalSize])
		if err != nil {
			return nil, fmt.Errorf("failed to decode message at offset %d: %v", offset, err)
		}

		messages = append(messages, msg)
		offset += totalSize
	}
	return messages, nil
}
