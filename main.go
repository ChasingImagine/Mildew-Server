package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	Text string `json:"text"`
	Sayi int    `json:"sayi"`
}

func main() {
	listener, err := net.Listen("tcp", "localhost:12345")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on port 12345")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	message := Message{
		Text: "Merhaba  unity",
		Sayi: 80,
	}

	data, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	// Unity'den gelen cevabı okuma
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	responseData := buffer[:n]
	var responseMessage Message
	err = json.Unmarshal(responseData, &responseMessage)
	if err != nil {
		fmt.Println("Error unmarshaling response JSON:", err)
		return
	}

	fmt.Println("Unity'den gelen cevap:", responseMessage.Text, responseMessage.Sayi)
}
