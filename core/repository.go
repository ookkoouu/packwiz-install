package core

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sync"

	"golang.org/x/sync/errgroup"
)

type hashObj struct {
	Hex    string
	Format string
}

type Repository struct {
	// LocalBaseDir base directory of local machine
	LocalBaseDir string
	// Url pack.toml url
	Url *url.URL
	// Pack parsed `pack.toml`
	Pack *PackToml
	// Pack parsed `index.toml`
	Index *IndexToml
	// Pack parsed `*.pw.toml`
	Metafiles []*MetafileToml
	// Mods in-memory representation of indexfiles used internally
	Mods []*Mod
	// PackHash hash of `pack.toml`. Use `WithHash`to reserve this when creation.
	PackHash   *hashObj
	httpClient *http.Client
}

type InstallStatus string

const (
	Skipped   = InstallStatus("skipped")
	Installed = InstallStatus("installed")
)

type InstallResult struct {
	Skipped   bool
	Installed bool
}

type RepoOptFn func(r *Repository)

func NewRepository(packUrl string, opts ...RepoOptFn) (*Repository, error) {
	parsedUrl, err := url.ParseRequestURI(packUrl)
	if err != nil {
		return nil, fmt.Errorf("invalid pack url: %w", err)
	}

	repo := &Repository{
		Url:        parsedUrl,
		httpClient: DefaultHttpClient,
	}

	for _, opt := range opts {
		opt(repo)
	}
	repo.LoadPack()
	return repo, nil
}

func (r *Repository) LoadPack() (*PackToml, error) {
	var (
		data []byte
		err  error
	)

	if r.PackHash == nil {
		data, err = GetFile(r.httpClient, r.Url.String())
		if err != nil {
			return nil, err
		}
	} else {
		data, err = GetFileVerify(r.httpClient, r.Url.String(), r.PackHash.Format, r.PackHash.Hex)
		if err != nil {
			return nil, err
		}
	}

	pack, err := ParsePack(data)
	if err != nil {
		return nil, err
	}
	r.Pack = pack
	return pack, nil
}

func (r *Repository) LoadIndex() (*IndexToml, error) {
	if r.Pack == nil {
		_, err := r.LoadPack()
		if err != nil {
			return nil, err
		}
	}

	indexUrl := r.BaseUrl().JoinPath(r.Pack.Index.File).String()
	data, err := GetFileVerify(r.httpClient, indexUrl, r.Pack.Index.HashFormat, r.Pack.Index.Hash)
	if err != nil {
		return nil, err
	}

	index, err := ParseIndex(data)
	if err != nil {
		return nil, err
	}
	r.Index = index
	return index, nil
}

func (r *Repository) LoadMetafiles() ([]*MetafileToml, error) {
	if r.Index == nil {
		_, err := r.LoadIndex()
		if err != nil {
			return nil, err
		}
	}
	var mods = make([]*MetafileToml, 0, len(r.Index.Files))

	eg := errgroup.Group{}
	mutex := sync.Mutex{}
	for _, file := range r.Index.Files {
		file := file
		eg.Go(func() error {
			if !file.Metafile {
				return nil
			}

			metafileUrl := r.IndexUrl().JoinPath("..", file.File).String()
			hashFmt := file.HashFormat
			if hashFmt == "" {
				hashFmt = r.Index.HashFormat
			}
			data, err := GetFileVerify(r.httpClient, metafileUrl, hashFmt, file.Hash)
			if err != nil {
				return err
			}

			mod, err := ParseMetafile(data)
			if err != nil {
				return err
			}
			mod.IndexName = file.File

			mutex.Lock()
			mods = append(mods, mod)
			mutex.Unlock()
			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	r.Metafiles = mods
	return mods, nil
}

func (r *Repository) LoadMods() ([]*Mod, error) {
	if r.Index == nil || r.Metafiles == nil {
		_, err := r.LoadMetafiles()
		if err != nil {
			return nil, err
		}
	}

	var mods = make([]*Mod, 0, len(r.Index.Files))
	for _, file := range r.Index.Files {
		if file.Metafile {
			metafile, ok := find(r.Metafiles, func(meta *MetafileToml) bool {
				return meta.IndexName == file.File
			})
			if !ok {
				return nil, fmt.Errorf("metafile not found: %s", file.File)
			}

			downloadUrl := metafile.Download.Url
			if downloadUrl == "" && metafile.Download.Mode == "metadata:curseforge" {
				cfData := metafile.Update["curseforge"]
				url, err := DefaultCurseClient.GetDownloadUrl(cfData.ProjectId, cfData.FileId)
				if err != nil {
					return nil, fmt.Errorf("curseforge fetch: %w", err)
				}
				downloadUrl = url
			}

			modDir := filepath.Join(r.IndexDir(), filepath.Dir(file.File))
			localPath := filepath.Join(modDir, metafile.Filename)

			mod := &Mod{
				LocalPath:   localPath,
				DownloadUrl: downloadUrl,
				Hash:        metafile.Download.Hash,
				HashFormat:  metafile.Download.HashFormat,
			}
			mods = append(mods, mod)
		} else {
			hashFmt := file.HashFormat
			if hashFmt == "" {
				hashFmt = r.Index.HashFormat
			}

			mod := &Mod{
				LocalPath:   filepath.Join(r.IndexDir(), file.File),
				DownloadUrl: r.IndexUrl().JoinPath("..", file.File).String(),
				Hash:        file.Hash,
				HashFormat:  hashFmt,
			}
			mods = append(mods, mod)
		}
	}

	r.Mods = mods
	return mods, nil
}

func (r *Repository) CheckIntegrity(m *Mod) (bool, error) {
	// existence
	p := filepath.Join(r.LocalBaseDir, m.LocalPath)
	stat, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		if stat.IsDir() {
			return false, fmt.Errorf("")
		}
		return false, err
	}

	// hash
	file, err := os.Open(p)
	if err != nil {
		return false, fmt.Errorf("open file: %w", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return false, fmt.Errorf("read file: %w", err)
	}
	valid, err := MatchHash(data, m.HashFormat, m.Hash)
	if err != nil {
		return false, fmt.Errorf("calc hash: %w", err)
	}

	return valid, nil
}

func (r *Repository) WriteFileMod(m *Mod) error {
	_, err := filepath.Abs(r.LocalBaseDir)
	if r.LocalBaseDir == "" || err != nil {
		return fmt.Errorf("BaseDir required")
	}

	p := filepath.Join(r.LocalBaseDir, m.LocalPath)
	err = os.MkdirAll(filepath.Dir(p), 0755)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}
	file, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0744)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer func() error {
		err := file.Close()
		if err != nil {
			fmt.Println("file close:", err)
		}
		return err
	}()
	_, err = file.Write(m.Data)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func (r *Repository) Install(ctx context.Context) error {
	if r.Mods == nil {
		_, err := r.LoadMods()
		if err != nil {
			return fmt.Errorf("loading mods: %w", err)
		}
	}

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(runtime.NumCPU())
	for _, m := range r.Mods {
		m := m
		eg.Go(func() error {
			fmt.Println("Attempt:", m.LocalPath)
			valid, err := r.CheckIntegrity(m)
			if err != nil {
				return err
			}
			fmt.Println("Hash:", valid, m.LocalPath)
			if valid {
				fmt.Println("Skipped:", filepath.Join(r.LocalBaseDir, m.LocalPath))
				return nil
			}

			err = m.Download(ctx, r.httpClient)
			if err != nil {
				return err
			}

			err = r.WriteFileMod(m)
			if err != nil {
				return err
			}
			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) IndexDir() string {
	return filepath.Dir(r.Pack.Index.File)
}

func (r *Repository) BaseUrl() *url.URL {
	return r.Url.JoinPath("..")
}

func (r *Repository) IndexUrl() *url.URL {
	return r.BaseUrl().JoinPath(r.Pack.Index.File)
}

func WithHash(hashFormat string, hash string) RepoOptFn {
	return func(r *Repository) {
		if hashFormat == "" || hash == "" {
			return
		}
		if !slices.Contains(PreferredHashList, hashFormat) {
			return
		}

		h := &hashObj{
			Hex:    hash,
			Format: hashFormat,
		}
		r.PackHash = h
	}
}

func WithClient(c *http.Client) RepoOptFn {
	return func(r *Repository) {
		r.httpClient = c
	}
}

func WithDir(localDir string) RepoOptFn {
	return func(r *Repository) {
		dir, err := filepath.Abs(localDir)
		if err != nil {
			return
		}
		r.LocalBaseDir = dir
	}
}

func WithProxy(host string) RepoOptFn {
	return func(r *Repository) {
		u, err := url.Parse(host)
		if host == "" || err != nil {
			return
		}
		r.httpClient.Transport.(*http.Transport).Proxy = http.ProxyURL(u)
	}
}
