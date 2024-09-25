package libs

import (
	"log"
	"net"
	"reflect"
	"sync"
	"syscall"

	"golang.org/x/sys/unix"
)

type Epoll struct {
	epfd        int
	connections map[int]net.Conn
	lock        *sync.RWMutex
}

func NewEpoll() *Epoll {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		panic(err)
	}
	return &Epoll{
		epfd:        fd,
		lock:        &sync.RWMutex{},
		connections: make(map[int]net.Conn),
	}
}

func (e *Epoll) Add(connection net.Conn) error {
	fd := getWebsocketFd(connection)
	err := unix.EpollCtl(e.epfd, syscall.EPOLL_CTL_ADD, fd, &unix.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	e.connections[fd] = connection
	return nil
}

func (e *Epoll) Remove(connection net.Conn) error {
	fd := getWebsocketFd(connection)
	err := unix.EpollCtl(e.epfd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	delete(e.connections, fd)
	return nil
}

func (e *Epoll) Wait() ([]net.Conn, error) {
	events := make([]unix.EpollEvent, 100)
	n, err := unix.EpollWait(e.epfd, events, 100)
	if err != nil {
		return nil, err
	}
	e.lock.RLock()
	defer e.lock.RUnlock()
	connections := make([]net.Conn, 0, n)
	for i := 0; i < n; i += 1 {
		connection := e.connections[int(events[i].Fd)]
		connections = append(connections, connection)
	}
	return connections, nil
}

func (e *Epoll) Close() {
	if err := unix.Close(e.epfd); err != nil {
		log.Printf("Failed to epoll close: %v", err)
	}
}

func getWebsocketFd(connection net.Conn) int {
	// tcp, err := connection.(*net.TCPConn).File()
	// log.Printf("tcp.Fd(): %+v\n", tcp.Fd())

	tcp := reflect.Indirect(reflect.ValueOf(connection)).FieldByName("conn")
	fd := tcp.FieldByName("fd")
	pfd := reflect.Indirect(fd).FieldByName("pfd")
	return int(pfd.FieldByName("Sysfd").Int())
}
