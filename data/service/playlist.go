package service

import (
	"net/http"

	"github.com/mownier/duyog/data/store"

	"github.com/gorilla/mux"
)

// Playlist interface
type Playlist interface {
	Create() ResourceHandler
	Update() ResourceHandler
	GetInfo() ResourceHandler
	GetSongs() ResourceHandler
	AddSongs() ResourceHandler
}

// CreatePlaylist method
func CreatePlaylist(p Playlist) ResourceHandler { return p.Create() }

// UpdatePlaylist method
func UpdatePlaylist(p Playlist) ResourceHandler { return p.Update() }

// GetPlaylistInfo method
func GetPlaylistInfo(p Playlist) ResourceHandler { return p.GetInfo() }

// GetPlaylistSongs method
func GetPlaylistSongs(p Playlist) ResourceHandler { return p.GetSongs() }

// AddPlaylistSongs method
func AddPlaylistSongs(p Playlist) ResourceHandler { return p.AddSongs() }

type playlist struct {
	repo store.PlaylistRepo
	Resource
}

func (p playlist) Create() ResourceHandler   { return p.create }
func (p playlist) Update() ResourceHandler   { return p.update }
func (p playlist) GetInfo() ResourceHandler  { return p.getInfo }
func (p playlist) GetSongs() ResourceHandler { return p.getSongs }
func (p playlist) AddSongs() ResourceHandler { return p.addSongs }

func (p playlist) create(w http.ResponseWriter, r *http.Request) {
	auth, ok := validate(w, r, p.Resource, pathCodePlaylistCreate)

	if ok == false {
		return
	}

	param, ok := decodePlaylist(w, r, p.Resource.ResponseWriter)

	if ok == false {
		return
	}

	param.Key = store.PlaylistKey(mux.Vars(r)["id"])
	userKey := store.UserKey(auth.UserKey)
	playlist, err := store.CreatePlaylist(p.repo, userKey, param)

	if err != nil {
		writeRespErr(w, r, p.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, p.ResponseWriter, playlist)
}

func (p playlist) update(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, p.Resource, pathCodePlaylistUpdate)

	if ok == false {
		return
	}

	param, ok := decodePlaylist(w, r, p.ResponseWriter)

	if ok == false {
		return
	}

	param.Key = store.PlaylistKey(mux.Vars(r)["id"])
	playlist, err := store.UpdatePlaylist(p.repo, param)

	if err != nil {
		writeRespErr(w, r, p.ResponseWriter, err)
		return
	}

	playlist.Key = ""

	marshalResponse(w, r, p.ResponseWriter, playlist)
}

func (p playlist) getInfo(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, p.Resource, pathCodePlaylistUpdate)

	if ok == false {
		return
	}

	key := store.PlaylistKey(mux.Vars(r)["id"])
	playlist, err := store.GetPlaylistByKey(p.repo, key)

	if err != nil {
		writeRespErr(w, r, p.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, p.ResponseWriter, playlist)
}

func (p playlist) getSongs(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, p.Resource, pathCodePlaylistGetSongs)

	if ok == false {
		return
	}

	key := store.PlaylistKey(mux.Vars(r)["id"])
	songs, err := store.GetPlaylistSongs(p.repo, key)

	if err != nil {
		writeRespErr(w, r, p.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, p.ResponseWriter, songs)
}

func (p playlist) addSongs(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, p.Resource, pathCodePlaylistGetSongs)

	if ok == false {
		return
	}

	param, ok := decodeSongKeys(w, r, p.ResponseWriter)

	if ok == false {
		return
	}

	key := store.PlaylistKey(mux.Vars(r)["id"])
	songs, err := store.AddPlaylistSongs(p.repo, key, param.Songs)

	if err != nil {
		writeRespErr(w, r, p.ResponseWriter, err)
		return
	}

	param.Songs = songs
	marshalResponse(w, r, p.ResponseWriter, param)
}

// PlaylistResource method
func PlaylistResource(res Resource, r store.PlaylistRepo) Playlist {
	return playlist{
		repo:     r,
		Resource: res,
	}
}
