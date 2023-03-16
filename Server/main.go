package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func PrintUsage() {
	fmt.Println("Usage example:")
	fmt.Printf("%s 1457\n", os.Args[0])
}

func main() {

	// Порт который будем слушать
	var port string

	if len(os.Args) > 1 {
		port = os.Args[1]
	} else {
		PrintUsage()
		return
	}

	// listen to incoming udp packets
	udpServer, err := net.ListenPacket("udp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()

	fmt.Println("Sever started and wait UDP packet to this port:", port)

	for {
		buf := make([]byte, 1024)
		_, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			continue
		}
		go response(udpServer, addr, buf)
	}

}

func response(udpServer net.PacketConn, addr net.Addr, buf []byte) {
	time := time.Now().Format(time.ANSIC)
	responseStr := fmt.Sprintf("time received: %v. Your message: %v!", time, string(buf))

	udpServer.WriteTo([]byte(responseStr), addr)
}
