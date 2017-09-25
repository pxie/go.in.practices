package main

import (
	"encoding/json"
	"log"
	"net/url"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
)

type message struct {
	ClientID string
	Body     []byte
}

func parseMessage(payload string, conn *websocket.Conn) {
	switch payload {
	case "profile":
		log.Println("retrieve local os profile.")
		doProfile(conn)
	default:
		log.Printf("unknow command. command=%s", payload)
	}
}

func doProfile(conn *websocket.Conn) {
	// log.Println("get local hardware profile")
	rtOS := runtime.GOOS

	// cpu - get CPU number of cores and speed
	cpuStat, _ := cpu.Info()
	cpuInfo, _ := json.Marshal(cpuStat)

	// host or machine kernel, uptime, platform Info
	hostStat, _ := host.Info()
	hostInfo, _ := json.Marshal(hostStat)

	html := "OS: " + rtOS + "\n"
	html = html + "CPU: " + string(cpuInfo) + "\n"
	html = html + "host: " + string(hostInfo) + "\n"

	if err := conn.WriteMessage(websocket.TextMessage, []byte(html)); err != nil {
		log.Fatalln("cannot send message back to agent.", err)
	}
}

func exec(msg string, conn *websocket.Conn) {
	switch msg {
	case "profile":
		log.Println("retrieve local os profile.")
		doProfile(conn)
	default:
		log.Printf("unknow command. command=%s", msg)
		conn.WriteMessage(websocket.TextMessage, []byte("unknown command"))
	}
}

func main() {
	// connect to websocket server

	// run against localhost
	// addr := "localhost:2222"
	// u := url.URL{Scheme: "ws", Host: addr, Path: "/register"}

	// TODO: Change addr to your websocket server, which pushed to Predix platform
	// run against websocket server on Predix
	addr := "websocket-server-multistory-width.run.aws-jp01-pr.ice.predix.io"
	u := url.URL{Scheme: "wss", Host: addr, Path: "/register"}

	log.Printf("connecting to %s", u.String())
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("connect to websocker server error.", err)
	}
	defer conn.Close()

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Fatalln("read message from websocket error.", err)
			}
			exec(string(msg), conn)
		}
	}
}
