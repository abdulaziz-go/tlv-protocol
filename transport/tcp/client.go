package tcp

import (
	"fmt"
	"net"
	"syscall"
	"tlv-protocol/tlv"
)

type Client struct {
	fd     int // client file descriptor
	addr   net.Addr
	server *TLVServer
	buffer []byte
}

func (c *Client) Close() {
	if c.fd != 0 {
		syscall.Close(c.fd)
	}
}

func (c *Client) SendTLV(msg *tlv.TLVMessage) error {
	encoded, err := msg.Encode()
	if err != nil {
		return err
	}

	totalWritten := 0

	for totalWritten < len(encoded) {
		n, err := syscall.Write(c.fd, encoded[totalWritten:])
		if err != nil {
			return fmt.Errorf("write error: %v", err)
		}
		totalWritten += n
	}

	fmt.Printf("Client %s ga message yuborildi: %s (%d bytes)",
		c.addr, msg, len(encoded))

	return nil
}
