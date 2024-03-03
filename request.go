package reqs2http

import (
	"bytes"
	"errors"
	"net/http"

	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
)

var (
	ErrUnexpectedHeader error = errors.New("unexpected header")
)

const ContentTypeDefault string = "application/octet-stream"

type RequestConverter interface {
	Convert(*rhp.Request) (*http.Request, error)
}

type RequestConvFn func(*rhp.Request) (*http.Request, error)

func (f RequestConvFn) Convert(p *rhp.Request) (*http.Request, error) {
	return f(p)
}

func (f RequestConvFn) AsIf() RequestConverter { return f }

func RequestConvNew(methodDefault rhp.Method) RequestConverter {
	return RequestConvFn(func(p *rhp.Request) (*http.Request, error) {
		var m rhp.Method = p.GetMethod()
		var method Method = Method{m}.Or(methodDefault)
		var ms string = method.String()
		var body []byte = p.GetBody()
		var rdr *bytes.Reader = bytes.NewReader(body)
		req, e := http.NewRequest(
			ms,
			p.GetUrl(),
			rdr,
		)
		if nil != e {
			return nil, e
		}

		req.Header.Set("Content-Type", ContentTypeDefault) // default type

		var hdr *rhp.Header = p.GetHeader()
		var items []*rhp.HeaderItem = hdr.GetItems()
		for _, item := range items {
			switch v := item.GetItem().(type) {
			case *rhp.HeaderItem_Custom:
				req.Header.Add(
					v.Custom.GetKey(),
					v.Custom.GetVal(),
				)
			case *rhp.HeaderItem_ContentType:
				req.Header.Set("Content-Type", v.ContentType)
			default:
				return nil, ErrUnexpectedHeader
			}
		}
		return req, e
	})
}

var RequestConvDefault RequestConverter = RequestConvNew(rhp.Method_METHOD_GET)
