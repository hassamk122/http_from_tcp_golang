package request

import (
	"io"
)

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

type Request struct {
	RequestLine RequestLine
	state       parserState
	// Headers     map[string]string
	// Body        []byte
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
outer:
	for {
		switch r.state {
		case StateInit:
		case StateDone:
			break outer
		}
	}

	return read, nil
}

func (r *Request) Done() bool {
	return r.state == StateDone
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

		readN, err := request.parse(buff[:bufLen+n])
		if err != nil {
			return nil, err
		}

		copy(buff, buff[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}
