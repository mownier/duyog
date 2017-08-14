package service

import (
	"github.com/mownier/duyog/auth/store"
	"github.com/mownier/duyog/progerr"
	"encoding/json"
	"net/http"
)

// User interface
type User interface {
	SignIn() ResourceHandler
	SignUp() ResourceHandler
	ChangePass() ResourceHandler
}

// SignInUser method
func SignInUser(u User) ResourceHandler { return u.SignIn() }

// SignUpUser method
func SignUpUser(u User) ResourceHandler { return u.SignUp() }

// ChangeUserPass method
func ChangeUserPass(u User) ResourceHandler { return u.ChangePass() }

type user struct {
	userRepo   store.UserRepo
	tokenRepo  store.TokenRepo
	clientRepo store.ClientRepo

	tokenExpiry int64

	Resource
}

func (u user) SignIn() ResourceHandler     { return u.signIn }
func (u user) SignUp() ResourceHandler     { return u.signUp }
func (u user) ChangePass() ResourceHandler { return u.changePass }

func (u user) signIn(w http.ResponseWriter, r *http.Request) {
	authenticate(w, r, u.Resource, pathCodeUserSignIn, u.clientRepo, u.tokenRepo, u.tokenExpiry,
		func(user store.User) (store.User, error) {
			return store.ValidateUserCredential(u.userRepo, user)
		})
}

func (u user) signUp(w http.ResponseWriter, r *http.Request) {
	authenticate(w, r, u.Resource, pathCodeUserSignUp, u.clientRepo, u.tokenRepo, u.tokenExpiry,
		func(user store.User) (store.User, error) {
			return store.CreateUser(u.userRepo, user)
		})
}

func (u user) changePass(w http.ResponseWriter, r *http.Request) {
	param := changePassParam{}
	err := json.NewDecoder(r.Body).Decode(&param)

	if err != nil {
		writeRespErr(w, r, u.ResponseWriter, progerr.RequestBodyNotJSON)
		return
	}

	_, ok := validate(w, r, u.Resource, pathCodeUserChangePass, u.clientRepo, param.APIkey)

	if ok == false {
		return
	}

	input := store.ChangePassInput{
		CurPass: param.CurPass,
		NewPass: param.NewPass,
		Email:   param.Email,
	}
	err = store.ChangeUserPass(u.userRepo, input)

	if err != nil {
		writeRespErr(w, r, u.ResponseWriter, err)
		return
	}

	resp := map[string]string{
		"message": "Changed password successfully",
	}

	marshalResp(w, r, u.ResponseWriter, resp)
}

// UserResource method
func UserResource(res Resource, u store.UserRepo, t store.TokenRepo, c store.ClientRepo, exp int64) User {
	return user{
		userRepo:   u,
		tokenRepo:  t,
		clientRepo: c,

		tokenExpiry: exp,

		Resource: res,
	}
}
