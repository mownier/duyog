package generator

import (
	"crypto/md5"
	"encoding/hex"
)

// Pass interface
type Pass interface {
	Generate(s string) string
}

type hashedPass struct{}

func (hashedPass) Generate(s string) string {
	hasher := md5.New()
	hasher.Write([]byte(s))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GeneratePass method
func GeneratePass(p Pass, s string) string {
	return p.Generate(s)
}

// HashedPass method
func HashedPass() Pass {
	return hashedPass{}
}
