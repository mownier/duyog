package validator

import (
	"bytes"
	"duyog/progerr"
	"net/http"

	"github.com/gorilla/rpc/v2/json2"
)

// AuthURL string
type AuthURL string

// AuthMethod string
type AuthMethod string

// AccessToken string
type AccessToken string

// AuthArgs data
type AuthArgs struct {
	AccessToken AccessToken `json:"access_token"`
}

// AuthReply struct
type AuthReply struct {
	UserKey string `json:"user_id"`
	Email   string `json:"email"`
	Role    string `json:"role"`
}

// Auth interface
type Auth interface {
	Validate(t AccessToken) (AuthReply, error)
}

// ValidateAuth method
func ValidateAuth(a Auth, t AccessToken) (AuthReply, error) {
	return a.Validate(t)
}

type rpcAuth struct {
	url    string
	method string
}

func (a rpcAuth) Validate(t AccessToken) (AuthReply, error) {
	var reply AuthReply

	if t == "" {
		return reply, progerr.TokenInvalidAccess
	}

	args := AuthArgs{
		AccessToken: t,
	}
	json, err := json2.EncodeClientRequest(a.method, args)

	if err != nil {
		return reply, progerr.Internal(err)
	}

	req, err := http.NewRequest(http.MethodPost, a.url, bytes.NewBuffer(json))

	if err != nil {
		return reply, progerr.Internal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return reply, progerr.TokenInvalidAccess
	}

	var tmp AuthReply
	err = json2.DecodeClientResponse(resp.Body, &tmp)

	if err != nil {
		return reply, progerr.TokenInvalidAccess
	}

	reply = tmp

	return reply, nil
}

// RPCAuth method
func RPCAuth(u AuthURL, m AuthMethod) Auth {
	return rpcAuth{
		url:    string(u),
		method: string(m),
	}
}
