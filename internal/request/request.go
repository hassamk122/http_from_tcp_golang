package request

import (
	"fmt"
	"io"
)

/*
GET /something HTTP/1.1
*/
type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     map[string]string
	Body        []byte
}

var ERR_BAD_START_LINE = fmt.Errorf("bad start line")

func ParseRequestLine(s string) (*Request, string, error) {
	// will do later
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	// will do later
}
