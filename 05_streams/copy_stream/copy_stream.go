// go run copy_stream.go

package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func CopyNOffset(dst io.Writer, src io.ReadSeeker, offset, length int64) (int64, error) {
	if _, err := src.Seek(offset, io.SeekStart); err != nil {
		return 0, err
	}

	return io.CopyN(dst, src, length)
}

func main() {
	src := strings.NewReader("This is an example of CopyN with offset")
	for i, length, step := int64(0), int64(src.Len()), int64(5); i < length; i += step {
		CopyNOffset(os.Stdout, src, i, step)
		fmt.Println()
	}
}
