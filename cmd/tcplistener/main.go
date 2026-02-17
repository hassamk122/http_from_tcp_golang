package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)
		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)

			if n > 0 {

				data = data[:n]

				idxNewline := bytes.IndexByte(data, '\n')
				if idxNewline != -1 {
					str += string(data[:idxNewline])
					data = data[idxNewline+1:]
					out <- str
					str = ""
				}

				str += string(data)
			}

			if err == io.EOF {
				break
			}

			if err != nil {
				return
			}
		}

		if len(str) != 0 {
			out <- str
		}
	}()

	return out
}

// wrote this for practice
func getBytesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer f.Close()
		defer close(out)

		for {
			data := make([]byte, 8)

			n, err := f.Read(data)
			if err != nil {
				break
			}

			data = data[:n]
			out <- string(data)
		}

	}()

	return out
}

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}
		for line := range getLinesChannel(conn) {
			fmt.Printf("read : %s \n", line)
		}
	}

}
