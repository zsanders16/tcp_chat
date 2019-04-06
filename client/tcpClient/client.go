package tcpClient

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"tcp_chat/types"
)

const Port = ":3333"

func NewEntryPoint() *EntryPoint {
	return &EntryPoint{}
}

type EntryPoint struct {
	Conn net.Conn
	Name string
}

func (e *EntryPoint) Open(address string) (*bufio.ReadWriter, error) {
	var err error
	e.Conn, err = net.Dial("tcp", address+Port)
	if err != nil {
		log.Println("Error creating connection on client: ", err.Error())
	}
	return bufio.NewReadWriter(bufio.NewReader(e.Conn), bufio.NewWriter(e.Conn)), nil
}

func (e *EntryPoint) SendMessage(message types.MessageCommand, rw *bufio.ReadWriter) {
	_, err := rw.WriteString("MESSAGE\n")
	if err != nil {
		log.Println("Error sending STRING command: ", err.Error())
	}

	enc := gob.NewEncoder(rw)
	if err := enc.Encode(message); err != nil {
		log.Println("Error encoding data: " + err.Error())
	}

	if err := rw.Flush(); err != nil {
		log.Println("Error flush string data: " + err.Error())
	}
}
