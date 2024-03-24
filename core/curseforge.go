package core

import (
	"context"
	"fmt"
	"os"
)

var (
	cf_api_key         = ""
	DefaultCurseClient = NewCurseClient(getApiKey())
)

type cfDownloadUrlRes struct {
	Data string `json:"data"`
}

type CurseClient struct {
	apiKey     string
	httpClient HttpClient
}

func NewCurseClient(apiKey string) *CurseClient {
	c := &CurseClient{
		apiKey: apiKey,
		httpClient: NewHttpClient(
			WithBaseURL("https://api.curseforge.com"),
			WithAccept(acceptJson),
			WithHeader("user-agent", userAgent),
			WithHeader("x-api-key", apiKey),
		),
	}
	return c
}

func getApiKey() string {
	key := os.Getenv("CF_API_KEY")
	if key == "" {
		key = cf_api_key
	}
	return key
}

func (c *CurseClient) getJson(ctx context.Context, path string, v any) error {
	if c.apiKey == "" {
		return fmt.Errorf("invalid curseforge api key")
	}
	err := httpGetJson(ctx, c.httpClient, path, &v)
	if err != nil {
		return fmt.Errorf("curseforge api: %w", err)
	}
	return nil
}

func (c *CurseClient) GetDownloadUrl(ctx context.Context, d *CurseforgeData) (string, error) {
	path := fmt.Sprintf("/v1/mods/%d/files/%d/download-url", d.ProjectID, d.FileID)
	var resUrl cfDownloadUrlRes
	err := c.getJson(ctx, path, &resUrl)
	if err != nil {
		return "", err
	}
	return resUrl.Data, nil
}
