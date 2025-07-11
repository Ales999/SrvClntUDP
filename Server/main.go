package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
	"unicode"
)

func PrintUsage() {
	fmt.Println("Пример использования:")
	fmt.Printf("%s <port>\n", os.Args[0])
}

var wg sync.WaitGroup

func main() {
	if len(os.Args) < 2 {
		PrintUsage()
		return
	}
	// Проверим корректность номера порта
	port, valid := validatePort(os.Args[1])
	if !valid {
		fmt.Printf("Ошибка: '%s' - не валидный порт (должно быть число от 1 до 65535)\n", os.Args[1])
		PrintUsage()
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	udpServer, err := net.ListenPacket("udp4", ":"+port)
	if err != nil {
		log.Fatal(err)
	}
	defer udpServer.Close()

	go func(conn net.PacketConn) {
		<-sigChan
		fmt.Println("\nОстанавливаем сервер...")
		cancel()
		udpServer.Close() // Закрываем соединение перед остановкой сервера
	}(udpServer)

	fmt.Println("Сервер запущен и ожидает UDP-пакеты на порту:", port)

	buf := make([]byte, 1024)

	for ctx.Err() == nil {
		n, addr, err := udpServer.ReadFrom(buf)
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			log.Println("Ошибка чтения из UDP:", err, "IP:", addr.String())
			continue
		}

		fmt.Printf("packet-received: bytes=%d from=%s\n", n, addr.String())
		wg.Add(1) // Инкрементируем WaitGroup перед запуском новой горутины
		go func(ctx context.Context, conn net.PacketConn, addr net.Addr, msg []byte) {
			defer wg.Done() // Декриментируем WaitGroup в конце горутины
			response(ctx, conn, addr, msg)
		}(ctx, udpServer, addr, buf[:n])
	}

	// Ждём завершения всех горутин перед закрытием сервера
	wg.Wait()
}

func response(ctx context.Context, conn net.PacketConn, addr net.Addr, msg []byte) {
	select {
	case <-ctx.Done():
		return
	default:
		timeStr := time.Now().Format(time.ANSIC)
		response := fmt.Sprintf("time received: %v. Your message: %v!", timeStr, string(msg))
		if _, err := conn.WriteTo([]byte(response), addr); err != nil {
			log.Println("Ошибка записи в UDP:", err)
		}
	}
}

func validatePort(portStr string) (string, bool) {
	portInt, err := strconv.Atoi(portStr)
	if err != nil || portInt < 1 || portInt > 65535 {
		return "0", false
	}

	// Check that all characters are digits
	for _, c := range portStr {
		if !unicode.IsDigit(c) {
			return "0", false
		}
	}
	return strconv.Itoa(portInt), true
}
