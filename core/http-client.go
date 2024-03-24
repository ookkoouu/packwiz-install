package core

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	acceptText        = acceptMime("text/plain")
	acceptHtml        = acceptMime("text/html")
	acceptJson        = acceptMime("application/json")
	acceptOctetStream = acceptMime("application/octet-stream")
)

type acceptMime string
type httpOptFn func(c *httpClient)

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpClient struct {
	baseURL string
	headers map[string]string
	client  *http.Client
}

func NewHttpClient(opts ...httpOptFn) *httpClient {
	c := &httpClient{
		headers: map[string]string{},
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *httpClient) Do(req *http.Request) (*http.Response, error) {
	// apply baseURL
	if c.baseURL != "" {
		u, err := url.ParseRequestURI(c.baseURL + req.URL.Path)
		if err != nil {
			return nil, err
		}
		req.URL = u
	}

	// apply headers
	if c.headers != nil {
		for k, v := range c.headers {
			req.Header.Set(k, v)
		}
	}

	return c.client.Do(req)
}

func WithBaseURL(u string) httpOptFn {
	return func(c *httpClient) {
		c.baseURL = u
	}
}

func WithAccept(mime acceptMime) httpOptFn {
	return func(c *httpClient) {
		c.headers["accept"] = string(mime)
	}
}

func WithHeader(key, value string) httpOptFn {
	return func(c *httpClient) {
		key = strings.ToLower(key)
		c.headers[key] = value
	}
}

func WithTimeout(t time.Duration) httpOptFn {
	return func(c *httpClient) {
		c.client.Timeout = t
	}
}

func WithProxy(url *url.URL) httpOptFn {
	return func(c *httpClient) {
		c.client.Transport.(*http.Transport).Proxy = http.ProxyURL(url)
	}
}
