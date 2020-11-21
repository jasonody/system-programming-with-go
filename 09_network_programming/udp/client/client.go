// go run client.go localhost:8900 <<< $'hello there\nPDU htiw nuf emos'

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/jasonody/system-programming-with-go/09_network_programming/udp/common"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalln("Please specify an address")
	}

	addr, err := net.ResolveUDPAddr("udp", os.Args[1])
	if err != nil {
		log.Fatalln("Invalid address:", os.Args[1], err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalln("-> Connection:", err)
	}

	log.Println("-> Connection to", addr)
	r := bufio.NewReader(os.Stdin)
	b := make([]byte, 256*256)

	for {
		fmt.Print("# ")
		msg, err := r.ReadBytes('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("-> Message error", err)
			}
			return
		}

		data, err := common.CreateMessage(bytes.TrimSpace(msg))
		if err != nil {
			log.Fatalln("-> Encode error:", err)
		}

		if _, err := conn.Write(data); err != nil {
			log.Println("-> Connection:", err)
			return
		}

		n, err := conn.Read(b)
		if err != nil {
			log.Println("<- Receive error:", err)
			continue
		}

		msg, err = common.MessageContent(b[:n])
		if err != nil {
			log.Println("<- Decpde error:", err)
			continue
		}
		log.Printf("<- %q", msg)
	}
}
