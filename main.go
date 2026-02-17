package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("open error : ", err)
	}
	for {
		data := make([]byte, 8)
		n, err := f.Read(data)
		if err != nil {
			if err == io.EOF {
				fmt.Println("file ended")
			}
			break
		}
		fmt.Printf("read %s\n", string(data[:n]))
	}
}
