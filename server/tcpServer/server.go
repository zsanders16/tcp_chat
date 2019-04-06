package tcpServer

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"tcp_chat/types"
)

const Port = ":3333"

type HandlerFunc func(*bufio.ReadWriter, *Endpoint)

type User struct {
	Connection net.Conn
	Name       string
}

func NewEndpoint() *Endpoint {
	return &Endpoint{
		Handler: map[string]HandlerFunc{},
		Users:   []*User{},
	}
}

type Endpoint struct {
	Listener net.Listener
	Users    []*User
	Handler  map[string]HandlerFunc
	M        sync.RWMutex
}

func (e *Endpoint) Open(address string) error {
	var err error
	e.Listener, err = net.Listen("tcp", "localhost"+Port)
	return err
}

func (e *Endpoint) Listen() {
	fmt.Println("Listening on ", e.Listener.Addr().String())
	for {
		conn, err := e.Listener.Accept()
		if err != nil {
			log.Println("Unable to Accept connection request: ", err)
			continue
		}
		go e.HandleMessage(conn)
	}
}

func (e *Endpoint) HandleMessage(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()
	for {
		cmd, err := rw.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Println("Error reading command. Got: ", cmd+"\n"+err.Error())
			return
		}

		cmd = strings.Trim(cmd, "\n ")

		if cmd == "NAME" {
			e.RegisterUser(rw, conn)
			continue
		}

		e.M.RLock()
		handleCmd, ok := e.Handler[cmd]
		e.M.RUnlock()
		if !ok {
			log.Println("Command '" + cmd + "' not found. Please register the command")
			return
		}
		handleCmd(rw, e)
	}
}

func (e *Endpoint) RegisterHandler(name string, f HandlerFunc) {
	e.M.RLock()
	e.Handler[name] = f
	e.M.RUnlock()
}

func (e *Endpoint) RegisterUser(wr *bufio.ReadWriter, conn net.Conn) {

	var nameCommand types.NameCommand
	dec := gob.NewDecoder(wr)
	if err := dec.Decode(&nameCommand); err != nil {
		log.Println("Error decoding name", err.Error())
	}
	name := strings.Trim(nameCommand.Name, "\n")

	user := &User{
		Connection: conn,
		Name:       name,
	}
	e.M.RLock()
	e.Users = append(e.Users, user)
	e.M.RUnlock()
}

func (e *Endpoint) DeregisterUser() {

}

func (e *Endpoint) Broadcast(m types.MessageCommand) {
	for _, user := range e.Users {
		w := bufio.NewWriter(user.Connection)
		enc := gob.NewEncoder(w)
		if err := enc.Encode(m); err != nil {
			log.Println("Unable to encode broadcast: ", err.Error())
		}
		if err := w.Flush(); err != nil {
			log.Println("Unable to flush broadcast: ", err.Error())
		}
	}

}
