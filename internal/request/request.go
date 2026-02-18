package request

import (
	"io"
)

type parserState string

const (
	StateInit  parserState = "init"
	StateDone  parserState = "done"
	StateError parserState = "error"
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
		switch r.state {
		case StateError:
			return 0, ERR_REQUEST_IN_ERROR_STATE
		case StateInit:
			rl, n, err := ParseRequestLine(data[read:])
			if err != nil {
				r.state = StateError
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n

			r.state = StateDone
		case StateDone:
			break outer
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

		readN, err := request.parse(buff[:bufLen+n])
		if err != nil {
			return nil, err
		}

		copy(buff, buff[readN:bufLen])
		bufLen -= readN
	}

	return request, nil
}
