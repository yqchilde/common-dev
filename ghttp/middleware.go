package ghttp

import (
	"fmt"
	"net/http"
	"net/url"
)

func WithCookie(rawUrl string, cookies ...*http.Cookie) Middleware  {
	return func(c *Client, h Handler) Handler {
		u, err := url.Parse(rawUrl)
		if err != nil {
			fmt.Println("add cookie to jar failed, err: ", err.Error())
		} else {
			c.client.Jar.SetCookies(u, cookies)
		}

		return func(req *Request) *Response {
			return h(req)
		}
	}
}

func WithDebug() Middleware {
	return func(c *Client, h Handler) Handler {
		return func(req *Request) *Response {
			res := h(req.SetDebug(true))
			return res
		}
	}
}