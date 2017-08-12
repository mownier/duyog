package store

import "duyog/generator"

// TokenRepo interface
type TokenRepo interface {
	Create(generator.AccessInput) (Token, error)
	Refresh(RefreshTokenInput) (Token, error)
}

// CreateToken method
func CreateToken(r TokenRepo, i generator.AccessInput) (Token, error) {
	return r.Create(i)
}

// RefreshToken method
func RefreshToken(r TokenRepo, i RefreshTokenInput) (Token, error) {
	return r.Refresh(i)
}
