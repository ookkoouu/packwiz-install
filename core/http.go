package core

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/carlmjohnson/requests"
)

var (
	userAgent             = "packwiz-install"
	defaultRequestBuilder = newRequestBuilder(http.DefaultClient)
)

func newRequestBuilder(c *http.Client) *requests.Builder {
	return requests.New().
		Client(c).
		UserAgent(userAgent)
}

func httpGetJson(ctx context.Context, c *http.Client, url string, v any) error {
	return defaultRequestBuilder.Clone().Client(c).BaseURL(url).ToJSON(&v).Fetch(ctx)
}

func httpGetBytes(ctx context.Context, c *http.Client, url string) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := defaultRequestBuilder.
		Clone().
		Client(c).
		BaseURL(url).
		ToBytesBuffer(buf).
		Fetch(context.WithoutCancel(ctx))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func httpGetValidBytes(ctx context.Context, c *http.Client, url string, hashFormat string, hash string) ([]byte, error) {
	data, err := httpGetBytes(ctx, c, url)
	if err != nil {
		return nil, err
	}

	valid, err := MatchHash(data, hashFormat, hash)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, fmt.Errorf("download hash mismatched: %s", url)
	}
	return data, nil
}
