package main

import (
	"strings"

	ui "github.com/gizak/termui"
	"github.com/john-pettigrew/irc/message"
)

type MessageList struct {
	ui.List
	Messages      []string
	MessageOffset int
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
		renderScreen()
	}
}

func (m *MessageList) SetMessageItems() {

	if len(m.Messages) == 0 {
		return
	}

	//calculate how many messages can fit
	var currentMessages []string
	availableHeight := m.InnerHeight()
	availableWidth := m.Width

	for i := len(m.Messages) - m.MessageOffset - 1; i >= 0; i-- {
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
}

func (m *MessageList) ListMove(y int) {
	if m.MessageOffset+y < 0 || m.MessageOffset+y > len(m.Messages)-1 {
		return
	}
	m.MessageOffset += y
	renderScreen()
}
