package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"syscall"
	"time"

	"redigo/constant"
	"redigo/internal/config"
	"redigo/internal/core"
	"redigo/internal/core/io_multiplexing"
)

func readCommand(fd int) (*core.RedigoCommand, error) {
	var buf = make([]byte, 512)
	n, err := syscall.Read(fd, buf)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, io.EOF
	}
	return core.ParseCmd(buf)
}

func RunIoMultiplexingServer() error {
	// start tcp server
	address := config.Host + ":" + strconv.Itoa(config.Port)
	listener, err := net.Listen(config.Protocol, address)
	if err != nil {
		log.Fatal()
	}
	defer listener.Close()

	// get file descriptor from listener
	tcpListener, ok := (listener).(*net.TCPListener)
	if !ok {
		log.Fatal("listener is not TCP server")
	}
	listenerFile, err := tcpListener.File()
	if err != nil {
		log.Fatal(err)
	}
	defer listenerFile.Close()
	serverFd := int(listenerFile.Fd())

	// create ioMultiplexer instance (epoll in Linux, kqueue in MacOS)
	ioMultipler, err := io_multiplexing.CreateIOMultiplexer()
	if err != nil {
		log.Fatal(err)
	}
	defer ioMultipler.Close()

	// Monitor <read> event on <serverFd>
	if err := ioMultipler.Monitor(io_multiplexing.Event{
		Fd: serverFd,
		Op: io_multiplexing.OpRead,
	}); err != nil {
		log.Fatal(err)
	}

	go func() {
		ticker := time.NewTicker(constant.ActiveExpireFrequency)
		defer ticker.Stop()
		for range ticker.C {
			core.ActiveDeleteExpiredKeys()
		}
	}()

	for {
		events, err := ioMultipler.Wait()
		if err != nil {
			continue
		}

		for i := 0; i < len(events); i++ {
			if events[i].Fd == serverFd {
				log.Println("connect to server")
				connFd, _, err := syscall.Accept(serverFd)
				if err != nil {
					log.Println("error", err)
					continue
				}

				log.Println("set up new connection")
				if err = ioMultipler.Monitor(io_multiplexing.Event{
					Fd: connFd,
					Op: io_multiplexing.OpRead,
				}); err != nil {
					log.Fatal(err)
				}
			} else {
				log.Println("read new command")
				cmd, err := readCommand(events[i].Fd)
				if err != nil {
					if err == io.EOF || err == syscall.ECONNRESET {
						log.Println("client disconnected")
						syscall.Close(events[i].Fd)
						continue
					}
					log.Println("read error", err)
					continue
				}
				log.Printf("parsed command: %+v\n", cmd)

				data := core.ExecuteCommand(cmd)
				if _, err := syscall.Write(events[i].Fd, data); err != nil {
					log.Println("error write: ", err)
				}
			}
		}
	}
}
