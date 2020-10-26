// go run anonymous_pipes.go

package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

type readerWriter struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

func main() {
	var cmds = []*exec.Cmd{
		exec.Command("cat", "book_list.txt"),
		exec.Command("grep", "Game"),
		exec.Command("wc", "-l"),
	}

	var pipes []readerWriter

	for i := range cmds {
		if i == len(cmds)-1 {
			pipes = append(pipes, readerWriter{reader: nil, writer: nil})
			cmds[i].Stdout = os.Stdout
		} else {
			r, w := io.Pipe()
			pipes = append(pipes, readerWriter{reader: r, writer: w})
			cmds[i+1].Stdin, cmds[i].Stdout = pipes[i].reader, pipes[i].writer
		}

		if err := cmds[i].Start(); err != nil {
			log.Fatalln("Start:", i, err)
		}
	}

	for i, pipe := range pipes {
		if err := cmds[i].Wait(); err != nil {
			log.Fatalln("Wait:", i, err)
		}

		if pipe.writer == nil {
			continue
		}

		if err := pipe.writer.Close(); err != nil {
			log.Fatalln("Close:", i, err)
		}
	}
}
