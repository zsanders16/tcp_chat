package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"tcp_chat/server/tcpServer"
	"tcp_chat/types"
)

func main() {

	endPoint := tcpServer.NewEndpoint()
	endPoint.RegisterHandler("MESSAGE", HandleMessageCommand)

	endPoint.Open("localhost")
	endPoint.Listen()

}

func HandleMessageCommand(rw *bufio.ReadWriter, e *tcpServer.Endpoint) {
	var messageCommand types.MessageCommand
	dec := gob.NewDecoder(rw)
	if err := dec.Decode(&messageCommand); err != nil {
		log.Println("Error decoding message", err.Error())
	}

	fmt.Printf("%s: %s\n", messageCommand.Name, messageCommand.Message)
	e.Broadcast(messageCommand)
}
