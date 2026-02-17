package request

import (
	"errors"
	"fmt"
	"io"
)

type Request struct {
	RequestLine RequestLine
	// Headers     map[string]string
	// Body        []byte
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, errors.Join(
			fmt.Errorf("unable to io.readALl"),
			err,
		)
	}

	str := string(data)

	rl, _, err := ParseRequestLine(str)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *rl,
	}, err
}
