package core

import "github.com/pelletier/go-toml/v2"

type PackToml struct {
	Name        string `toml:"name"`
	Author      string `toml:"author,omitempty"`
	Version     string `toml:"version,omitempty"`
	Description string `toml:"description,omitempty"`
	PackFormat  string `toml:"pack-format"`
	Index       struct {
		File       string `toml:"file"`
		HashFormat string `toml:"hash-format"`
		Hash       string `toml:"hash"`
	} `toml:"index"`
	Versions map[string]string `toml:"versions"`
}

type IndexToml struct {
	HashFormat string            `toml:"hash-format"`
	Files      []IndexedfileToml `toml:"files,omitempty"`
}

type IndexedfileToml struct {
	File       string `toml:"file"`
	Hash       string `toml:"hash"`
	Alias      string `toml:"alias,omitempty"`
	HashFormat string `toml:"hash-format,omitempty"`
	Metafile   bool   `toml:"metafile,omitempty"`
	Preserve   bool   `toml:"preserve,omitempty"`
}

type UpdateModrinth struct {
	ModId   string `toml:"mod-id,omitempty"`
	Version string `toml:"version,omitempty"`
}

type UpdateCurseForge struct {
	FileId    uint32 `toml:"file-id,omitempty"`
	ProjectId uint32 `toml:"project-id,omitempty"`
}

type MetafileToml struct {
	Filename string `toml:"filename"`
	Name     string `toml:"name"`
	Side     string `toml:"side,omitempty"`
	Download struct {
		HashFormat string `toml:"hash-format"`
		Hash       string `toml:"hash"`
		Url        string `toml:"url,omitempty"`
		Mode       string `toml:"mode,omitempty"`
	} `toml:"download"`
	Update *struct {
		Modrinth   *UpdateModrinth   `toml:"modrinth,omitempty,nullable"`
		CurseForge *UpdateCurseForge `toml:"curseforge,omitempty,nullable"`
	} `toml:"update,omitempty"`
	Option *struct {
		Optional    bool   `toml:"optional"`
		Default     bool   `toml:"default,omitempty"`
		Description string `toml:"description,omitempty"`
	} `toml:"option,omitempty"`
	IndexName string
}

func ParsePack(data []byte) (*PackToml, error) {
	var pack = new(PackToml)
	if err := toml.Unmarshal(data, pack); err != nil {
		return nil, err
	}
	return pack, nil
}

func ParseIndex(data []byte) (*IndexToml, error) {
	var index = new(IndexToml)
	if err := toml.Unmarshal(data, index); err != nil {
		return nil, err
	}
	return index, nil
}

func ParseMetafile(data []byte) (*MetafileToml, error) {
	var mod = new(MetafileToml)
	if err := toml.Unmarshal(data, mod); err != nil {
		return nil, err
	}
	return mod, nil
}
