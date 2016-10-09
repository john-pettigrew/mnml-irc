package main

import (
	"strings"

	ui "github.com/gizak/termui"
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
	newServerList.Servers = []*Server{NewServer(newMsgCh, "mnmlIRC", false)}
	return &newServerList
}

func (s *ServerList) PrevWindow() {
	if len(s.Servers) == 0 {
		return
	}
	//First attempt to go to prev channel in current server.
	if !s.Servers[s.CurrentServerIndex].PrevChannel() {
		//Next try to go to last channel of prev server.
		if s.CurrentServerIndex-1 < 0 {
			s.CurrentServerIndex = len(s.Servers) - 1
		} else {
			s.CurrentServerIndex -= 1
		}

		//set correct channel in server
		s.Servers[s.CurrentServerIndex].SetChannel(len(s.Servers[s.CurrentServerIndex].Channels) - 1)
	}
}

func (s *ServerList) NextWindow() {
	if len(s.Servers) == 0 {
		return
	}
	//First attempt to go to next channel in current server.
	if !s.Servers[s.CurrentServerIndex].NextChannel() {
		//Next try to go to first channel of next server.
		if s.CurrentServerIndex+1 > len(s.Servers)-1 {
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
	newServer := NewServer(msgCh, strings.Split(serverURL, ":")[0], false)
	go newServer.ListenForMessages()

	// Switch to new server
	s.Servers = append(s.Servers, newServer)
	s.CurrentServerIndex = len(s.Servers) - 1

	msgCh <- message.Message{Command: "IRC", Options: []string{"Connecting to server..."}}
	ircConn, err = irc.NewClient(serverURL)
	if err != nil {

		msgCh <- message.Message{Command: "Error", Options: []string{"Error connecting to host: " + err.Error()}}
	}

	newServer.IrcConn = &ircConn

	msgCh <- message.Message{Command: "IRC", Options: []string{"Connected to server"}}
	go ircConn.SubscribeForMessages(&msgCh)
}

func (s *ServerList) Join(channelName string) {
	s.Servers[s.CurrentServerIndex].Join(channelName)
}

func (s ServerList) CurrentChannel() *Channel {
	return s.Servers[s.CurrentServerIndex].CurrentChannel()
}

func (s ServerList) CurrentServer() *Server {
	return s.Servers[s.CurrentServerIndex]
}

func (s *ServerList) AddMessage(m message.Message) {
	s.Servers[s.CurrentServerIndex].MsgCh <- m
}

func (s *ServerList) HandleInput(input string) {
	if input == "" {
		return
	}

	//default PRIVMSG to current channel
	if input[0] != '/' {
		input = "/PRIVMSG " + serverList.CurrentChannel().Name + " " + input
	}

	msg, err := message.ParseCommand(input)
	if err != nil {
		s.AddMessage(message.Message{Command: "Error", Options: []string{"Error parsing message: " + err.Error()}})
		return
	}

	switch msg.Command {
	case "QUIT":
		ui.StopLoop()
		return
	case "CONNECT":
		// Connect to server
		s.Connect(msg.Options[0])
	default:
		if s.Servers[s.CurrentServerIndex].IrcConn == nil {
			s.AddMessage(message.Message{Command: "Error", Options: []string{"You must connect to a server first"}})
			return
		}

		if msg.Command == "JOIN" {
			// Join Channel
			s.Servers[s.CurrentServerIndex].Join(msg.Options[0])
		}

		err = s.Servers[s.CurrentServerIndex].SendMessage(msg)
		if err != nil {
			s.AddMessage(message.Message{Command: "Error", Options: []string{"Error sending message: " + err.Error()}})
			return
		}

		s.AddMessage(msg)
	}
}
