package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
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

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("open error : ", err)
	}

	lines := getLinesChannel(f)
	for line := range lines {
		fmt.Println(line)
	}
}
