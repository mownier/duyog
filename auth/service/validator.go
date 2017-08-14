package service

import (
	"github.com/mownier/duyog/auth/store"
	"github.com/mownier/duyog/logger"
	"github.com/mownier/duyog/progerr"
	"github.com/mownier/duyog/validator"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/rpc/v2/json2"
)

// Validator interface
type Validator interface {
	ValidateAccessToken() RPCHandler
}

// ValidateAccessToken method
func ValidateAccessToken(v Validator) RPCHandler { return v.ValidateAccessToken() }

type authValidator struct {
	userRepo   store.UserRepo
	clientRepo store.ClientRepo

	log logger.Request
}

func (a authValidator) ValidateAccessToken() RPCHandler { return a.validateAccessToken }

func (a authValidator) validateAccessToken(r *http.Request, args *validator.AuthArgs, rep *validator.AuthReply) error {
	logger.LogRequest(a.log, r)

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

	client, err := store.GetClientByAccessToken(a.clientRepo, string(args.AccessToken))

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

	user, err := store.GetUserByKey(a.userRepo, claims.UserKey)

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

// AuthValidator method
func AuthValidator(u store.UserRepo, c store.ClientRepo, l logger.Request) Validator {
	return authValidator{
		userRepo:   u,
		clientRepo: c,

		log: l,
	}
}
