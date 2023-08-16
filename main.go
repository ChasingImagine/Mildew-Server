package main

import (
	"aftermildewserver/transforms"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"time"
)

var clientCount = 0

func main() {
	listener, err := net.Listen("tcp", "localhost:12345")
	if err != nil {
		fmt.Println("Hata:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Sunucu 12345 portunda dinliyor")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Hata bağlantı kabul edilirken:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	clientCount++
	defer func() {
		fmt.Printf("Bağlantı kapatıldı. Toplam istemci sayısı: %d\n", clientCount)
		conn.Close()
		clientCount--
	}()

	rand.Seed(time.Now().UnixNano())

	go sendResponse(conn)

	receivedData := make([]byte, 4096)
	for {
		n, err := conn.Read(receivedData)
		if err != nil {
			fmt.Println("Hata veri alırken:", err)
			return
		}

		receivedMessage := transforms.Transforms{}
		err = json.Unmarshal(receivedData[:n], &receivedMessage)
		if err != nil {
			fmt.Println("Hata JSON çözme sırasında:", err)
			return
		}

		fmt.Printf("Gelen veri: %+v\n", receivedMessage)
	}
}

func sendResponse(conn net.Conn) {
	for {
		message := transforms.Transforms{
			Position: transforms.Positions{X: float64(rand.Intn(10)), Y: float64(rand.Intn(10)), Z: float64(rand.Intn(10))},
			Rotation: transforms.Rotations{X: float64(rand.Intn(10)), Y: float64(rand.Intn(10)), Z: float64(rand.Intn(10))},
		}

		data, err := json.Marshal(message)
		if err != nil {
			fmt.Println("Hata JSON kodlamada:", err)
			return
		}

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Hata veri gönderirken:", err)
			return
		}

		//fmt.Printf("Sunucudan gönderilen mesaj: %s \n", data)

		time.Sleep(time.Second) // 1 saniye bekle
	}
}
