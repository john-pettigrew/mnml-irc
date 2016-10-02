package main

import (
	"bytes"
	"os"
	"time"

	ui "github.com/gizak/termui"
	"github.com/john-pettigrew/irc"
	"github.com/john-pettigrew/irc/message"
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

			if buffer.Contents == "" {
				return
			}

			if buffer.Contents == "/quit" {
				ui.StopLoop()
				return
			} else if len(buffer.Contents) >= 8 && buffer.Contents[:8] == "/connect" {

				// Connect to server
				serverList.Connect(buffer.Contents[9:])

				buffer.Clear()
			} else if ircConn == (irc.Client{}) {
				serverList.AddMessage(message.Message{Command: "Error", Options: []string{"You must connect to a server first"}})
				buffer.Clear()
			} else {
				msg, err := message.ParseCommand(buffer.Contents)
				if err != nil {
					serverList.AddMessage(message.Message{Command: "Error", Options: []string{"Error parsing message: " + err.Error()}})
					break
				}
				err = ircConn.SendMessage(msg)
				if err != nil {
					serverList.AddMessage(message.Message{Command: "Error", Options: []string{"Error sending message: " + err.Error()}})
					break

				}

				serverList.AddMessage(msg)
				buffer.Clear()
			}
		} else {
			buffer.Type(r)
		}

		renderScreen()
	}

	renderScreen()
}
