package main

import (
	"log"
	"net/http"

	"github.com/satori/go.uuid"

	"github.com/gorilla/websocket"
)

type agent struct {
	// broker to know every agent
	b *broker

	// connect to remote client via websocket
	conn *websocket.Conn

	// channel get message from broker
	agentCh chan *message
}

func (a *agent) start() {
	for {
		select {
		case msg := <-a.agentCh:
			// send message to remote client
			log.Printf("send message to remote client. message=%v", msg)
			if err := a.conn.WriteMessage(websocket.TextMessage, msg.Body); err != nil {
				log.Fatalln("write message back to websocket error.", err)
			}

			// wait for the response sent back from client
			_, payload, err := a.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Fatalln("websocket close error.", err)
				}
			}
			log.Printf("Get response payload from remote client. payload=%s", payload)

			// client cannot handle specific message.
			if string(payload) == "unknown command" {
				log.Println("Get `unknown command` from client side.")
				payload = []byte("Server Internel Error")
			}
			msg.CallbackCh <- payload
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serve websocket request send from remote clients
func serveWs(b *broker, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("init websocket conn error.", err)
	}

	ag := &agent{b: b, conn: conn, agentCh: make(chan *message, 256)}
	agentID := uuid.NewV4().String()

	// build map to register agent
	reg := make(map[string]*agent)
	reg[agentID] = ag
	log.Printf("register agent to broker. agentID=%s, agent=%v", agentID, ag)
	b.register <- reg

	// one goroutine to handle communication between agent to remote client, and broker
	go ag.start()
}
