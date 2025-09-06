package tlv_test

import (
	"testing"
	"tlv-protocol/tlv"
)

func TestTLVMessage(t *testing.T) {
	t.Run("NumberMessage", func(t *testing.T) {
		numberMsg := tlv.NewNumberMessage(12345)
		number, err := numberMsg.GetNumber()
		if err != nil {
			t.Fatalf("unexpected error getting number: %v", err)
		}

		if number != 12345 {
			t.Errorf("expected 12345 %d", number)
		}
	})

	t.Run("StringMessage", func(t *testing.T) {
		stringMsg := tlv.NewStringMessage("Hello world")

		text, err := stringMsg.GetString()
		if err != nil {
			t.Fatalf("unexpected error getting string: %v", err)
		}
		if text != "Hello world" {
			t.Errorf("expected 'Hello world' not '%s'", text)
		}
	})

	t.Run("FileMessage", func(t *testing.T) {
		fileData := []byte("This is file contents Hello world")
		fileMsg := tlv.NewFileMessage(fileData)

		retrieved, err := fileMsg.GetFileData()
		if err != nil {
			t.Fatalf("unexpected error getting file data: %v", err)
		}
		if string(retrieved) != string(fileData) {
			t.Errorf("expected '%s' not '%s'", string(fileData), string(retrieved))
		}
	})

	t.Run("TypeValidation", func(t *testing.T) {
		stringMsg := tlv.NewStringMessage("Test")
		_, err := stringMsg.GetNumber()
		if err == nil {
			t.Errorf("expected error when getting number from string message but got nil")
		}
	})
}
