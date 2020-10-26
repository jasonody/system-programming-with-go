// go run broadcast.go

package main

import (
	"bytes"
	"io"
	"log"
	"os/exec"
)

func main() {
	var (
		words   = []string{"Game", "Feast", "Dragons", "of"}
		cmds    = make([][2]*exec.Cmd, len(words))
		writers = make([]io.Writer, len(words))
		buffers = make([]bytes.Buffer, len(words))
		err     error
	)

	for i := range words {
		cmds[i][0] = exec.Command("grep", words[i])
		// collect all the greps' stdin pipes in the writers slices
		// (these will be connected to cat's stdout using the MultiWriter)
		if writers[i], err = cmds[i][0].StdinPipe(); err != nil {
			log.Fatal("in pipe:", i, err)
		}

		cmds[i][1] = exec.Command("wc", "-l")
		// connect each wc's stdin to the corresponing grep's stdout pipe
		if cmds[i][1].Stdin, err = cmds[i][0].StdoutPipe(); err != nil {
			log.Fatal("in pipe:", i, err)
		}
		// connect each wc's stdout to a buffer
		cmds[i][1].Stdout = &buffers[i]
	}

	cat := exec.Command("cat", "book_list.txt")
	// broadcast cat's stdout to each of the grep commands
	cat.Stdout = io.MultiWriter(writers...)
	for i := range cmds {
		for j := range cmds[i] {
			if err := cmds[i][j].Start(); err != nil {
				log.Fatalln("start:", i, j, err)
			}
		}
	}
	if err := cat.Run(); err != nil {
		log.Fatalln("cat:", err)
	}
	// close the set of pipes writing to grep so it can terminate
	// when the StdinPipe method is used, the resulting writer must be closed
	for i := range cmds {
		if err := writers[i].(io.Closer).Close(); err != nil {
			log.Fatalln("close 0:", i, err)
		}
	}

	// wait on grep to finish
	for i := range cmds {
		if err := cmds[i][0].Wait(); err != nil {
			log.Fatalln("grep wait:", i, err)
		}
	}

	// wait on wc to finish and then log output
	for i := range cmds {
		if err := cmds[i][1].Wait(); err != nil {
			log.Fatalln("wc wait:", i, err)
		}
		count := bytes.TrimSpace(buffers[i].Bytes())
		log.Printf("%10q %s entries", cmds[i][0].Args[1], count)
	}
}
