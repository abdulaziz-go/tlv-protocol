package tcp

import (
	"fmt"
	"log"
	"net"
	"syscall"
	"tlv-protocol/tlv"
)

type TLVServer struct {
	port     int
	listener int // file descriptor
	running  bool
	handlers map[uint8]MessageHandler
	clients  map[int]*Client
}

type MessageHandler func(client *Client, msg *tlv.TLVMessage) error

func NewTlVServer(port int) *TLVServer {
	return &TLVServer{
		port:     port,
		running:  false,
		handlers: make(map[uint8]MessageHandler),
		clients:  make(map[int]*Client),
	}
}

func (s *TLVServer) RegisterHandler(msgType uint8, handler MessageHandler) {
	s.handlers[msgType] = handler
}

func (s *TLVServer) Start() error {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		return fmt.Errorf("failed to create socket %v", err)
	}

	s.listener = fd

	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		syscall.Close(fd)
		return fmt.Errorf("failed to set SO_REUSEADDR: %v", err)
	}

	addr := syscall.SockaddrInet4{
		Port: s.port,
		Addr: [4]byte{0, 0, 0, 0},
	}

	if err := syscall.Bind(fd, &addr); err != nil {
		syscall.Close(fd)
		return fmt.Errorf("failed to bind to port %d: %v", s.port, err)
	}

	if err := syscall.Listen(fd, 10); err != nil {
		syscall.Close(fd)
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.running = true

	fmt.Println(fmt.Sprintf("TLV Server ishga tushdi port %d", s.port))

	for s.running {
		clientFd, clientAddr, err := syscall.Accept(fd)
		if err != nil {
			if !s.running {
				break
			}
			log.Printf("Accept error: %v", err)
			continue
		}
		client := &Client{
			fd:     clientFd,
			addr:   parseAddr(clientAddr),
			server: s,
			buffer: make([]byte, 0, 4096),
		}

		s.clients[clientFd] = client

		fmt.Println(fmt.Sprintf("Yangi client ulandi: %s", client.addr))

		go s.handleClient(client)
	}
	return nil
}

func (s *TLVServer) handleClient(client *Client) {
	defer func() {
		client.Close()
		delete(s.clients, client.fd)
		log.Printf("Client ulanishi yopildi: %s", client.addr)
	}()

	readBuffer := make([]byte, 4096)

	for s.running {
		n, err := syscall.Read(client.fd, readBuffer)
		if err != nil {
			if err == syscall.EAGAIN || err == syscall.EWOULDBLOCK {
				fmt.Println("syscall eagain e wouldblock detected")
				continue
			}
			log.Printf("Client %s dan o'qishda xatolik: %v", client.addr, err)
			break
		}

		if n == 0 {
			log.Printf("Client %s ulanishni yopdi", client.addr)
			break
		}

		client.buffer = append(client.buffer, readBuffer[:n]...)

		if err := s.processClientBuffer(client); err != nil {
			log.Printf("Buffer processing error for client %s: %v", client.addr, err)
			break
		}
	}
}

func (s *TLVServer) processClientBuffer(client *Client) error {
	for len(client.buffer) >= 3 {
		length := uint32(client.buffer[1]) | uint32(client.buffer[2])<<8
		totalSize := 5 + int(length)

		if len(client.buffer) < totalSize {
			break
		}

		msg, err := tlv.Decode(client.buffer[:totalSize])
		if err != nil {
			return fmt.Errorf("decode error: %v", err)
		}

		log.Printf("Client %s dan message qabul qilindi: %s", client.addr, msg)

		if handler, exists := s.handlers[msg.Type]; exists {
			if err := handler(client, msg); err != nil {
				log.Printf("Handler error for type %d: %v", msg.Type, err)
			}
		} else {
			log.Printf("Handler topilmadi message type: %d", msg.Type)
			response := tlv.NewStringMessage("Unknown message type")
			client.SendTLV(response)
		}

		client.buffer = client.buffer[totalSize:]
	}

	return nil
}

func parseAddr(addr syscall.Sockaddr) net.Addr {
	switch a := addr.(type) {
	case *syscall.SockaddrInet4:
		return &net.TCPAddr{
			IP:   net.IPv4(a.Addr[0], a.Addr[1], a.Addr[2], a.Addr[3]),
			Port: a.Port,
		}
	case *syscall.SockaddrInet6:
		return &net.TCPAddr{
			IP:   net.IP(a.Addr[:]),
			Port: a.Port,
		}
	default:
		return &net.TCPAddr{IP: net.IPv4(0, 0, 0, 0), Port: 0}
	}
}

func (s *TLVServer) Stop() {
	if !s.running {
		return
	}

	s.running = false

	for fd, client := range s.clients {
		client.Close()
		delete(s.clients, fd)
	}

	if s.listener != 0 {
		syscall.Close(s.listener)
	}

	log.Println("Server to'xtatildi")
}
