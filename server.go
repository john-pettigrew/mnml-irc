package main

import (
	"strings"

	"github.com/john-pettigrew/irc/message"
)

type Server struct {
	Channels            []*Channel
	CurrentChannelIndex int
	MsgCh               chan message.Message
	Closable            bool
}

func NewServer(msgCh chan message.Message, closable bool) *Server {
	newServer := Server{MsgCh: msgCh, Closable: closable}
	newServer.Channels = []*Channel{NewChannel(false)}
	return &newServer
}

func (s *Server) ListenForMessages() {
	var newMsg message.Message
	var open bool
	for {
		newMsg, open = <-s.MsgCh
		if !open {
			break
		}

		msgStr := message.Marshal(newMsg)

		//remove extra characters
		msgStr = strings.Replace(msgStr, "\n", "", -1)
		msgStr = strings.Replace(msgStr, "\r", "", -1)
		s.Channels[s.CurrentChannelIndex].Messages = append(s.Channels[s.CurrentChannelIndex].Messages, msgStr)
		renderScreen()
	}
}

func (s *Server) NextChannel() bool {
	if s.CurrentChannelIndex+1 > len(s.Channels)-1 {
		return false
	}

	s.CurrentChannelIndex += 1
	return true
}

func (s *Server) PrevChannel() bool {
	if s.CurrentChannelIndex-1 < 0 {
		return false
	}

	s.CurrentChannelIndex -= 1
	return true
}

func (s *Server) SetChannel(newChannel int) bool {
	if len(s.Channels)-1 < newChannel {
		return false
	}

	s.CurrentChannelIndex = newChannel
	return true
}

func (s *Server) CurrentChannel() *Channel {
	return s.Channels[s.CurrentChannelIndex]
}
