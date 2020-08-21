// go run filesystem_search_engine.go ../. := range

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type queryWriter struct {
	Query     []byte
	io.Writer // short-hand for Writer io.Writer?
}

func (q queryWriter) Write(input []byte) (totalBytesWritten int, err error) {
	lines := bytes.Split(input, []byte{'\n'})
	queryLength := len(q.Query)

	for _, line := range lines {
		i := bytes.Index(line, q.Query) // is the query contained in the line?
		if i == -1 {                    // no match
			continue
		}

		for _, s := range [][]byte{
			line[:i],                // what's before the match
			[]byte("\x1b[31m"),      // set star red color
			line[i : i+queryLength], // the matched part
			[]byte("\x1b[39m"),      // set default color
			line[i+queryLength:],    // whatever is remaining
		} {
			bytesWritten, err := q.Writer.Write(s)
			totalBytesWritten += bytesWritten
			if err != nil {
				return 0, err
			}
		}

		fmt.Fprintln(q.Writer) // format what is written to the Writer
	}

	return len(input), nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please specify a path and a search string.")
		return
	}

	root, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Println("Cannot get absolute path:", err)
		return
	}

	query := []byte(strings.Join(os.Args[2:], " ")) // create query from all args after path (position 1)
	fmt.Printf("Searching for %q in %s...\n", query, root)

	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		fmt.Println(path)
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		// io.TeeReader (read->write) is the oposite to io.Pipe (write->read)
		_, err = ioutil.ReadAll(io.TeeReader(file, queryWriter{Query: query, Writer: os.Stdout})) // read all of the file and "write" it to the queryWriter
		return err
	})

	if err != nil {
		fmt.Println(err)
	}
}
