package core

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var cfApiKey = "$2a$10$6BqncsODxEQp0JyQ.sExoeF5C44DL2Bunh2iyJ2jnDUeoV/hGyaya"

type CurseClient struct {
	httpClient *http.Client
	apiKey     string
	host       *url.URL
}

type CurseOptFn func(c *CurseClient)

type resDownloadUrl struct {
	Data string `json:"data"`
}

func NewCurseClient(apiKey string, opts ...CurseOptFn) *CurseClient {
	u, _ := url.ParseRequestURI("https://api.curseforge.com")
	c := &CurseClient{
		httpClient: DefaultHttpClient,
		apiKey:     apiKey,
		host:       u,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

var DefaultCurseClient = NewCurseClient(cfApiKey)

func (c *CurseClient) get(path string) ([]byte, error) {
	res, err := httpGet(c.httpClient, c.host.JoinPath(path).String(), WithHeader("x-api-key", c.apiKey))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	return data, nil
}

func (c *CurseClient) getJson(path string, v any) error {
	data, err := c.get(path)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, v)
	if err != nil {
		return fmt.Errorf("decode json: %w", err)
	}
	return nil
}

func (c *CurseClient) GetDownloadUrl(projId uint32, fileId uint32) (string, error) {
	path := fmt.Sprintf("/v1/mods/%d/files/%d/download-url", projId, fileId)
	resUrl := new(resDownloadUrl)
	if err := c.getJson(path, resUrl); err != nil {
		return "", err
	}
	return resUrl.Data, nil
}

func WithHost(host string) CurseOptFn {
	u, err := url.ParseRequestURI(host)
	return func(c *CurseClient) {
		if err == nil {
			c.host = u
		}
	}
}
