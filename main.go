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
	/*
		interfaces, err := net.Interfaces()
		if err != nil {
			fmt.Println("Ağ arabirimleri alınamadı:", err)
			return
		}

		fmt.Println("Aşağıdaki ağ arabirimlerinden birini seçin:")
		for i, iface := range interfaces {
			fmt.Printf("%d: %s\n", i+1, iface.Name)
		}

		var choice int
		fmt.Print("Seçiminizi yapın (1, 2, 3, ...): ")
		_, err = fmt.Scan(&choice)
		if err != nil || choice < 1 || choice > len(interfaces) {
			fmt.Println("Geçersiz seçim.")
			return
		}

		selectedInterface := interfaces[choice-1]

		var selectedIP string

		addrs, err := selectedInterface.Addrs()
		if err != nil {
			fmt.Printf("Ağ arabirimi %s için IP adresleri alınamadı: %v\n", selectedInterface.Name, err)
			return
		}

		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				selectedIP = ipnet.IP.String()
				fmt.Printf("Seçilen ağ arabirimi %s için yerel IPv4 adresi: %s\n", selectedInterface.Name, selectedIP)
			}
		}
	*/
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
