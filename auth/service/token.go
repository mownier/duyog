package service

import (
	"github.com/mownier/duyog/auth/store"
	"github.com/mownier/duyog/progerr"
	"encoding/json"
	"net/http"
)

// Token interface
type Token interface {
	Refresh() ResourceHandler
}

// RefreshToken method
func RefreshToken(t Token) ResourceHandler { return t.Refresh() }

type token struct {
	userRepo   store.UserRepo
	tokenRepo  store.TokenRepo
	clientRepo store.ClientRepo

	Resource
}

func (t token) Refresh() ResourceHandler { return t.refresh }

func (t token) refresh(w http.ResponseWriter, r *http.Request) {
	param := refreshTokenParam{}
	err := json.NewDecoder(r.Body).Decode(&param)

	if err != nil {
		writeRespErr(w, r, t.ResponseWriter, progerr.RequestBodyNotJSON)
		return
	}

	_, ok := validate(w, r, t.Resource, pathCodeTokenRefresh, t.clientRepo, param.APIKey)

	if ok == false {
		return
	}

	input := store.RefreshTokenInput{
		Expiry:       param.Expiry,
		RefreshToken: param.RefreshToken,
	}
	token, err := store.RefreshToken(t.tokenRepo, input)

	if err != nil {
		writeRespErr(w, r, t.ResponseWriter, err)
		return
	}

	resp := authResponse{
		Expiry:       token.Expiry,
		AccessToken:  token.Access,
		RefreshToken: token.Refresh,
	}

	marshalResp(w, r, t.ResponseWriter, resp)
}

// TokenResource method
func TokenResource(res Resource, u store.UserRepo, t store.TokenRepo, c store.ClientRepo) Token {
	return token{
		userRepo:   u,
		tokenRepo:  t,
		clientRepo: c,

		Resource: res,
	}
}
