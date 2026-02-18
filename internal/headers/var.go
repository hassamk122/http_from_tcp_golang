package headers

import "fmt"

var SEPARATOR = []byte("\r\n")
var ERR_BAD_HEADER_FIELD_LINE = fmt.Errorf("bad header field line")
var ERR_BAD_HEADER_FIELD_NAME = fmt.Errorf("bad header field name")
