package reqs2http_test

import (
	"testing"

	"bytes"
	"io"
	"net/http"
	"net/url"

	r2h "github.com/takanoriyanagitani/go-reqs2http"
	rhp "github.com/takanoriyanagitani/go-reqs2http/reqs2http/v1"
)

func TestRequest(t *testing.T) {
	t.Parallel()

	t.Run("RequestConverter", func(t *testing.T) {
		t.Parallel()

		t.Run("RequestConvDefault", func(t *testing.T) {
			t.Parallel()

			var conv r2h.RequestConverter = r2h.RequestConvDefault

			t.Run("Method", func(t *testing.T) {
				t.Parallel()

				t.Run("GET", func(t *testing.T) {
					t.Parallel()

					req := &rhp.Request{Method: rhp.Method_METHOD_GET}
					var converted *http.Request = must(conv.Convert(req))
					var method string = converted.Method
					t.Run("same method", assertEqual(method, http.MethodGet))
				})
			})

			t.Run("Url", func(t *testing.T) {
				t.Parallel()

				t.Run("example", func(t *testing.T) {
					t.Parallel()

					const testUrl string = "https://example.com"
					req := &rhp.Request{Url: testUrl}
					var converted *http.Request = must(conv.Convert(req))
					var u *url.URL = converted.URL
					t.Run("same url", assertEqual(u.String(), testUrl))
				})
			})

			t.Run("Body", func(t *testing.T) {
				t.Parallel()

				t.Run("empty", func(t *testing.T) {
					t.Parallel()

					req := &rhp.Request{}
					var converted *http.Request = must(conv.Convert(req))
					var buf bytes.Buffer
					_, _ = io.Copy(&buf, converted.Body)
					t.Run("empty bytes", assertEqual(len(buf.Bytes()), 0))
				})

				t.Run("string", func(t *testing.T) {
					t.Parallel()

					req := &rhp.Request{Body: []byte("helo")}
					var converted *http.Request = must(conv.Convert(req))
					var buf bytes.Buffer
					_, _ = io.Copy(&buf, converted.Body)
					t.Run("same bytes", assertEqual(buf.String(), "helo"))
				})
			})

			t.Run("Header", func(t *testing.T) {
				t.Parallel()

				t.Run("Content-Type", func(t *testing.T) {
					t.Parallel()

					t.Run("empty", func(t *testing.T) {
						t.Parallel()

						req := &rhp.Request{}
						var converted *http.Request = must(conv.Convert(req))
						var hdr http.Header = converted.Header
						var typ string = hdr.Get("Content-Type")
						const otyp string = "application/octet-stream"
						t.Run("default type", assertEqual(typ, otyp))
					})

					t.Run("json", func(t *testing.T) {
						t.Parallel()

						req := &rhp.Request{
							Header: &rhp.Header{
								Items: []*rhp.HeaderItem{
									{
										Item: &rhp.HeaderItem_ContentType{
											ContentType: "application/json",
										},
									},
								},
							},
						}
						var converted *http.Request = must(conv.Convert(req))
						var hdr http.Header = converted.Header
						var typ string = hdr.Get("Content-Type")
						const jtyp string = "application/json"
						t.Run("default type", assertEqual(typ, jtyp))
					})
				})

				t.Run("User-Agent", func(t *testing.T) {
					t.Parallel()

					t.Run("empty", func(t *testing.T) {
						t.Parallel()

						req := &rhp.Request{}
						var converted *http.Request = must(conv.Convert(req))
						var hdr http.Header = converted.Header
						var agent string = hdr.Get("User-Agent")
						t.Run("no agent", assertEqual(agent, ""))
					})

					t.Run("curl", func(t *testing.T) {
						t.Parallel()

						req := &rhp.Request{
							Header: &rhp.Header{
								Items: []*rhp.HeaderItem{
									{
										Item: &rhp.HeaderItem_Custom{
											Custom: &rhp.CustomHeader{
												Key: "user-agent",
												Val: "curl/8.4.0",
											},
										},
									},
								},
							},
						}
						var converted *http.Request = must(conv.Convert(req))
						var hdr http.Header = converted.Header
						var agent string = hdr.Get("User-Agent")
						t.Run("curl", assertEqual(agent, "curl/8.4.0"))
					})
				})
			})
		})
	})
}
