package generator

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// AccessInput struct
type AccessInput struct {
	UserKey     string
	ClientKey   string
	SecretToken string
	Expiry      int64
}

// TokenOutput struct
type TokenOutput struct {
	Token     string
	CreatedOn int64
}

// TokenClaims data struct
type tokenClaims struct {
	UserKey   string `json:"user_id"`
	ClientKey string `json:"client_id"`
	jwt.StandardClaims
}

// Token interface
type Token interface {
	Generate() (TokenOutput, error)
}

// Access interface
type Access interface {
	Generate(i AccessInput) (TokenOutput, error)
}

type access struct{}

func (g access) Generate(i AccessInput) (TokenOutput, error) {
	var output TokenOutput

	now := time.Now()

	claims := tokenClaims{
		i.UserKey,
		i.ClientKey,
		jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(time.Second * time.Duration(i.Expiry)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(i.SecretToken))

	if err != nil {
		return output, err
	}

	output = TokenOutput{
		Token:     token,
		CreatedOn: now.Unix(),
	}

	return output, nil
}

type token struct {
	length int
}

func (g token) Generate() (TokenOutput, error) {
	var output TokenOutput

	token, err := g.generateRandomString()

	if err != nil {
		return output, err
	}

	output = TokenOutput{
		Token:     token,
		CreatedOn: time.Now().Unix(),
	}

	return output, nil
}

func (g token) generateRandomString() (string, error) {
	bytes, err := g.generateRandomBytes()
	return base64.URLEncoding.EncodeToString(bytes), err
}

func (g token) generateRandomBytes() ([]byte, error) {
	bytes := make([]byte, g.length)
	_, err := rand.Read(bytes)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

// GenerateToken method
func GenerateToken(g Token) (TokenOutput, error) { return g.Generate() }

// GenerateAccess method
func GenerateAccess(a Access, i AccessInput) (TokenOutput, error) { return a.Generate(i) }

// AccessToken method
func AccessToken() Access { return access{} }

// RefreshToken method
func RefreshToken(l int) Token { return token{length: l} }

// SecretToken method
func SecretToken(l int) Token { return token{length: l} }

// APIToken method
func APIToken(l int) Token { return token{length: l} }
