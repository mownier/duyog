package generator

import "github.com/rs/xid"

// Key interface
type Key interface {
	Generate() string
}

type xidKeyGenerator struct{}

func (xidKeyGenerator) Generate() string {
	return xid.New().String()
}

// GenerateKey generates key
func GenerateKey(k Key) string { return k.Generate() }

// XIDKey method
func XIDKey() Key { return xidKeyGenerator{} }
