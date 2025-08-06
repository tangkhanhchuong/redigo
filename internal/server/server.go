package server

import (
	"log"
	"net"

	"redigo/internal/core"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	for {
		var buf []byte = make([]byte, 1000)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println("read error:", err)
			return
		}
		data := buf[:n]

		cmd, err := core.ReadRESPCommand(data)
		if err != nil {
			log.Println("client disconnected", err)
			return
		}

		log.Printf("parsed command: %+v\n", cmd)
		conn.Write([]byte("+PONG\r\n"))
	}
}

func Run(protocol string, addr string) error {
	listener, err := net.Listen(protocol, addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	log.Println("listening on", addr)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error", err)
			panic(err)
		}
		log.Println("a new client connected with address:", conn.RemoteAddr())

		go handleConnection(conn)
	}
}
