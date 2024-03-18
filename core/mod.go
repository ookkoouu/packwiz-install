package core

import (
	"context"
	"net/http"
	"runtime"

	"golang.org/x/sync/errgroup"
)

type Mod struct {
	// LocalPath relative path from base dir
	LocalPath   string `json:"localPath"`
	Hash        string `json:"hash"`
	HashFormat  string `json:"hashFormat"`
	DownloadUrl string `json:"downloadUrl"`
	Data        []byte
}

func (m *Mod) Download(ctx context.Context, c *http.Client) error {
	data, err := GetFileVerify(c, m.DownloadUrl, m.HashFormat, m.Hash)
	if err != nil {
		return err
	}

	m.Data = data
	return nil
}

type ModList []*Mod

func (mods ModList) DownloadModAll(ctx context.Context, c *http.Client) error {
	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(runtime.NumCPU())
	for _, mod := range mods {
		mod := mod
		eg.Go(func() error {
			return mod.Download(ctx, c)
		})
	}

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}
