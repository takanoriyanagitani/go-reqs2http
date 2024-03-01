package reqs2http

import (
	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
)

type Method struct {
	rhp.Method
}

var m2s map[rhp.Method]string = map[rhp.Method]string{
	rhp.Method_METHOD_UNSPECIFIED: "GET",
	rhp.Method_METHOD_GET:         "GET",
	rhp.Method_METHOD_HEAD:        "HEAD",
	rhp.Method_METHOD_POST:        "POST",
	rhp.Method_METHOD_PUT:         "PUT",
	rhp.Method_METHOD_DELETE:      "DELETE",
	rhp.Method_METHOD_CONNECT:     "CONNECT",
	rhp.Method_METHOD_OPTIONS:     "OPTIONS",
	rhp.Method_METHOD_TRACE:       "TRACE",
	rhp.Method_METHOD_PATCH:       "PATCH",
}

func init() {
}

func (m Method) String() string { return m2s[m.Method] }
func (m Method) Or(alt rhp.Method) Method {
	switch m.Method {
	case rhp.Method_METHOD_UNSPECIFIED:
		return Method{alt}
	default:
		return m
	}
}
