package extractor

import (
	"duyog/validator"
	"strings"
)

// Auth interface
type Auth interface {
	Extract(s string) validator.AccessToken
}

type auth struct{}

func (auth) Extract(s string) validator.AccessToken {
	if s == "" {
		return ""
	}

	words := strings.Split(s, " ")

	if len(words) != 2 {
		return ""
	}

	if strings.ToLower(words[0]) != "bearer" {
		return ""
	}

	if len(words[1]) <= 0 {
		return ""
	}

	return validator.AccessToken(words[1])
}

// ExtractAuth method
func ExtractAuth(a Auth, s string) validator.AccessToken {
	return a.Extract(s)
}

// AccessToken method
func AccessToken() Auth {
	return auth{}
}
