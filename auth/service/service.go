package service

import (
	"github.com/mownier/duyog/auth/store"
	"github.com/mownier/duyog/generator"
	"github.com/mownier/duyog/progerr"
	"github.com/mownier/duyog/validator"
	"github.com/mownier/duyog/writer"
	"encoding/json"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
)

// ResourceHandler function
type ResourceHandler func(http.ResponseWriter, *http.Request)

// RPCHandler function
type RPCHandler func(*http.Request, *validator.AuthArgs, *validator.AuthReply) error

// Resource struct
type Resource struct {
	ResponseWriter writer.Response
}

type authParam struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	APIKey   string `json:"api_key"`
}

type authResponse struct {
	RefreshToken string `json:"refresh_token,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	Email        string `json:"email,omitempty"`
	UserKey      string `json:"user_id,omitempty"`
	Expiry       int64  `json:"expiry,omitempty"`
}

type changePassParam struct {
	Email   string `json:"email"`
	NewPass string `json:"new_password"`
	CurPass string `json:"current_password"`
	APIkey  string `json:"api_key"`
}

type refreshTokenParam struct {
	RefreshToken string `json:"refresh_token"`
	Expiry       int64  `json:"expiry,omitempty"`
	APIKey       string `json:"api_key"`
}

type tokenClaims struct {
	UserKey   string `json:"user_id"`
	ClientKey string `json:"client_id"`
	jwt.StandardClaims
}

func validate(w http.ResponseWriter, r *http.Request, res Resource, c pathCode, clientRepo store.ClientRepo, apiKey string) (store.Client, bool) {
	var client store.Client
	var err error

	if r.Header.Get("Content-Type") != "application/json" {
		writeRespErr(w, r, res.ResponseWriter, progerr.ContentTypeNotJSON)
		return client, false
	}

	cl, err := store.GetClientByAPIKey(clientRepo, apiKey)

	if err != nil {
		writeRespErr(w, r, res.ResponseWriter, err)
		return client, false
	}

	if isForbidden(c) && cl.Role != "admin" {
		writeRespErr(w, r, res.ResponseWriter, progerr.RequestPathForbidden)
		return client, false
	}

	client = cl

	return client, true
}

func authenticate(w http.ResponseWriter, r *http.Request, res Resource, c pathCode, clientRepo store.ClientRepo, tokenRepo store.TokenRepo, expiry int64, userGen func(store.User) (store.User, error)) bool {
	var err error

	if r.Header.Get("Content-Type") != "application/json" {
		writeRespErr(w, r, res.ResponseWriter, progerr.ContentTypeNotJSON)
		return false
	}

	param := authParam{}
	err = json.NewDecoder(r.Body).Decode(&param)

	if err != nil {
		writeRespErr(w, r, res.ResponseWriter, progerr.RequestBodyNotJSON)
		return false
	}

	client, ok := validate(w, r, res, c, clientRepo, param.APIKey)

	if ok == false {
		return false
	}

	user := store.User{
		Email:    param.Email,
		Password: param.Password,
	}
	user, err = userGen(user)

	if err != nil {
		writeRespErr(w, r, res.ResponseWriter, err)
		return false
	}

	input := generator.AccessInput{
		ClientKey:   client.Key,
		UserKey:     user.Key,
		SecretToken: client.SecretToken,
		Expiry:      expiry,
	}

	token, err := tokenRepo.Create(input)

	if err != nil {
		writeRespErr(w, r, res.ResponseWriter, err)
		return false
	}

	resp := authResponse{
		AccessToken:  token.Access,
		RefreshToken: token.Refresh,
		Expiry:       token.Expiry,

		UserKey: user.Key,
		Email:   user.Email,
	}

	return marshalResp(w, r, res.ResponseWriter, resp)
}

func writeRespErr(w http.ResponseWriter, r *http.Request, wr writer.Response, e error) {
	var err progerr.Err

	switch e.(type) {
	case progerr.Err:
		err = e.(progerr.Err)

	default:
		err = progerr.Internal(err)
	}

	writer.WriteResponse(wr, w, r, err.HTTPStatus, err.Data())
}

func marshalResp(w http.ResponseWriter, r *http.Request, wr writer.Response, v interface{}) bool {
	data, err := json.Marshal(v)

	if err != nil {
		writeRespErr(w, r, wr, err)
		return false
	}

	writer.WriteResponse(wr, w, r, http.StatusOK, data)

	return true
}
