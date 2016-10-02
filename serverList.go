package main

import (
	"github.com/john-pettigrew/irc"
	"github.com/john-pettigrew/irc/message"
)

type ServerList struct {
	Servers            []*Server
	CurrentServerIndex int
}

func NewServerList() *ServerList {
	newMsgCh := make(chan message.Message)
	newServerList := ServerList{}
	newServerList.Servers = []*Server{NewServer(newMsgCh, false)}
	return &newServerList
}

func (s *ServerList) NextWindow() {
	//First attempt to go to next channel in current server.
	if !s.Servers[s.CurrentServerIndex].NextChannel() {
		//Next try to go to first channel of next server.
		if s.CurrentServerIndex+1 > len(s.Servers)-1 {
			if len(s.Servers) == 0 {
				return
			}
			s.CurrentServerIndex = 0
		} else {
			s.CurrentServerIndex += 1
		}

		//set correct channel in server
		s.Servers[s.CurrentServerIndex].SetChannel(0)
	}
}

func (s *ServerList) Connect(serverURL string) {

	// Add new server with base channel
	msgCh := make(chan message.Message)
	newServer := NewServer(msgCh, false)
	go newServer.ListenForMessages()

	// Switch to new server
	s.Servers = append(s.Servers, newServer)
	s.CurrentServerIndex = len(s.Servers) - 1

	msgCh <- message.Message{Command: "IRC", Options: []string{"Connecting to server..."}}
	ircConn, err = irc.NewClient(serverURL)
	if err != nil {

		msgCh <- message.Message{Command: "Error", Options: []string{"Error connecting to host: " + err.Error()}}
	}

	msgCh <- message.Message{Command: "IRC", Options: []string{"Connected to server"}}
	go ircConn.SubscribeForMessages(&msgCh)
}

func (s ServerList) CurrentChannel() *Channel {
	return s.Servers[s.CurrentServerIndex].CurrentChannel()
}

func (s *ServerList) AddMessage(m message.Message) {
	s.Servers[s.CurrentServerIndex].MsgCh <- m
}
