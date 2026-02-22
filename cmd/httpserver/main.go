package main

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/html") {
		target := req.RequestLine.RequestTarget
		resp, err := http.Get("http://httpbin.org/" + target[len("/httpbin/"):])
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
			fullBody := []byte{}
			for {
				data := make([]byte, 32)
				n, err := resp.Body.Read(data)
				if err != nil {
					break
				}
				fullBody = append(fullBody, data[:n]...)
				w.WriteBody([]byte(fmt.Sprintf("%x\r\n", n)))
				w.WriteBody(data[:n])
				w.WriteBody([]byte("\r\n"))
			}
			w.WriteBody([]byte("0\r\n"))

			out := sha256.Sum256(fullBody)
			tailers := headers.NewHeaders()
			tailers.Set("X-Content-SHA256", toStr(out[:]))
			tailers.Set("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))
			w.WriteHeaders(*tailers)
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
