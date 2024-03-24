package core

import (
	"fmt"
	"strings"

	packwiz "github.com/packwiz/packwiz/core"
)

var PreferredHashList = []string{
	"murmur2",
	"md5",
	"sha1",
	"sha256",
	"sha512",
}

func MatchHash(data []byte, hashFormat string, hash string) (bool, error) {
	hasher, err := packwiz.GetHashImpl(hashFormat)
	if err != nil {
		return false, err
	}
	_, err = hasher.Write(data)
	if err != nil {
		return false, err
	}
	hashgot := fmt.Sprintf("%x", hasher.Sum(nil))
	if !strings.EqualFold(hash, hashgot) {
		return false, nil
	}
	return true, nil
}
