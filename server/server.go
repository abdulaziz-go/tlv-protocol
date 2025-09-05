package main

import (
	"fmt"
	"syscall"
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		panic(err)
	}

	addr := syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{0, 0, 0, 0},
	}

	if err = syscall.Bind(fd, &addr); err != nil {
		panic(err)
	}

	fmt.Println("Server running on port 8080...")

	for {
		//clientFd, _, err := syscall.Accept(fd)
		//if err != nil {
		//	fmt.Println("Accept error ", err)
		//	continue
		//}

		//go handleClient(clientFd)
	}
}
