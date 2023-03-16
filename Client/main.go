package main

import (
	"fmt"
	"net"
	"os"
)

func PrintUsage() {
	fmt.Println("Usage example:")
	fmt.Printf("%s serverip port\n", os.Args[0])
	fmt.Println("Example: client.exe 192.168.1.1 1845")
}

func main() {

	var srvhost string
	// Порт который будем слушать
	var srvport string

	if len(os.Args) > 2 {
		srvhost = os.Args[1]
		srvport = os.Args[2]
	} else {
		PrintUsage()
		return
	}

	fmt.Print("Connected to: ", srvhost)
	fmt.Println(" with port", srvport)

	udpServer, err := net.ResolveUDPAddr("udp", srvhost+":"+srvport)

	if err != nil {
		println("ResolveUDPAddr failed:", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, udpServer)
	if err != nil {
		println("Listen failed:", err.Error())
		os.Exit(1)
	}

	//close the connection
	defer conn.Close()

	_, err = conn.Write([]byte("Work a UDP message"))
	if err != nil {
		println("Write data failed:", err.Error())
		os.Exit(1)
	}

	// buffer to get data
	received := make([]byte, 1024)
	// Number of received bytes
        var n int
	n, err = conn.Read(received)
	if err != nil {
		println("Read data failed:", err.Error())
		os.Exit(1)
	}

	println(string(received[:n]))
}
