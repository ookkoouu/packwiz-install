package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	neturl "net/url"
	"time"
)

var UserAgent = "packwiz/packwiz-installer"

type UserAgentTransport struct {
	http.Transport
}

func (t *UserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("user-agent", UserAgent)
	return t.Transport.RoundTrip(req)
}

func NewUATransport() *UserAgentTransport {
	trans := &UserAgentTransport{
		Transport: http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			MaxConnsPerHost:       20,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}

	return trans
}

var DefaultHttpClient = &http.Client{
	Timeout:   30 * time.Second,
	Transport: NewUATransport(),
}

type httpOptFn func(c *http.Client, r *http.Request)

func httpGet(c *http.Client, url string, opts ...httpOptFn) (*http.Response, error) {
	parsedUrl, err := neturl.ParseRequestURI(url)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("GET", parsedUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	for _, opt := range opts {
		opt(c, req)
	}

	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("response error: %s %s", res.Status, req.URL.String())
	}
	return res, nil
}

func GetJson(c *http.Client, url string, v any, opts ...httpOptFn) error {
	res, err := httpGet(c, url, WithContentType("application/json"))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(v)
}

func GetFile(c *http.Client, url string, opts ...httpOptFn) ([]byte, error) {
	res, err := httpGet(c, url, WithContentType("application/octet-stream"))
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetFileVerify(c *http.Client, url string, hashFormat string, hash string, opts ...httpOptFn) ([]byte, error) {
	data, err := GetFile(c, url, opts...)
	if err != nil {
		return nil, err
	}

	valid, err := MatchHash(data, hashFormat, hash)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, errors.New("data hash mismatched")
	}
	return data, nil
}

func WithContentType(contentType string) httpOptFn {
	return func(c *http.Client, r *http.Request) {
		r.Header.Set("content-type", contentType)
	}
}

func WithHeader(key string, value string) httpOptFn {
	return func(c *http.Client, r *http.Request) {
		r.Header.Add(key, value)
	}
}
