package main

import (
	"strings"

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

		msgStr := message.Marshal(newMsg)

		//remove extra characters
		msgStr = strings.Replace(msgStr, "\n", "", -1)
		msgStr = strings.Replace(msgStr, "\r", "", -1)

		m.Messages = append(m.Messages, msgStr)
		m.setMessageItems()
	}
}

func (m *MessageList) setMessageItems() {

	//calculate how many messages can fit
	var currentMessages []string
	availableHeight := m.InnerHeight()
	availableWidth := m.Width

	for i := len(m.Messages) - 1; i >= 0; i-- {
		availableHeight -= len(m.Messages[i]) / availableWidth
		if len(m.Messages[i])%availableWidth > 0 {
			availableHeight -= 1
		}

		if availableHeight < 1 {
			break
		}

		currentMessages = append(currentMessages, m.Messages[i])
	}

	//reverse message order
	var reversedMessages []string
	for i := len(currentMessages) - 1; i >= 0; i-- {
		reversedMessages = append(reversedMessages, currentMessages[i])
	}

	m.Items = reversedMessages
	renderScreen()
}

func (m *MessageList) ListMove(y int) {
	m.SetY(m.Y + y)
}
