package tlv_test

import (
	"testing"
	"tlv-protocol/tlv"
)

func TestTlvMessagesEncoding(t *testing.T) {
	t.Run("NumberMessageEncode", func(t *testing.T) {
		originalNumber := int64(12345)
		msg := tlv.NewNumberMessage(originalNumber)
		number, err := msg.GetNumber()
		if err != nil {
			t.Fatalf("Get Number error %v", err)
		}
		if number != originalNumber {
			t.Errorf("expected %d actual %d", originalNumber, number)
		}

		buffer, err := msg.Encode()
		if err != nil {
			t.Fatalf("Encode error %v", err)
		}
		if len(buffer) != 1+4+8 {
			t.Errorf("Encode length %d actual %d", len(buffer), 1+4+8)
		}
	})

	t.Run("StringMessageEncode", func(t *testing.T) {
		text := "Hello world"
		msg := tlv.NewStringMessage(text)

		got, err := msg.GetString()
		if err != nil {
			t.Fatalf("Get string error = %v", err)
		}
		if got != text {
			t.Errorf("Get string = %q; want %q", got, text)
		}

		buf, err := msg.Encode()
		if err != nil {
			t.Fatalf("Encdoe error = %v", err)
		}
		if len(buf) != 1+4+len(text) {
			t.Errorf("Encode length = %d; actual %d", len(buf), 1+4+len(text))
		}
	})

	t.Run("FileMessageEncode", func(t *testing.T) {
		fileData := []byte("This is file content hello world")
		msg := tlv.NewFileMessage(fileData)

		got, err := msg.GetFileData()
		if err != nil {
			t.Fatalf("GetFileData error = %v", err)
		}
		if string(got) != string(fileData) {
			t.Errorf("GetFileData = %q; actual %q", string(got), string(fileData))
		}

		buf, err := msg.Encode()
		if err != nil {
			t.Fatalf("Encode error = %v", err)
		}
		if len(buf) != 1+4+len(fileData) {
			t.Errorf("Encode length = %d; actual %d", len(buf), 1+4+len(fileData))
		}
	})

}
