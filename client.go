package main

import (
	ui "github.com/gizak/termui"
	"github.com/john-pettigrew/irc/message"
)

var err error
var inputPrefix = "> "
var buffer *InputBuffer
var serverList *ServerList

func main() {

	serverList = NewServerList()
	buffer = new(InputBuffer)

	err = initializeView()
	if err != nil {
		panic(err)
	}

	defer ui.Close()

	//Handle input
	go listenForInput()

	//add welcome message
	serverList.AddMessage(message.Message{Command: "IRC", Options: []string{"Welcome to mnmlIRC!"}})
	renderScreen()

	ui.Loop()
}
