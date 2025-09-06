package tlv_test

import (
	"testing"
	tlv2 "tlv-protocol/tlv"
)

func TestDecodeFunctions(t *testing.T) {
	t.Run("Decode Number Message", func(t *testing.T) {
		number := int64(12345678)
		msg := tlv2.NewNumberMessage(number)

		buf, err := msg.Encode()
		if err != nil {
			t.Fatalf("Encode error: %v", err)
		}

		decoded, err := tlv2.Decode(buf)
		if err != nil {
			t.Fatalf("Decode error: %v", err)
		}

		got, err := decoded.GetNumber()
		if err != nil {
			t.Fatalf("GetNumber error: %v", err)
		}

		if got != number {
			t.Errorf("decoded number = %d; want %d", got, number)
		}
	})

	t.Run("Decode String Message", func(t *testing.T) {
		text := "Hello world"
		msg := tlv2.NewStringMessage(text)

		buf, err := msg.Encode()
		if err != nil {
			t.Fatalf("Encode error: %v", err)
		}

		decoded, err := tlv2.Decode(buf)
		if err != nil {
			t.Fatalf("Decode error: %v", err)
		}

		got, err := decoded.GetString()
		if err != nil {
			t.Fatalf("GetString error: %v", err)
		}

		if got != text {
			t.Errorf("decoded string = %q want %q", got, text)
		}
	})

	t.Run("Decode File Message", func(t *testing.T) {
		fileData := []byte("File content here hello world")
		msg := tlv2.NewFileMessage(fileData)

		buf, err := msg.Encode()
		if err != nil {
			t.Fatalf("Encode error: %v", err)
		}

		decoded, err := tlv2.Decode(buf)
		if err != nil {
			t.Fatalf("Decode error: %v", err)
		}

		got, err := decoded.GetFileData()
		if err != nil {
			t.Fatalf("GetFileData error: %v", err)
		}

		if string(got) != string(fileData) {
			t.Errorf("decoded file data = %q; want %q", string(got), string(fileData))
		}
	})

	t.Run("DecodeMultiple Messages", func(t *testing.T) {
		msg1 := tlv2.NewNumberMessage(42)
		msg2 := tlv2.NewStringMessage("Test")
		msg3 := tlv2.NewFileMessage([]byte("Data"))

		mulBytes, err := tlv2.EncodeMultiple([]*tlv2.TLVMessage{msg1, msg2, msg3})
		if err != nil {
			t.Fatalf("Encode multiple error %v", err)
		}

		msgs, err := tlv2.DecodeMultiple(mulBytes)
		if err != nil {
			t.Fatalf("DecodeMultiple error: %v", err)
		}

		if len(msgs) != 3 {
			t.Fatalf("DecodeMultiple returned %d messages; want 3", len(msgs))
		}

		if n, _ := msgs[0].GetNumber(); n != 42 {
			t.Errorf("first message number = %d; want 42", n)
		}
		if s, _ := msgs[1].GetString(); s != "Test" {
			t.Errorf("second message string = %q; want 'Test'", s)
		}
		if f, _ := msgs[2].GetFileData(); string(f) != "Data" {
			t.Errorf("third message file = %q; want 'Data'", string(f))
		}
	})
}
