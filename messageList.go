package main

import (
	ui "github.com/gizak/termui"
	"github.com/john-pettigrew/irc/message"
)

type MessageList struct {
	ui.List
	Messages []string
}

func (m *MessageList) AddMessage(speaker, newMessage string) {
	m.Messages = append(m.Messages, speaker+"\t\t\t"+newMessage)
	m.setMessageItems()
}

func (m *MessageList) ListenForMessages(msgCh chan message.Message) {
	// m.Items = append(m.Items, speaker+"\t\t\t"+newMessage)
	var newMsg message.Message
	var open bool
	for {
		newMsg, open = <-msgCh
		if !open {
			break
		}
		m.Messages = append(m.Messages, message.Marshal(newMsg))
		m.setMessageItems()
	}
}

func (m *MessageList) setMessageItems() {
	start := len(m.Messages) - 7
	if start < 0 {
		start = 0
	}

	m.Items = m.Messages[start:]
	renderScreen()
}

func (m *MessageList) ListMove(y int) {
	m.SetY(m.Y + y)
}
