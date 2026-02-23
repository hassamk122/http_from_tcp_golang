package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hassamk122/http_from_tcp_golang/internal/headers"
	"github.com/hassamk122/http_from_tcp_golang/internal/request"
	"github.com/hassamk122/http_from_tcp_golang/internal/response"
	"github.com/hassamk122/http_from_tcp_golang/internal/server"
)

const port = 42069

func toStr(bytes []byte) string {
	out := ""
	for _, b := range bytes {
		out += fmt.Sprintf("%02x", b)
	}
	return out
}

func serveBasicHtmlAndGetChunkedData(w *response.Writer, req *request.Request) {
	h := response.GetDefaultHeaders(0)
	body := respond200()
	status := response.StatusOK

	if req.RequestLine.RequestTarget == "/yourproblem" {
		body = respond400()
		status = response.StatusBadRequest
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		body = respond500()
		status = response.StatusInternalServerError
	} else if req.RequestLine.RequestTarget == "/video" {
		f, _ := os.ReadFile("assets/cat.mp4")
		h.Replace("content-type", "video/mp4")
		h.Replace("content-length", fmt.Sprintf("%d", len(f)))
		w.WriteStatusLine(status)
		w.WriteHeaders(*h)
		w.WriteBody(f)

	} else if req.RequestLine.RequestTarget == "/chunked" {
		fileData, err := os.ReadFile("assets/big.txt")
		if err != nil {
			body = respond500()
			status = response.StatusInternalServerError
		} else {
			w.WriteStatusLine(response.StatusOK)
			h.Delete("content-length")
			h.Set("transfer-encoding", "chunked")
			h.Set("Trailer", "X-Content-SHA256")
			h.Set("Trailer", "X-Content-Length")
			h.Replace("content-type", "text/plain")
			w.WriteHeaders(*h)

			const chunkSize = 32
			fullBody := []byte{}
			for i := 0; i < len(fileData); i += chunkSize {
				end := i + chunkSize
				if end > len(fileData) {
					end = len(fileData)
				}
				chunk := fileData[i:end]
				fullBody = append(fullBody, chunk...)
				w.WriteBody([]byte(fmt.Sprintf("%x\r\n", len(chunk))))
				w.WriteBody(chunk)
				w.WriteBody([]byte("\r\n"))
			}

			w.WriteBody([]byte("0\r\n"))
			out := sha256.Sum256(fullBody)
			trailers := headers.NewHeaders()
			trailers.Set("X-Content-SHA256", toStr(out[:]))
			trailers.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
			w.WriteHeaders(*trailers)
			return
		}
	}

	h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
	h.Replace("Content-type", "text/html")
	w.WriteStatusLine(status)
	w.WriteHeaders(*h)
	w.WriteBody(body)

}

func main() {
	s, err := server.Serve(port, serveBasicHtmlAndGetChunkedData)
	if err != nil {
		log.Fatal("Error starting server : %v", err)
	}

	defer s.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func respond200() []byte {
	return []byte(`
<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>
	`)
}

func respond400() []byte {
	return []byte(`
	<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>
	`)
}

func respond500() []byte {
	return []byte(`
<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>
	`)
}
