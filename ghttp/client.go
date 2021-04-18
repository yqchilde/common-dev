package ghttp

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var client = NewClient()
var debug = false
var ReqRejectedErr = errors.New("request is rejected")

type Handler func(*Request) *Response
type Middleware func(*Client, Handler) Handler

type Client struct {
	client  *http.Client
	handler Handler
}

type RequestError struct {
	error
}

func NewClient(m ...Middleware) *Client {
	jar, _ := cookiejar.New(nil)

	c := &Client{
		client: &http.Client{
			Jar: jar,
			Transport: &http.Transport{
				Proxy: func(req *http.Request) (*url.URL, error) {
					if addr, ok := req.Context().Value("proxy").(string); ok && addr != "" {
						return url.Parse(addr)
					}
					return nil, nil
				},
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	}
	c.handler = basicHttpDo(c, nil)
	c.Use(m...)
	return c
}

func Do(req *Request) *Response {
	return client.Do(req)
}

func (c *Client) Use(m ...Middleware) *Client {
	for _, m := range m {
		c.handler = m(c, c.handler)
	}

	return c
}

func (c *Client) Do(req *Request) *Response {
	if req.Err != nil {
		return &Response{
			Req: req,
			Err: RequestError{req.Err},
		}
	}

	res := c.handler(req)
	if res == nil {
		return &Response{
			Req: req,
			Err: ReqRejectedErr,
		}
	}

	if res.Err == nil {
		res.Err = res.DecodeAndParse()
	}

	return res
}

func basicHttpDo(c *Client, next Handler) Handler {
	return func(req *Request) *Response {
		resp := &Response{
			Req:  req,
			Text: "",
			Body: []byte{},
		}

		resp.Response, resp.Err = c.client.Do(req.Request)
		if resp.Err != nil {
			return resp
		}
		defer resp.Response.Body.Close()

		resp.Body, resp.Err = ioutil.ReadAll(resp.Response.Body)
		if resp.Err != nil {
			return resp
		}
		resp.Err = resp.DecodeAndParse()
		return resp
	}
}
