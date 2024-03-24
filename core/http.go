package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	userAgent         = "packwiz-install"
	DefaultHttpClient = NewHttpClient(WithHeader("user-agent", userAgent))
)

func httpGetJson(ctx context.Context, c HttpClient, url string, v any) error {
	req, err := http.NewRequestWithContext(context.WithoutCancel(ctx), "GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("accept", string(acceptJson))
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(v)
}

func httpGetBytes(ctx context.Context, c HttpClient, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(context.WithoutCancel(ctx), "GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("accept", string(acceptOctetStream))

	res, err := c.Do(req)
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

func httpGetValidBytes(ctx context.Context, c HttpClient, url string, hashFormat string, hash string) ([]byte, error) {
	data, err := httpGetBytes(context.WithoutCancel(ctx), c, url)
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
