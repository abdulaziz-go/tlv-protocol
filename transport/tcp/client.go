package tcp

import (
	"net"
	"syscall"
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
