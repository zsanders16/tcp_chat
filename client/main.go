package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/zsanders16/tcp_chat/client/tcpClient"
	"github.com/zsanders16/tcp_chat/types"
)

func main() {
	entryPoint := tcpClient.NewEntryPoint()
	rw, err := entryPoint.Open("localhost")
	if err != nil {
		panic("Unable to open connection to server")
	}
	defer entryPoint.Conn.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("What is your name? ")
	name, _ := reader.ReadString('\n')
	name = strings.Trim(name, "\n")
	entryPoint.Name = name

	// Send NameCommand
	nameCommand := types.NameCommand{
		Name: name,
	}

	_, err = rw.WriteString("NAME\n")
	if err != nil {
		log.Println("Error sending STRING command: ", err.Error())
	}

	enc := gob.NewEncoder(rw)
	if err := enc.Encode(nameCommand); err != nil {
		log.Println("Error encoding data: " + err.Error())
	}

	if err := rw.Flush(); err != nil {
		log.Println("Error flush string data: " + err.Error())
	}

	// Start UI
	ui := tcpClient.NewUI(name, rw, *entryPoint)
	ui.StartUI()

}
