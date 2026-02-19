package request

import (
	"io"

	"github.com/hassamk122/http_from_tcp_golang/internal/headers"
)

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateDone    parserState = "done"
	StateError   parserState = "error"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
	Headers     *headers.Headers
	// Body        []byte
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

func (r *Request) Done() bool {
	return r.state == StateDone || r.state == StateError
}

func (r *Request) Error() bool {
	return r.state == StateError
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentData := data[read:]
		switch r.state {
		case StateError:
			return 0, ERR_REQUEST_IN_ERROR_STATE
		case StateInit:
			rl, n, err := ParseRequestLine(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n

			r.state = StateHeaders
		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}

			if done {
				r.state = StateDone
			}

			if n == 0 {
				break outer
			}

			read += n
		case StateDone:
			break outer
		default:
			panic("something went wrong")
		}
	}

	return read, nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buff := make([]byte, 1024)
	bufLen := 0
	for !request.Done() {
		n, err := reader.Read(buff[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n

		readN, err := request.parse(buff[:bufLen])
		if err != nil {
			return nil, err
		}

		copy(buff, buff[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}
