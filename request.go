package reqs2http

import (
	"bytes"
	"net/http"

	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
)

type RequestConverter interface {
	Convert(*rhp.Request) (*http.Request, error)
}

type RequestConvFn func(*rhp.Request) (*http.Request, error)

func (f RequestConvFn) Convert(p *rhp.Request) (*http.Request, error) { return f(p) }
func (f RequestConvFn) AsIf() RequestConverter                        { return f }

func RequestConvNew(methodDefault rhp.Method) RequestConverter {
	return RequestConvFn(func(p *rhp.Request) (*http.Request, error) {
		var m rhp.Method = p.GetMethod()
		var method Method = Method{m}.Or(methodDefault)
		var ms string = method.String()
		var body []byte = p.GetBody()
		var rdr *bytes.Reader = bytes.NewReader(body)
		return http.NewRequest(
			ms,
			p.GetUrl(),
			rdr,
		)
	})
}

var RequestConvDefault RequestConverter = RequestConvNew(rhp.Method_METHOD_GET)
