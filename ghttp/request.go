package ghttp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/context"
)

type Request struct {
	*http.Request
	RespEncode string
	Writer     io.Writer
	Debug      bool
	Err        error

	client   *Client
	callback func(resp *Response) *Response
}

func NewRequest(method, url string) *Request {
	req, err := http.NewRequest(method, url, nil)
	return &Request{
		Request:    req,
		RespEncode: "",
		Debug:      debug,
		Err:        err,
		client:     client,
		callback: func(resp *Response) *Response {
			return resp
		},
	}
}

func Get(url string) *Request {
	return NewRequest("GET", url)
}

func Post(url string) *Request {
	return NewRequest("POST", url)
}

func Put(url string) *Request {
	return NewRequest("PUT", url)
}

func Delete(url string) *Request {
	return NewRequest("DELETE", url)
}

func Head(url string) *Request {
	return NewRequest("HEAD", url)
}

func Connect(url string) *Request {
	return NewRequest("CONNECT", url)
}

func Options(url string) *Request {
	return NewRequest("OPTIONS", url)
}

func Trace(url string) *Request {
	return NewRequest("TRACE", url)
}

func Patch(url string) *Request {
	return NewRequest("PATCH", url)
}

func (req *Request) addContextValue(k, v interface{}) *Request {
	req.Request = req.WithContext(context.WithValue(req.Request.Context(), k, v))
	return req
}

func (req *Request) SetDebug(d bool) *Request {
	req.Debug = d
	return req
}

func (req *Request) SetTimeout(t time.Duration) *Request {
	ctx, _ := context.WithTimeout(req.Context(), t)
	req.Request = req.WithContext(ctx)
	return req
}

func (req *Request) AddCookie(c *http.Cookie) *Request {
	req.Request.AddCookie(c)
	return req
}

func (req *Request) AddCookies(cs ...*http.Cookie) *Request {
	for _, c := range cs {
		req.Request.AddCookie(c)
	}
	return req
}

func (req *Request) AddHeader(k, v string) *Request {
	req.Request.Header.Add(k, v)
	return req
}

func (req *Request) AddHeaders(m map[string]string) *Request {
	for k, v := range m {
		req.AddHeader(k, v)
	}
	return req
}

func (req *Request) AddParam(k, v string) *Request {
	if len(req.Request.URL.RawQuery) > 0 {
		req.Request.URL.RawQuery += "&"
	}
	req.Request.URL.RawQuery += url.QueryEscape(k) + "=" + url.QueryEscape(v)
	return req
}

func (req *Request) AddParams(m map[string]string) *Request {
	for k, v := range m {
		req.AddParam(k, v)
	}
	return req
}

func (req *Request) SetUA(ua string) *Request {
	req.AddHeader("User-Agent", ua)
	return req
}

func (req *Request) SetBasicAuth(username, password string) *Request {
	req.Request.SetBasicAuth(username, password)
	return req
}

func (req *Request) SetBody(body io.Reader) *Request {
	rc, ok := body.(io.ReadCloser)
	if !ok && body != nil {
		rc = io.NopCloser(body)
	}
	req.Request.Body = rc

	switch v := body.(type) {
	case *bytes.Buffer:
		req.ContentLength = int64(v.Len())
		buf := v.Bytes()
		req.GetBody = func() (io.ReadCloser, error) {
			r := bytes.NewReader(buf)
			return io.NopCloser(r), nil
		}
	case *bytes.Reader:
		req.ContentLength = int64(v.Len())
		snapshot := *v
		req.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return io.NopCloser(&r), nil
		}
	case *strings.Reader:
		req.ContentLength = int64(v.Len())
		snapshot := *v
		req.GetBody = func() (io.ReadCloser, error) {
			r := snapshot
			return io.NopCloser(&r), nil
		}
	default:
	}

	if req.GetBody != nil && req.ContentLength == 0 {
		req.Body = http.NoBody
		req.GetBody = func() (io.ReadCloser, error) {
			return http.NoBody, nil
		}
	}

	return req
}

func (req *Request) SetRawBody(b []byte) *Request {
	req.SetBody(bytes.NewReader(b))
	return req
}

func (req *Request) SetFormBody(m map[string]string) *Request {
	var u url.URL
	var q = u.Query()

	for k, v := range m {
		q.Add(k, v)
	}
	req.SetRawBody([]byte(q.Encode()))
	req.AddHeader("Content-Type", "application/x-www-form-urlencoded")

	return req
}

func (req *Request) SetJsonBody(v interface{}) *Request {
	body, err := json.Marshal(v)
	req.SetRawBody(body)
	req.Err = err
	req.AddHeader("Content-Type", "application/json")

	return req
}

func (req *Request) SetCallback(f func(resp *Response) *Response) *Request {
	req.callback = f
	return req
}

func (req *Request) SetClient(c *Client) *Request {
	req.client = c
	return req
}

func (req *Request) String() string {
	return req.URL.String()
}

func (req *Request) Do() *Response {
	return req.callback(req.client.Do(req))
}
