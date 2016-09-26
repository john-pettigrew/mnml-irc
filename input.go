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
		messagesList.ListMove(-1)
		return
	}
	//down
	if bytes.Equal([]byte(input)[:3], []byte{27, 79, 66}) {
		messagesList.ListMove(1)
		return
	}

	for _, r := range input {
		if r == rune(0) {
			continue
		}
		if r == '\r' {
			//handle buffer
			if buffer.Contents == "/quit" {
				ui.StopLoop()
				return
			} else if len(buffer.Contents) >= 8 && buffer.Contents[:8] == "/connect" {

				messagesList.AddMessage("IRC", "Connecting to host...")
				ircConn, err = irc.NewClient(buffer.Contents[9:])
				if err != nil {
					messagesList.AddMessage("Error", "Error connecting to host: "+err.Error())
					buffer.Clear()
				}

				messagesList.AddMessage("IRC", "Connected to server")
				go ircConn.SubscribeForMessages(&msgCh)
				buffer.Clear()
			} else if ircConn == (irc.Client{}) {
				messagesList.AddMessage("Error", "You must connect to a server first")
				buffer.Clear()
			} else {
				msg, err := message.ParseCommand(buffer.Contents)
				if err != nil {
					panic(err)
				}
				err = ircConn.SendMessage(msg)
				if err != nil {
					panic(err)
				}
				msgCh <- msg

				buffer.Clear()
			}
		} else {
			buffer.Type(r)
		}

		textInput.Text = inputPrefix + buffer.Contents
		ui.Render(ui.Body)
	}
}
