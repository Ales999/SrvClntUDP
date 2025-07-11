package main

import (
	"fmt"
	"net"
	"os"
)

func PrintUsage() {
	fmt.Println("Пример использования:")
	fmt.Printf("%s serverip port\n", os.Args[0])
	fmt.Println("Пример: udpClient.exe 192.168.1.1 1845")
}

const message = "Работа с UDP сообщением"

func main() {
	if len(os.Args) < 3 {
		PrintUsage()
		return
	}

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%s", os.Args[1], os.Args[2]))
	if err != nil {
		fmt.Printf("Не удалось разрешить адрес: %v\n", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("Не удалось установить UDP соединение: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Подключено к: %s с портом %s\n", os.Args[1], os.Args[2])

	if _, err := conn.Write([]byte(message)); err != nil {
		fmt.Printf("Не удалось отправить сообщение: %v\n", err)
		os.Exit(1)
	}

	received := make([]byte, 1024)
	n, err := conn.Read(received)
	if err != nil {
		fmt.Printf("Не удалось прочитать ответ: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(received[:n]))
}
