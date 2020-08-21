// go run multi_writer.go

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	reader := strings.NewReader("Here's a message for us all to read.\n")
	buffer := bytes.NewBuffer(nil)
	writer := io.MultiWriter(buffer, os.Stdout, os.Stdout)
	io.Copy(writer, reader) // copy from our reader stream and write to out multi writer

	fmt.Print("Buffer: ", buffer.String()) // buffer also contains string
}
