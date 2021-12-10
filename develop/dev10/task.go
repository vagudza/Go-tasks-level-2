package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.


Telnet - это сетевая утилита, которая позволяет соединиться с удаленным портом любого компьютера и установить интерактивный канал связи,
 например, для передачи команд или получения информации. Можно сказать, что это универсальный браузер в терминале, который умеет работать
 со множеством сетевых протоколов.

Пример использования:
 go run task.go --timeout=10s opennet.ru 80
 GET /		- после подключения выполняем запрос

 go run task.go --timeout=10s 127.0.0.1 8080    - подключение к локальному TCP серверу. Можно подключиться из консоли дополнительно: telnet 127.0.0.1 8080
*/

type Config struct {
	Timeout          time.Duration
	Host             string
	Port             string
	LocalServerStart bool
	LocalServerHost  string
	LocalServerPort  string
}

// NewConfig - конструктор, парсящий флаги и аргументы
func NewConfig() *Config {
	conf := Config{}

	flag.Usage = func() {
		fmt.Println("Usage flags: [--timeout] host port")
		flag.PrintDefaults()
	}

	timeout := flag.Duration("timeout", time.Second*10, "timeout")

	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	conf.Host = args[0]
	conf.Port = args[1]
	conf.Timeout = *timeout
	conf.LocalServerHost = "127.0.0.1"
	conf.LocalServerPort = "8080"
	conf.LocalServerStart = true

	return &conf
}

// Start - Точка входа в программу
func Start(conf *Config) {
	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, syscall.SIGINT)
	go func() {
		<-signalChan
		cancel()
	}()

	if conf.LocalServerStart {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go startLocalTCPServer(conf, wg)
		wg.Wait()
	}

	addr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(conf.Host, conf.Port))
	if err != nil {
		log.Fatal("can not to resolve tcp addres:", err)
	}

	conn, err := net.DialTimeout(addr.Network(), addr.String(), conf.Timeout)
	if err != nil {
		log.Fatal("timeout to dial connection:", err)
	}

	defer conn.Close()

	go read(conn, cancel)
	go write(conn, cancel)

	<-ctx.Done()
	log.Println("finish telnet client")
}

func read(conn net.Conn, cancelFunc context.CancelFunc) {
	// читаем из conn, записываем в stdout
	scanner := bufio.NewScanner(conn)
	for {
		if !scanner.Scan() {
			log.Printf("read: connection was closed")
			cancelFunc()
			return
		}
		text := scanner.Text()
		fmt.Printf("%s\n", text)
	}
}

func write(conn net.Conn, cancelFunc context.CancelFunc) {
	// читаем из stdin, записываем в conn
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			log.Printf("write: can not scan")
			cancelFunc()
			return
		}
		str := scanner.Text()

		_, err := conn.Write([]byte(fmt.Sprintln(str)))
		if err != nil {
			log.Println("write: can not send to server", err)
			cancelFunc()
			return
		}
	}
}

func startLocalTCPServer(conf *Config, wg *sync.WaitGroup) {
	listner, err := net.Listen("tcp", fmt.Sprintf("%s:%s", conf.LocalServerHost, conf.LocalServerPort))
	if err != nil {
		log.Fatal("local server: can not listen", err)
	}

	wg.Done()

	for {
		conn, err := listner.Accept()
		if err != nil {
			log.Fatal("can not Accept connection", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	conn.Write([]byte("Hello, " + remoteAddr + "\n\r"))
	log.Printf("%+v connected\n", remoteAddr)
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()

		if text == "exit" {
			conn.Write([]byte(fmt.Sprintf("Bye, %+v\n\r", remoteAddr)))
			log.Printf("%+v disconnected\n", remoteAddr)
			break
		} else if text != "" {
			conn.Write([]byte(fmt.Sprintf("%+v message is '%s'\n\r", remoteAddr, text)))
		}
	}
}

func main() {
	conf := NewConfig()
	Start(conf)
}
