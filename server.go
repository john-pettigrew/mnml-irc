package main

import (
	"errors"
	"strings"

	"github.com/john-pettigrew/irc"
	"github.com/john-pettigrew/irc/message"
)

type Server struct {
	Channels            []*Channel
	CurrentChannelIndex int
	MsgCh               chan message.Message
	Closable            bool
	IrcConn             *irc.Client
}

func NewServer(msgCh chan message.Message, closable bool) *Server {
	newServer := Server{MsgCh: msgCh, Closable: closable}
	newServer.Channels = []*Channel{NewChannel("", false)}
	go newServer.ListenForMessages()
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

		//send to correct channel(s)
		sendToAll := true
		if newMsg.Options[0][0] == '#' {
			//send to a certain channel
			sendToAll = false
		}
		for i := 0; i < len(s.Channels); i++ {
			if sendToAll || s.Channels[i].Name == newMsg.Options[0] {
				s.Channels[i].Messages = append(s.Channels[i].Messages, msgStr)
			}
		}

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

func (s *Server) Join(channelName string) {
	s.Channels = append(s.Channels, NewChannel(channelName, true))
	s.CurrentChannelIndex = len(s.Channels) - 1
}

func (s *Server) SendMessage(msg message.Message) error {
	if s.IrcConn == nil {
		return errors.New("Error: Not connected to server.")
	}
	return s.IrcConn.SendMessage(msg)
}
