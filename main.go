package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("open error : ", err)
	}

	str := ""
	for {
		data := make([]byte, 8)

		n, err := f.Read(data)
		if err != nil {
			break
		}

		data = data[:n]

		idxNewline := bytes.IndexByte(data, '\n')
		if idxNewline != -1 {
			str += string(data[:idxNewline])
			data = data[idxNewline+1:]
			fmt.Printf("read : %s\n", str)
			str = ""
		}

		str += string(data)
	}

	if len(str) != 0 {
		fmt.Printf("read : %s\n", str)
	}
}
