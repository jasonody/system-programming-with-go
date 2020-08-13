//

package main

import (
	"bufio"
	"fmt"
	"os"
)

const grr = "G. R. R. Martin"

type book struct {
	Author, Title string
	Year          int
}

func main() {
	dst, err := os.OpenFile("book_list.txt", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer dst.Close()

	bookList := []book{
		{Author: grr, Title: "A Game of Thrones", Year: 1996},
		{Author: grr, Title: "A Clash of Kings", Year: 1998},
		{Author: grr, Title: "A Storm of Swords", Year: 2000},
		{Author: grr, Title: "A Feast for Crows", Year: 2005},
		{Author: grr, Title: "A Dance with Dragons", Year: 2011},
		{Author: grr, Title: "The Winds of Winter"},
		{Author: grr, Title: "A Dream of Spring"},
	}

	w := bufio.NewWriter(dst)
	defer w.Flush()

	for _, v := range bookList {
		// prints a message with arguements to writer
		fmt.Fprintf(w, "%s - %s", v.Title, v.Author)
		if v.Year > 0 { // do not print the year of it's not present
			fmt.Fprintf(w, "(%d)", v.Year)
		}

		w.WriteRune('\n')
	}
}
