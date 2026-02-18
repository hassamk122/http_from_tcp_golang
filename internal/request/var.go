package request

import "fmt"

var ERR_BAD_REQUESTLINE = fmt.Errorf("bad request line")
var ERR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
var ERR_REQUEST_IN_ERROR_STATE = fmt.Errorf("request in error state")
var SEPARATOR = []byte("\r\n")
