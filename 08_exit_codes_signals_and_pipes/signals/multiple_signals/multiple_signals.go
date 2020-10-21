// go build multiple_signals.go
// shell #1:
// ./multiple_signals
// shell #2:
// increase duration: kill -SIGUSR1 $(pgrep multi)
// decrease duration: kill -SIGUSR2 $(pgrep multi)
// load settings: kill -SIGHUP $(pgrep multi)
// save settings: kill -SIGALRM $(pgrep multi)
// save settings and exit: kill -SIGINT $(pgrep multi)
// exit with saving settings: kill -SIGQUIT $(pgrep multi)

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"syscall"
	"time"
)

var cfgPath string

func init() {
	u, err := user.Current()
	if err != nil {
		log.Fatal("user: ", err)
	}
	cfgPath = filepath.Join(u.HomeDir, ".multi")
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2, syscall.SIGALRM)

	// initial load
	d := time.Second * 4
	if err := handleSignals(syscall.SIGHUP, &d); err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	for {
		select {
		case s := <-c:
			if err := handleSignals(s, &d); err != nil {
				log.Printf("Error handling %s: %s", s, err)
				continue
			}
		default:
			time.Sleep(d)
			log.Println("After", d, "Executing action!")
		}
	}
}

func handleSignals(s os.Signal, d *time.Duration) error {
	switch s {
	case syscall.SIGHUP:
		return loadSettings(d)
	case syscall.SIGALRM:
		return saveSetting(d)
	case syscall.SIGINT:
		if err := saveSetting(d); err != nil {
			log.Println("Cannot save:", err)
		}
		fallthrough
	case syscall.SIGQUIT:
		os.Exit(1)
	case syscall.SIGUSR1:
		changeSettings(d, (*d)*2)
		return nil
	case syscall.SIGUSR2:
		changeSettings(d, (*d)/2)
		return nil
	}
	return nil
}

func changeSettings(d *time.Duration, v time.Duration) {
	*d = v
	log.Println("Changed:", v)
}

func loadSettings(d *time.Duration) error {
	b, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		return err
	}

	var v time.Duration
	if v, err = time.ParseDuration(string(b)); err != nil {
		return err
	}

	*d = v
	log.Println("Loaded: ", v)
	return nil
}

func saveSetting(d *time.Duration) error {
	f, err := os.OpenFile(cfgPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := fmt.Fprintf(f, d.String()); err != nil {
		return err
	}

	log.Println("Saved:", *d)
	return nil
}
