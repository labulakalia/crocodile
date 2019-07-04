package util

import (
	"crypto/md5"
	"fmt"
	"hash"
	"io"
)

func GenerateAvatar(email string, size int) (avatar_url string, err error) {
	var (
		h hash.Hash
	)
	h = md5.New()
	_, err = io.WriteString(h, email)
	avatar_url = fmt.Sprintf("https://www.gravatar.com/avatar/%x?d=identicon&s=%d", h.Sum(nil), size)
	return
}
