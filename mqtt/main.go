package main

import (
	"time"
)

func main() {
	broker := "tcp://localhost:1883"
	topic := "demo/qos2"
	qos := byte(2)

	subcriber := &client{}
	subcriber.init(broker, "subClient")
	defer subcriber.distroy()

	subcriber.sub(topic, qos, handler)
	defer subcriber.unsub(topic)

	publisher := &client{}
	clientID := "pubClient"
	publisher.init(broker, clientID)
	defer publisher.distroy()

	publisher.pub(topic, qos, createMsg(clientID, 0))

	time.Sleep(1 * time.Second)
}
