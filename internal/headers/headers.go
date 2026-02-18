package headers

import "bytes"

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", ERR_BAD_HEADER_FIELD_LINE
	}
	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", ERR_BAD_HEADER_FIELD_NAME
	}

	return string(name), string(value), nil
}

/*
Remember
field-name ":" OWS field-value OWS
*/

func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], SEPARATOR)
		if idx == -1 {
			break
		}

		if idx == 0 {
			done = true
			read += len(SEPARATOR)
			break
		}

		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}

		read += idx + len(SEPARATOR)

		h[name] = value
	}

	return read, done, nil
}
