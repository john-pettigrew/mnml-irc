package main

import (
	ui "github.com/gizak/termui"
	"github.com/john-pettigrew/irc/message"
)

var err error
var msgCh = make(chan message.Message)
var inputPrefix = "> "
var buffer *InputBuffer

func main() {

	buffer = new(InputBuffer)

	err = initializeView()
	if err != nil {
		panic(err)
	}

	defer ui.Close()

	//Handle input
	go listenForInput()

	ui.Loop()
}
