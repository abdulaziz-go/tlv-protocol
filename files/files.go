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

func ReadFile(fd int, size int) ([]byte, error) {
	buffer := make([]byte, size)
	read, err := syscall.Read(fd, buffer)
	if err != nil {
		return nil, err
	}
	return buffer[:read], nil
}

func GetFdAlreadyExistFile(filename string) (int, error) {
	fd, err := syscall.Open(filename, syscall.O_RDONLY, 0)
	if err != nil {
		return -1, err
	}
	return fd, nil
}
