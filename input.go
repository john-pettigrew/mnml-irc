package main

import (
	"bytes"
	"os"
	"time"

	"github.com/john-pettigrew/irc"
)

var ircConn irc.Client

func listenForInput() {
	input := make([]byte, 100)
	for {
		input = make([]byte, 100)
		_, err := os.Stdin.Read(input)
		if err != nil {
			break
		}
		handleInput(string(input))
		time.Sleep(time.Second / 60)
	}
}

func handleInput(input string) {

	//Check for direction input

	//left
	if bytes.Equal([]byte(input)[:3], []byte{27, 79, 68}) {
		buffer.CursorMove(-1)
		return
	}
	//right
	if bytes.Equal([]byte(input)[:3], []byte{27, 79, 67}) {
		buffer.CursorMove(1)
		return
	}

	//up
	if bytes.Equal([]byte(input)[:3], []byte{27, 79, 65}) {

		serverList.CurrentChannel().ListMove(1)
		return
	}
	//down
	if bytes.Equal([]byte(input)[:3], []byte{27, 79, 66}) {
		serverList.CurrentChannel().ListMove(-1)
		return
	}

	for _, r := range input {
		if r == rune(0) {
			continue
		}
		if r == '\r' {
			//handle buffer

			serverList.HandleInput(buffer.Contents)
			buffer.Clear()

		} else {
			buffer.Type(r)
		}

		renderScreen()
	}

	renderScreen()
}
