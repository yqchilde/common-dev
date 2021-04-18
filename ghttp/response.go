package ghttp

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	"golang.org/x/net/html/charset"
)

type Response struct {
	*http.Response
	Body []byte
	Text string
	Req  *Request
	Err  error
}

func (r *Response) Resp() (*Response, error) {
	return r, r.Err
}

func (r *Response) Txt() (string, error) {
	return r.Text, r.Err
}

func (r *Response) HTML() (*goquery.Document, error) {
	if r.Err == nil {
		return nil, r.Err
	}

	return goquery.NewDocumentFromReader(bytes.NewReader(r.Body))
}

func (r *Response) JSON() (gjson.Result, error) {
	return gjson.Parse(r.Text), r.Err
}

func (r *Response) Error() error {
	return r.Err
}

func (r *Response) IsHTML() bool {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	return strings.Contains(contentType, "/html")
}

func (r *Response) IsJSON() bool {
	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	return strings.Contains(contentType, "/json")
}

func (r *Response) BindJSON(i interface{}) error {
	if r.Err != nil {
		return r.Err
	}

	return json.Unmarshal(r.Body, i)
}

func (r *Response) EncodeBytes(b []byte, contentType string) ([]byte, error) {
	reader, err := charset.NewReader(bytes.NewReader(b), contentType)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(reader)
}

func (r *Response) DecodeAndParse() error {
	if r.Err != nil {
		return r.Err
	}

	if len(r.Body) == 0 {
		return nil
	}

	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if strings.Contains(contentType, "text/") || strings.Contains(contentType, "/json") {
		if !strings.Contains(contentType, "charset") {
			if r.Req.RespEncode != "" {
				contentType += "; charset=" + r.Req.RespEncode
			} else {
				// TODO
			}
		}

		if strings.Contains(contentType, "utf-8") || strings.Contains(contentType, "utf8") {
			r.Text = string(r.Body)
		} else {
			encodeBytes, err := r.EncodeBytes(r.Body, contentType)
			if err != nil {
				return err
			}
			r.Body = encodeBytes
			r.Text = string(r.Body)
		}
	}

	return nil
}
