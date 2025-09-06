package files

import (
	"fmt"
	"syscall"
)

func CreateFile(name string) (int, error) {
	fDescriptor, err := syscall.Open(name, syscall.O_CREAT|syscall.O_RDWR, 0644)
	if err != nil {
		return -1, err
	}

	return fDescriptor, nil
}

func WriteFile(fd int, data []byte) error {
	byteCount, err := syscall.Write(fd, data)
	if err != nil {
		return err
	}
	fmt.Println("written bytes ", byteCount)
	return err
}

func ReadFile(fd int, size int) (string, error) {
	buffer := make([]byte, 0)
	read, err := syscall.Read(fd, buffer)
	if err != nil {
		return "", err
	}

	return string(buffer[:read]), nil
}
