// go build custom_signal_handling.go
// run application in the background
// ./custom_signal_handling &
// kill the appliation (PID will be logged when application is run in the background)
// kill -6 <PID>

package main

import (
	"log"
	"os"
	"os/signal"
)

func main() {
	log.Println("Start application...")
	c := make(chan os.Signal)
	signal.Notify(c) // notify on any signal
	s := <-c         // wait for signal (this should be done in separate goroutine)
	log.Println("Exit with signal:", s)
}
