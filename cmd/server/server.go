package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"tlv-protocol/tlv"
	"tlv-protocol/transport/tcp"
)

func main() {
	log.Println("=== TLV Protocol Server ===")

	server := tcp.NewTlVServer(8080)
	setupHandlers(server)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("Server start error:", err)
		}
	}()

	<-c
	log.Println("Shutdown signal olindi...")

	server.Stop()
	log.Println("Server to'xtatildi")
}

func setupHandlers(server *tcp.TLVServer) {
	server.RegisterHandler(tlv.TypeNumber, handleNumber)
	server.RegisterHandler(tlv.TypeString, handleString)
	server.RegisterHandler(tlv.TypeFile, handleFile)
}

func handleNumber(client *tcp.Client, msg *tlv.TLVMessage) error {
	number, err := msg.GetNumber()
	if err != nil {
		return err
	}
	log.Printf("Number qabul qilindi: %d", number)

	result := number * 2
	response := tlv.NewNumberMessage(result)

	if err := client.SendTLV(response); err != nil {
		return err
	}

	log.Printf("Number response yuborildi: %d -> %d", number, result)
	return nil
}

func handleString(client *tcp.Client, msg *tlv.TLVMessage) error {
	text, err := msg.GetString()
	if err != nil {
		return err
	}

	log.Printf("String qabul qilindi: %q", text)

	response := tlv.NewStringMessage("Echo: " + text)

	if err := client.SendTLV(response); err != nil {
		return err
	}

	log.Printf("String echo yuborildi: %q", text)
	return nil
}

func handleFile(client *tcp.Client, msg *tlv.TLVMessage) error {
	fileData, err := msg.GetFileData()
	if err != nil {
		return err
	}

	log.Printf("File qabul qilindi: %d bytes", len(fileData))

	info := fmt.Sprintf("File received: %d bytes",
		len(fileData))

	response := tlv.NewStringMessage(info)

	if err := client.SendTLV(response); err != nil {
		return err
	}

	log.Printf("File info response yuborildi")
	return nil
}
