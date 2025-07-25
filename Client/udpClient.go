package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
	"unicode"
)

func PrintUsage() {
	fmt.Println("Условие использования:")
	fmt.Printf("%s <server_ip> <port>\n", os.Args[0])
	fmt.Println("Пример: client.exe 192.168.1.1 1845")
}

func validateIP(ipStr string) bool {
	if ip := net.ParseIP(ipStr); ip != nil {
		return true
	}
	return false
}

func validatePort(portStr string) (int, bool) {
	portInt, err := strconv.Atoi(portStr)
	if err != nil || portInt < 0 || portInt > 65535 {
		return 0, false
	}

	// Проверяем, что все символы - цифры
	for _, c := range portStr {
		if !unicode.IsDigit(c) {
			return 0, false
		}
	}

	return portInt, true
}

const message = "Работа с UDP сообщением"

func main() {
	if len(os.Args) != 3 {
		PrintUsage()
		return
	}

	host := os.Args[1]
	portStr := os.Args[2]

	// Проверяем вводы
	if !validateIP(host) {
		fmt.Printf("Ошибка: '%s' - не валидный IP адрес\n", host)
		PrintUsage()
		return
	}

	port, portValid := validatePort(portStr)
	if !portValid {
		fmt.Printf("Ошибка: '%s' - не валидный порт (должно быть число от 0 до 65535)\n", portStr)
		PrintUsage()
		return
	}

	fmt.Printf("Отправляю UDP сообщение на адрес: %s:%d\n", host, port)

	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Printf("Ошибка при разрешении адреса: %v\n", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		fmt.Printf("Ошибка при установлении соединения: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(message)); err != nil {
		fmt.Printf("Ошибка при отправке сообщения: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Сообщение успешно отправлено, ожидаем ответ в течение трех секунд...")

	// Создаем контекст с таймаутом в 3 секунды
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dataChan := make(chan []byte)

	// Горутина для чтения данных
	go func() {
		received := make([]byte, 2048)
		n, err := conn.Read(received)
		if err != nil {
			fmt.Printf("Ошибка при чтении ответа: %v\n", err)
			dataChan <- nil // или можно отправить ошибку
			return
		}
		dataChan <- received[:n]
	}()

	// Ожидаем данные или completion контекста
	select {
	case data := <-dataChan:
		if data == nil {
			fmt.Println("Ошибка при чтении ответа")
			os.Exit(1)
		}
		fmt.Println(string(data))
	case <-ctx.Done():
		fmt.Println("!!! Внимание !!! - ответ так не получен в течение 3 секунд :-(")
		os.Exit(1)
	}
}
