package io_multiplexing

import (
	"log"
	"redigo/internal/config"
	"syscall"
)

type Epoll struct {
	fd            int
	epollEvents   []syscall.EpollEvent
	genericEvents []Event
}

func CreateIOMultiplexer() (*Epoll, error) {
	epollFD, err := syscall.EpollCreate1(0)
	if err != nil {
		log.Fatal()
		return nil, err
	}

	return &Epoll{
		fd:            epollFD,
		epollEvents:   make([]syscall.EpollEvent, config.MaxConnection),
		genericEvents: make([]Event, config.MaxConnection),
	}, nil
}

func (epoll *Epoll) Monitor(event Event) error {
	epollEvent := event.toNativeEvent()
	return syscall.EpollCtl(epoll.fd, syscall.EPOLL_CTL_ADD, event.Fd, &epollEvent)
}

func (epoll *Epoll) Wait() ([]Event, error) {
	n, err := syscall.EpollWait(epoll.fd, epoll.epollEvents, -1)
	if err != nil {
		return nil, err
	}
	for i := 0; i < n; i++ {
		epoll.genericEvents[i] = createGenericEvent(epoll.epollEvents[i])
	}
	return epoll.genericEvents[:n], nil
}

func (epoll *Epoll) Close() error {
	return syscall.Close(epoll.fd)
}
