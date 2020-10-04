// go run child_processes.go

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Accessing child processes
	cmd := exec.Command("ls", "-l")
	// cmd.Stdout = os.Stdout
	var out outstream
	cmd.Stdout = out
	if err := cmd.Start(); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("Cmd: ", cmd.Args[0])
	fmt.Println("Args: ", cmd.Args[1])
	fmt.Println("PID: ", cmd.Process.Pid)
	cmd.Wait()

	// Standard input
	b := bytes.NewBuffer(nil)
	cmd = exec.Command("cat")
	cmd.Stdin = b
	cmd.Stdout = os.Stdout
	fmt.Fprintf(b, "Hello world! I'm using memory address: %p\n", b)
	if err := cmd.Start(); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	cmd.Wait()
}

type outstream struct{}

func (out outstream) Write(p []byte) (int, error) {
	fmt.Println(string(p))
	return len(p), nil
}
