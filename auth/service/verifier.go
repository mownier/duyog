package service

import (
	"net/http"

	"github.com/mownier/duyog/auth/store"
	"github.com/mownier/duyog/logger"
	"github.com/mownier/duyog/progerr"
	"github.com/mownier/duyog/validator"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/rpc/v2/json2"
)

// Verifier interface
type Verifier interface {
	ValidateAccessToken(*http.Request, *validator.AuthArgs, *validator.AuthReply) error
}

type verifier struct {
	userRepo   store.UserRepo
	clientRepo store.ClientRepo

	log logger.Request
}

func (v verifier) ValidateAccessToken(r *http.Request, args *validator.AuthArgs, rep *validator.AuthReply) error {
	logger.LogRequest(v.log, r)

	if r.Method != http.MethodPost {
		return &json2.Error{
			Code:    json2.E_NO_METHOD,
			Message: progerr.MethodNotAllowed.Message,
		}
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return &json2.Error{
			Code:    json2.E_INVALID_REQ,
			Message: progerr.RequestBodyNotJSON.Message,
		}
	}

	if len(args.AccessToken) == 0 {
		return &json2.Error{
			Code:    json2.E_INVALID_REQ,
			Message: progerr.TokenInvalidAccess.Message,
		}
	}

	client, err := store.GetClientByAccessToken(v.clientRepo, string(args.AccessToken))

	if err != nil {
		return &json2.Error{
			Code:    json2.E_BAD_PARAMS,
			Message: err.(progerr.Err).Message,
		}
	}

	secret := client.SecretToken
	claims := tokenClaims{}
	_, err = jwt.ParseWithClaims(string(args.AccessToken), &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return &json2.Error{
			Code:    json2.E_BAD_PARAMS,
			Message: progerr.TokenInvalidAccess.Message,
		}
	}

	user, err := store.GetUserByKey(v.userRepo, claims.UserKey)

	if err != nil {
		return &json2.Error{
			Code:    json2.E_BAD_PARAMS,
			Message: err.(progerr.Err).Message,
		}
	}

	if claims.ClientKey != client.Key {
		return &json2.Error{
			Code:    json2.E_BAD_PARAMS,
			Message: progerr.TokenInvalidAccess.Message,
		}
	}

	rep.UserKey = user.Key
	rep.Email = user.Email
	rep.Role = client.Role

	return nil
}

// NewVerifier method
func NewVerifier(u store.UserRepo, c store.ClientRepo, l logger.Request) Verifier {
	return verifier{
		userRepo:   u,
		clientRepo: c,

		log: l,
	}
}
