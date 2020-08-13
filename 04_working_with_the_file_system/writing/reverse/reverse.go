// go run reverse.go source.txt destination.txt

package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Please supply a read and write path.")
		return
	}

	src, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Cannot open:", err)
		return
	}
	defer src.Close()

	dst, err := os.Create(os.Args[2])
	if err != nil {
		fmt.Println("Cannot create:", err)
		return
	}
	defer dst.Close()

	cur, err := src.Seek(0, io.SeekEnd)
	if err != nil {
		fmt.Println("Error seeking:", err)
		return
	}

	b := make([]byte, 16)

	for step, r, w := int64(16), 0, 0; cur != 0; {
		if cur < step { // ensure cursor is 0 at max
			b, step = b[:cur], cur
		}
		cur = cur - step

		_, err = src.Seek(cur, io.SeekStart) // go backwards // os.SEEK_SET
		if err != nil {
			break
		}

		if r, err = src.Read(b); err != nil || r != len(b) {
			if err == nil { // all of buffer should be read
				err = fmt.Errorf("read: expected %d bytes, got %d", len(b), r)
			}
			break
		}

		// reverse
		for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
			switch { // swap "\r\n" so they are written in the correct order
			case b[i] == '\r' && b[i+1] == '\n':
				b[i], b[i+1] = b[i+1], b[i]
			case j != len(b)-1 && b[j-1] == '\r' && b[j] == '\n':
				b[j], b[j-1] = b[j-1], b[j]
			}

			b[i], b[j] = b[j], b[i] // swap bytes
		}

		if w, err = dst.Write(b); err != nil || w != len(b) {
			if err != nil {
				err = fmt.Errorf("write: expected %d bytes, got %d", len(b), w)
			}
		}
	}

	if err != nil && err != io.EOF { // we expect an EOF
		fmt.Println("\n\nError:", err)
	}
}
