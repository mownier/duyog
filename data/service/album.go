package service

import (
	"net/http"

	"github.com/mownier/duyog/data/store"

	"github.com/gorilla/mux"
)

// Album interface
type Album interface {
	Create() ResourceHandler
	Update() ResourceHandler
	GetInfo() ResourceHandler
	GetSongs() ResourceHandler
	GetArtists() ResourceHandler
}

// CreateAlbum method
func CreateAlbum(a Album) ResourceHandler { return a.Create() }

// UpdateAlbum method
func UpdateAlbum(a Album) ResourceHandler { return a.Update() }

// GetAlbumInfo method
func GetAlbumInfo(a Album) ResourceHandler { return a.GetInfo() }

// GetAlbumSongs method
func GetAlbumSongs(a Album) ResourceHandler { return a.GetSongs() }

// GetAlbumArtists method
func GetAlbumArtists(a Album) ResourceHandler { return a.GetArtists() }

type album struct {
	repo store.AlbumRepo

	Resource
}

func (a album) Create() ResourceHandler     { return a.create }
func (a album) Update() ResourceHandler     { return a.update }
func (a album) GetInfo() ResourceHandler    { return a.getInfo }
func (a album) GetSongs() ResourceHandler   { return a.getSongs }
func (a album) GetArtists() ResourceHandler { return a.getArtists }

func (a album) create(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, a.Resource, pathCodeAlbumCreate)

	if ok == false {
		return
	}

	param, ok := decodeAlbum(w, r, a.ResponseWriter)

	if ok == false {
		return
	}

	param.Key = store.AlbumKey(mux.Vars(r)["id"])
	album, err := store.CreateAlbum(a.repo, param)

	if err != nil {
		writeRespErr(w, r, a.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, a.ResponseWriter, album)
}

func (a album) update(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, a.Resource, pathCodeAlbumUpdate)

	if ok == false {
		return
	}

	param, ok := decodeAlbum(w, r, a.ResponseWriter)

	if ok == false {
		return
	}

	param.Key = store.AlbumKey(mux.Vars(r)["id"])
	album, err := store.UpdateAlbum(a.repo, param)

	if err != nil {
		writeRespErr(w, r, a.ResponseWriter, err)
		return
	}

	album.Key = ""

	marshalResponse(w, r, a.ResponseWriter, album)
}

func (a album) getInfo(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, a.Resource, pathCodeAlbumGetInfo)

	if ok == false {
		return
	}

	key := store.AlbumKey(mux.Vars(r)["id"])
	album, err := store.GetAlbumByKey(a.repo, key)

	if err != nil {
		writeRespErr(w, r, a.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, a.ResponseWriter, album)
}

func (a album) getSongs(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, a.Resource, pathCodeAlbumGetSongs)

	if ok == false {
		return
	}

	key := store.AlbumKey(mux.Vars(r)["id"])
	songs, err := store.GetAlbumSongs(a.repo, key)

	if err != nil {
		writeRespErr(w, r, a.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, a.ResponseWriter, songs)
}

func (a album) getArtists(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, a.Resource, pathCodeAlbumGetArtists)

	if ok == false {
		return
	}

	key := store.AlbumKey(mux.Vars(r)["id"])
	artists, err := store.GetAlbumArtists(a.repo, key)

	if err != nil {
		writeRespErr(w, r, a.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, a.ResponseWriter, artists)
}

// AlbumResource method
func AlbumResource(res Resource, repo store.AlbumRepo) Album {
	return album{
		repo:     repo,
		Resource: res,
	}
}
