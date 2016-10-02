package main

import ui "github.com/gizak/termui"

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

	ui.Loop()
}
