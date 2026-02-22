package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/hassamk122/http_from_tcp_golang/internal/request"
	"github.com/hassamk122/http_from_tcp_golang/internal/response"
	"github.com/hassamk122/http_from_tcp_golang/internal/server"
)

const port = 42069

func serveBasicHtml(w *response.Writer, req *request.Request) {
	h := response.GetDefaultHeaders(0)
	body := respond200()
	status := response.StatusOK

	if req.RequestLine.RequestTarget == "/yourproblem" {
		body = respond400()
		status = response.StatusBadRequest
	} else if req.RequestLine.RequestTarget == "/myproblem" {
		body = respond500()
		status = response.StatusInternalServerError
	}

	h.Replace("Content-length", fmt.Sprintf("%d", len(body)))
	h.Replace("Content-type", "text/html")
	w.WriteStatusLine(status)
	w.WriteHeaders(*h)
	w.WriteBody(body)

}

func main() {
	s, err := server.Serve(port, serveBasicHtml)
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
