// go build exit_cleanup.go
// run application in the background
// ./exit_cleanup &
// kill the appliation (PID will be logged when application is run in the background)
// kill -6 <PID>
// verify data was written to file.txt

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c := make(chan os.Signal)
	// note: kill command issues a SIGABRT signal
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT) // notify on termination, interrupt, and abort signal

	f, err := os.OpenFile("file.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	// this go routine will ensure whenever is written to the buffer gets persisted
	// even if the application is terminated before it completes
	go func() {
		<-c // wait for signal on channel c
		w.Flush()
		log.Println("Application terminated gracefully")
		os.Exit(0)
	}()

	for i := 0; i <= 10; i++ {
		fmt.Fprintln(w, "hello")
		// log.Println(i)
		time.Sleep(time.Second)
	}

	log.Println("Application exited normally")
}
