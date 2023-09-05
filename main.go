package main

import (
	"aftermildewserver/players"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

var clientCount = 0

var corectionsMap = make(map[net.Addr]int)
var idMap = make(map[string]players.Player)

var mutexCorectionsMap sync.Mutex
var mutexIdMap sync.Mutex

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

		clientAddr := conn.RemoteAddr()
		clientIP, clientPort, err := net.SplitHostPort(clientAddr.String())
		log.Println(clientIP, "(:?/)", clientPort)
		for true {
			min := 1
			max := 999
			randomNumber := rand.Intn(max-min+1) + min
			mutexIdMap.Lock()
			if _, ok := idMap[strconv.Itoa(randomNumber)]; !ok {
				idMap[strconv.Itoa(randomNumber)] = players.Player{}
				mutexIdMap.Unlock()
				corectionsMap[clientAddr] = randomNumber
				break
			}
			mutexIdMap.Unlock()

		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	clientCount++
	defer func() {

		mutexIdMap.Lock()
		mutexCorectionsMap.Lock()

		fmt.Printf("Bağlantı kapatıldı. Toplam istemci sayısı: %d\n", clientCount)
		delete(idMap, strconv.Itoa(corectionsMap[conn.RemoteAddr()]))
		delete(corectionsMap, conn.RemoteAddr())

		log.Println(idMap)

		mutexIdMap.Unlock()
		mutexCorectionsMap.Unlock()

		conn.Close()
		clientCount--
	}()

	go sendResponse(conn)

	receivedData := make([]byte, 4096)
	for {
		n, err := conn.Read(receivedData)
		if err != nil {
			fmt.Println("Hata veri alırken:", err)
			return
		}

		var receivedMessage = make(map[string]players.Player)

		//receivedMessage :=  players.Player{}
		err = json.Unmarshal(receivedData[:n], &receivedMessage)
		if err != nil {
			fmt.Println("Hata JSON çözme sırasında:", err)
			return
		}

		mutexIdMap.Lock()
		mutexCorectionsMap.Lock()

		var id string = strconv.Itoa(corectionsMap[conn.RemoteAddr()])
		idMap[id] = receivedMessage["Palyer"]

		mutexCorectionsMap.Unlock()
		mutexIdMap.Unlock()

		//fmt.Printf("Gelen veri: %+v\n", receivedMessage["0"])
	}
}

func sendResponse(conn net.Conn) {
	time.Sleep(time.Second)
	for {

		mutexIdMap.Lock()
		mutexCorectionsMap.Lock()
		message1 := players.Player{
			Id:         strconv.Itoa(corectionsMap[conn.RemoteAddr()]),
			Transforms: idMap[strconv.Itoa(corectionsMap[conn.RemoteAddr()])].Transforms,
		}

		idMap[strconv.Itoa(corectionsMap[conn.RemoteAddr()])] = message1

		var message = make(map[string]string)

		jsondatas, _ := json.Marshal(message1)
		message["Player"] = string(jsondatas)

		delete(idMap, "0")

		jsondatasId, _ := json.Marshal(idMap)

		mutexCorectionsMap.Unlock()
		mutexIdMap.Unlock()

		message["Players"] = string(jsondatasId)

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

		time.Sleep(time.Second / 60) // 1/60 saniye bekle

		mutexIdMap.Lock()
		log.Println(idMap)
		mutexIdMap.Unlock()

	}
}
