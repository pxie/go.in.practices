package main

import (
	"fmt"
	"log"
)

type message struct {
	AgentID string
	Body    []byte

	// channel is used to send callback message to http server
	CallbackCh chan []byte
}

type broker struct {
	// register channel to add new websocket connect to broker
	register chan map[string]*agent

	// channel to get message from http server to broker
	httpCh chan *message

	// record all websocket conns
	conns map[string]*agent
}

func newBroker() *broker {
	return &broker{
		register: make(chan map[string]*agent),
		httpCh:   make(chan *message, 512),
		conns:    make(map[string]*agent),
	}
}

func (b *broker) start() {
	for {
		select {
		case item := <-b.register:
			for id, ag := range item {
				log.Printf("add agent to broker map. AgentID=%s, RemoteAddr=%s", id, ag.conn.RemoteAddr())
				b.conns[id] = ag
			}
		case msg := <-b.httpCh:
			if ag, ok := b.conns[msg.AgentID]; ok {
				log.Printf("send message to specific agent. AgentID=%s, msg=%v", msg.AgentID, msg)
				ag.agentCh <- msg
			} else {
				log.Printf("cannot passdown message to agent, since agent does not exist. AgentID=%s", msg.AgentID)
				text := []byte(fmt.Sprintf("Agent does not exist. AgentID=%s", msg.AgentID))
				msg.CallbackCh <- text
			}
		}
	}
}

func (b *broker) getConns() map[string]*agent {
	return b.conns
}
