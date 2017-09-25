package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cloudfoundry-community/go-cfenv"
)

func agentsHandler(b *broker, w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var response []byte
	if val, ok := r.Form["id"]; ok {
		// has id parameter, therefore send message to broker to parse parameter
		msg := &message{}
		msg.AgentID = strings.Join(val, "")
		msg.Body = []byte(profileCmd)
		msg.CallbackCh = make(chan []byte)
		defer close(msg.CallbackCh)

		log.Printf("send message to broker. message=%v", msg)
		b.httpCh <- msg

		// wait for callback message as response
		response = <-msg.CallbackCh
	} else {
		// No id parameter, return all agent list
		log.Println("list all agents from REST interface.")
		response = listAgents(b)
	}
	fmt.Fprintln(w, string(response))
}

func listAgents(b *broker) []byte {

	type record struct {
		AgentID    string
		RemoteAddr string
	}
	var text []record

	for id, agent := range b.conns {
		r := &record{}
		r.AgentID = id
		r.RemoteAddr = agent.conn.RemoteAddr().String()
		text = append(text, *r)
	}
	if len(text) == 0 {
		return []byte("no agent has been registered.")
	}
	res, _ := json.Marshal(text)
	return res
}

func getPort() string {
	env, err := cfenv.Current()
	if err != nil {
		return "2222"
	}

	return strconv.Itoa(env.Port)
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "welcome!\n For more details, https://github.com/pxie/go.in.practices/tree/master/websocket")
}

func main() {
	broker := newBroker()
	go broker.start()

	http.HandleFunc("/", home)

	http.HandleFunc("/agents", func(w http.ResponseWriter, r *http.Request) {
		agentsHandler(broker, w, r)
	})

	// create agent, and register it to broker
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		serveWs(broker, w, r)
	})

	addr := ":" + getPort()
	log.Println("start listen", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal("http ListenAndServe error.", err)
	}
}
