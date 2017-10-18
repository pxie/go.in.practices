package main

import "time"
import "log"

func main() {
	// ticker := time.NewTicker(2 * time.Second)
	quit := make(chan bool)

	// define time out function first
	timer := time.NewTimer(5 * time.Second)
	go func() {
		<-timer.C
		log.Println("send message to quit channel.")
		quit <- true
		close(quit)
	}()

	// run endless loop
	go func() {
		i := 0
		for {
			select {
			case <-quit:
				log.Print("time is up, quit.")
				return
			default:
				// endless loop
				log.Println(i)
				i++
				time.Sleep(time.Second)
			}
		}
	}()

	time.Sleep(15 * time.Second)
	log.Println("end game")
}
