package ghttp

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
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

func (r *Request) addContextValue(k, v interface{}) *Request {
	r.Request = r.WithContext(context.WithValue(r.Request.Context(), k, v))
	return r
}

func (r *Request) SetDebug(d bool) *Request {
	r.Debug = d
	return r
}

func (r *Request) SetTimeout(t time.Duration) *Request {
	ctx, _ := context.WithTimeout(r.Context(), t)
	r.Request = r.WithContext(ctx)
	return r
}

func (r *Request) AddCookie(c *http.Cookie) *Request {
	r.Request.AddCookie(c)
	return r
}

func (r *Request) AddCookies(cs ...*http.Cookie) *Request {
	for _, c := range cs {
		r.Request.AddCookie(c)
	}
	return r
}

func (r *Request) AddHeader(k, v string) *Request {
	r.Request.Header.Add(k, v)
	return r
}

func (r *Request) AddHeaders(m map[string]string) *Request {
	for k, v := range m {
		r.AddHeader(k, v)
	}
	return r
}

func (r *Request) AddParam(k, v string) *Request {
	if len(r.Request.URL.RawQuery) > 0 {
		r.Request.URL.RawQuery += "&"
	}
	r.Request.URL.RawQuery += url.QueryEscape(k) + "=" + url.QueryEscape(v)
	return r
}

func (r *Request) AddParams(m map[string]string) *Request {
	for k, v := range m {
		r.AddParam(k, v)
	}
	return r
}

func (r *Request) SetUA(ua string) *Request {
	r.AddHeader("User-Agent", ua)
	return r
}

func (r *Request) SetBasicAuth(username, password string) *Request {
	r.Request.SetBasicAuth(username, password)
	return r
}

func (r *Request) SetBody(b io.Reader) *Request {
	rc, ok := b.(io.ReadCloser)
	if !ok && b != nil {
		rc = ioutil.NopCloser(b)
	}
	r.Request.Body = rc

	switch v := b.(type) {
	case *bytes.Buffer:
		r.ContentLength = int64(v.Len())
		buf := v.Bytes()
		r.GetBody = func() (io.ReadCloser, error) {
			reader := bytes.NewReader(buf)
			return ioutil.NopCloser(reader), nil
		}
	case *bytes.Reader:
		r.ContentLength = int64(v.Len())
		r.GetBody = func() (io.ReadCloser, error) {
			return ioutil.NopCloser(v), nil
		}
	case *strings.Reader:
		r.ContentLength = int64(v.Len())
		r.GetBody = func() (io.ReadCloser, error) {
			return ioutil.NopCloser(v), nil
		}
	default:
	}

	return r
}

func (r *Request) SetRawBody(b []byte) *Request {
	r.SetBody(bytes.NewReader(b))
	return r
}

func (r *Request) SetFormBody(m map[string]string) *Request {
	var u url.URL
	var q = u.Query()

	for k, v := range m {
		q.Add(k, v)
	}
	r.SetRawBody([]byte(q.Encode()))
	r.AddHeader("Content-Type", "application/x-www-form-urlencoded")

	return r
}

func (r *Request) SetJSONBody(v interface{}) *Request {
	body, err := json.Marshal(v)
	r.SetRawBody(body)
	r.Err = err
	r.AddHeader("Content-Type", "application/json")

	return r
}

func (r *Request) SetCallback(f func(resp *Response) *Response) *Request {
	r.callback = f
	return r
}

func (r *Request) SetClient(c *Client) *Request {
	r.client = c
	return r
}

func (r *Request) String() string {
	return r.URL.String()
}

func (r *Request) Do() *Response {
	return r.callback(r.client.Do(r))
}
