package request

import (
	"io"
	"strconv"

	"github.com/hassamk122/http_from_tcp_golang/internal/headers"
)

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateDone    parserState = "done"
	StateBody    parserState = "body"
	StateError   parserState = "error"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
	Headers     *headers.Headers
	Body        string
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
		Body:    "",
	}
}

func getInt(headers *headers.Headers, name string, defaultValue int) int {
	valueStr, exists := headers.Get(name)
	if !exists {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

func (r *Request) Done() bool {
	return r.state == StateDone || r.state == StateError
}

func (r *Request) Error() bool {
	return r.state == StateError
}

func (r *Request) hasBody() bool {
	length := getInt(r.Headers, "content-length", 0)
	return length > 0
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		currentData := data[read:]
		if len(currentData) == 0 {
			break
		}

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
				if r.hasBody() {
					r.state = StateBody
				} else {
					r.state = StateDone
				}
			}

			if n == 0 {
				break outer
			}
			read += n
		case StateBody:
			lengthStr := getInt(r.Headers, "content-length", 0)
			if lengthStr == 0 {
				panic("chunked not implemented")
			}

			remaining := min(lengthStr-len(r.Body), len(currentData))
			r.Body += string(currentData[:remaining])
			read += remaining

			if len(r.Body) == lengthStr {
				r.state = StateDone
			}

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
