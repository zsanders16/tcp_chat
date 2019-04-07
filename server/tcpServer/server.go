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

	"github.com/zsanders16/tcp_chat/types"
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
		user := &User{
			Connection: conn,
		}
		go e.HandleMessage(user)
	}
}

func (e *Endpoint) HandleMessage(user *User) {
	rw := bufio.NewReadWriter(bufio.NewReader(user.Connection), bufio.NewWriter(user.Connection))
	// defer conn.Close()
	defer e.DeregisterUser(user)
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
			e.RegisterUser(rw, user)
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

func (e *Endpoint) RegisterUser(wr *bufio.ReadWriter, user *User) {

	var nameCommand types.NameCommand
	dec := gob.NewDecoder(wr)
	if err := dec.Decode(&nameCommand); err != nil {
		log.Println("Error decoding name", err.Error())
	}
	name := strings.Trim(nameCommand.Name, "\n")
	user.Name = name
	e.M.RLock()
	e.Users = append(e.Users, user)
	e.M.RUnlock()
}

func (e *Endpoint) DeregisterUser(u *User) {
	e.M.Lock()
	defer e.M.Unlock()

	for i, user := range e.Users {
		if u == user {
			e.Users = append(e.Users[:i], e.Users[i+1:]...)
		}
	}
	u.Connection.Close()
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
