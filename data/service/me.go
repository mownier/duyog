package service

import (
	"duyog/data/store"
	"net/http"

	"github.com/gorilla/mux"
)

// Me interface
type Me interface {
	Update() ResourceHandler
	AddPlaylists() ResourceHandler
	GetPlaylists() ResourceHandler
	UpdatePlaylist() ResourceHandler
	AddPlaylistSongs() ResourceHandler
}

// UpdateMe method
func UpdateMe(m Me) ResourceHandler { return m.Update() }

// AddToMyPlaylists method
func AddToMyPlaylists(m Me) ResourceHandler { return m.AddPlaylists() }

// GetMyPlaylists method
func GetMyPlaylists(m Me) ResourceHandler { return m.GetPlaylists() }

// UpdateMyPlaylist method
func UpdateMyPlaylist(m Me) ResourceHandler { return m.UpdatePlaylist() }

// AddSongsToMyPlaylist method
func AddSongsToMyPlaylist(m Me) ResourceHandler { return m.AddPlaylistSongs() }

type me struct {
	userRepo     store.UserRepo
	playlistRepo store.PlaylistRepo

	Resource
}

func (m me) Update() ResourceHandler           { return m.update }
func (m me) AddPlaylists() ResourceHandler     { return m.addPlaylists }
func (m me) GetPlaylists() ResourceHandler     { return m.getPlaylists }
func (m me) UpdatePlaylist() ResourceHandler   { return m.updatePlaylist }
func (m me) AddPlaylistSongs() ResourceHandler { return m.addPlaylists }

func (m me) update(w http.ResponseWriter, r *http.Request) {
	auth, ok := validate(w, r, m.Resource, pathCodeMeUpdate)

	if ok == false {
		return
	}

	param, ok := decodeUser(w, r, m.ResponseWriter)

	if ok == false {
		return
	}

	param.Key = store.UserKey(auth.UserKey)
	user, err := store.UpdateUser(m.userRepo, param)

	if err != nil {
		writeRespErr(w, r, m.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, m.ResponseWriter, user)
}

func (m me) addPlaylists(w http.ResponseWriter, r *http.Request) {
	auth, ok := validate(w, r, m.Resource, pathCodeMeAddPlaylists)

	if ok == false {
		return
	}

	param, ok := decodePlaylistKeys(w, r, m.ResponseWriter)

	if ok == false {
		return
	}

	key := store.UserKey(auth.UserKey)
	playlists, err := store.AddUserPlaylists(m.userRepo, key, param.Playlists)

	if err != nil {
		writeRespErr(w, r, m.ResponseWriter, err)
		return
	}

	param.Playlists = playlists
	marshalResponse(w, r, m.ResponseWriter, param)
}

func (m me) getPlaylists(w http.ResponseWriter, r *http.Request) {
	auth, ok := validate(w, r, m.Resource, pathCodeMeGetPlaylists)

	if ok == false {
		return
	}

	key := store.UserKey(auth.UserKey)
	playlists, err := store.GetPlaylistsByUser(m.playlistRepo, key)

	if err != nil {
		writeRespErr(w, r, m.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, m.ResponseWriter, playlists)
}

func (m me) updatePlaylist(w http.ResponseWriter, r *http.Request) {
	auth, ok := validate(w, r, m.Resource, pathCodeMeGetPlaylists)

	if ok == false {
		return
	}

	userKey := store.UserKey(auth.UserKey)
	playlistKey := store.PlaylistKey(mux.Vars(r)["id"])

	err := store.UserHasPlaylist(m.userRepo, userKey, playlistKey)

	if err != nil {
		writeRespErr(w, r, m.ResponseWriter, err)
		return
	}

	param, ok := decodePlaylist(w, r, m.ResponseWriter)

	if ok == false {
		return
	}

	param.Key = playlistKey
	playlist, err := store.UpdatePlaylist(m.playlistRepo, param)

	if err != nil {
		writeRespErr(w, r, m.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, m.ResponseWriter, playlist)
}

func (m me) addPlaylistSongs(w http.ResponseWriter, r *http.Request) {
	auth, ok := validate(w, r, m.Resource, pathCodeMeGetPlaylists)

	if ok == false {
		return
	}

	userKey := store.UserKey(auth.UserKey)
	playlistKey := store.PlaylistKey(mux.Vars(r)["id"])

	err := store.UserHasPlaylist(m.userRepo, userKey, playlistKey)

	if err != nil {
		writeRespErr(w, r, m.ResponseWriter, err)
		return
	}

	param, ok := decodeSongKeys(w, r, m.ResponseWriter)

	if ok == false {
		return
	}

	songs, err := store.AddPlaylistSongs(m.playlistRepo, playlistKey, param.Songs)

	if err != nil {
		writeRespErr(w, r, m.ResponseWriter, err)
		return
	}

	param.Songs = songs
	marshalResponse(w, r, m.ResponseWriter, param)
}

// MeResource method
func MeResource(res Resource, u store.UserRepo, p store.PlaylistRepo) Me {
	return me{
		userRepo:     u,
		playlistRepo: p,

		Resource: res,
	}
}
