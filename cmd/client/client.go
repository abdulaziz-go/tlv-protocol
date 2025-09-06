package main

import (
	"fmt"
	"net"
	"syscall"
	"time"
	"tlv-protocol/files"
	"tlv-protocol/tlv"
)

func main() {
	fmt.Println("=== TLV Protocol Client ===")

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if err := testNumber(conn, 42); err != nil {
		fmt.Printf("Number test error: %v", err)
	}

	if err := testString(conn, "Salom Server!"); err != nil {
		fmt.Printf("String test error: %v", err)
	}

	if err := testFile(conn, "test.txt"); err != nil {
		fmt.Printf("File test error: %v", err)
	}

	time.Sleep(5 * time.Second)
}

func testNumber(conn net.Conn, number int64) error {
	fmt.Printf("Yuborish: %d\n", number)

	msg := tlv.NewNumberMessage(number)
	if err := sendTLV(conn, msg); err != nil {
		return err
	}

	response, err := receiveTLV(conn)
	if err != nil {
		return err
	}

	if response.Type != tlv.TypeNumber {
		return fmt.Errorf("expected number response, got type %d", response.Type)
	}

	result, err := response.GetNumber()
	if err != nil {
		return err
	}

	fmt.Printf("Qabul qilindi: %d\n", result)
	fmt.Printf("Server %d ni %d ga aylantirdi ✓\n", number, result)

	return nil
}

func testString(conn net.Conn, text string) error {
	fmt.Printf("Yuborish: %q\n", text)

	msg := tlv.NewStringMessage(text)
	if err := sendTLV(conn, msg); err != nil {
		return err
	}

	response, err := receiveTLV(conn)
	if err != nil {
		return err
	}

	if response.Type != tlv.TypeString {
		return fmt.Errorf("expected string response, got type %d", response.Type)
	}

	result, err := response.GetString()
	if err != nil {
		return err
	}

	fmt.Printf("Qabul qilindi: %q\n", result)

	return nil
}

func testFile(conn net.Conn, filename string) error {
	fd, err := files.GetFdAlreadyExistFile(filename)
	if err != nil {
		return err
	}

	fileData, err := files.ReadFile(fd, 4*1024*1024)
	if err != nil {
		return err
	}

	defer syscall.Close(fd)

	msg := tlv.NewFileMessage(fileData)
	if err := sendTLV(conn, msg); err != nil {
		return err
	}

	response, err := receiveTLV(conn)
	if err != nil {
		return err
	}

	if response.Type != tlv.TypeString {
		return fmt.Errorf("expected string response, got type %d", response.Type)
	}

	result, err := response.GetString()
	if err != nil {
		return err
	}

	fmt.Printf("Server response: %s\n", result)
	fmt.Println("File upload muvaffaqiyatli ✓")

	return nil
}

func sendTLV(conn net.Conn, msg *tlv.TLVMessage) error {
	encoded, err := msg.Encode()
	if err != nil {
		return fmt.Errorf("encode error: %v", err)
	}

	_, err = conn.Write(encoded)
	return err
}

func receiveTLV(conn net.Conn) (*tlv.TLVMessage, error) {
	header := make([]byte, 5)
	if _, err := conn.Read(header); err != nil {
		return nil, fmt.Errorf("header read error: %v", err)
	}

	length := uint32(header[1]) | uint32(header[2])<<8 | uint32(header[3])<<16 | uint32(header[4])<<24

	value := make([]byte, length)
	if length > 0 {
		if _, err := conn.Read(value); err != nil {
			return nil, fmt.Errorf("value read error: %v", err)
		}
	}

	fullMessage := append(header, value...)

	return tlv.Decode(fullMessage)
}
