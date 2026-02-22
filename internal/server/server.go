package server

import (
	"fmt"
	"io"
	"net"

	"github.com/hassamk122/http_from_tcp_golang/internal/request"
	"github.com/hassamk122/http_from_tcp_golang/internal/response"
)

type Handler func(w *response.Writer, req *request.Request)

type Server struct {
	closed  bool
	handler Handler
}

func (s *Server) handle(conn io.ReadWriteCloser) {
	defer conn.Close()

	responseWriter := response.NewWriter(conn)

	r, err := request.RequestFromReader(conn)
	if err != nil {
		responseWriter.WriteStatusLine(response.StatusBadRequest)
		responseWriter.WriteHeaders(*response.GetDefaultHeaders(0))
		return
	}

	s.handler(responseWriter, r)
}

func listen(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if s.closed {
			return
		}
		if err != nil {
			return
		}

		go s.handle(conn)

	}
}

func Serve(port uint16, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	server := &Server{
		closed:  false,
		handler: handler,
	}
	go listen(server, listener)
	return server, nil
}

func (s *Server) Close() error {
	s.closed = true
	return nil
}
