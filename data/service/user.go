package service

import (
	"github.com/mownier/duyog/data/store"
	"net/http"

	"github.com/gorilla/mux"
)

// User interface
type User interface {
	GetInfo() ResourceHandler
	Register() ResourceHandler
	GetPlaylists() ResourceHandler
}

// GetUserInfo method
func GetUserInfo(u User) ResourceHandler { return u.GetInfo() }

// RegisterUser method
func RegisterUser(u User) ResourceHandler { return u.Register() }

// GetUserPlaylists method
func GetUserPlaylists(u User) ResourceHandler { return u.GetPlaylists() }

type user struct {
	userRepo     store.UserRepo
	playlistRepo store.PlaylistRepo

	Resource
}

func (u user) GetInfo() ResourceHandler      { return u.getInfo }
func (u user) Register() ResourceHandler     { return u.register }
func (u user) GetPlaylists() ResourceHandler { return u.getPlaylists }

func (u user) getInfo(w http.ResponseWriter, r *http.Request) {
	auth, ok := validate(w, r, u.Resource, pathCodeUserGetInfo)

	if ok == false {
		return
	}

	key := store.UserKey(mux.Vars(r)["id"])
	user, err := store.GetUserByKey(u.userRepo, key)

	if err != nil {
		writeRespErr(w, r, u.ResponseWriter, err)
		return
	}

	if user.Key != store.UserKey(auth.UserKey) {
		user.Email = ""
	}

	marshalResponse(w, r, u.ResponseWriter, user)
}

func (u user) register(w http.ResponseWriter, r *http.Request) {
	auth, ok := validate(w, r, u.Resource, pathCodeUserRegister)

	if ok == false {
		return
	}

	param, ok := decodeUser(w, r, u.ResponseWriter)

	if ok == false {
		return
	}

	param.Key = store.UserKey(auth.UserKey)
	user, err := store.RegisterUser(u.userRepo, param)

	if err != nil {
		writeRespErr(w, r, u.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, u.ResponseWriter, user)
}

func (u user) getPlaylists(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, u.Resource, pathCodeUserRegister)

	if ok == false {
		return
	}

	key := store.UserKey(mux.Vars(r)["id"])
	playlists, err := store.GetPlaylistsByUser(u.playlistRepo, key)

	if err != nil {
		writeRespErr(w, r, u.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, u.ResponseWriter, playlists)
}

// UserResource method
func UserResource(res Resource, u store.UserRepo, p store.PlaylistRepo) User {
	return user{
		userRepo:     u,
		playlistRepo: p,

		Resource: res,
	}
}
