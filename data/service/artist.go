package service

import (
	"github.com/mownier/duyog/data/store"
	"net/http"

	"github.com/gorilla/mux"
)

// Artist interface
type Artist interface {
	Create() ResourceHandler
	Update() ResourceHandler
	GetInfo() ResourceHandler
	GetSongs() ResourceHandler
	GetAlbums() ResourceHandler
}

// CreateArtist method
func CreateArtist(a Artist) ResourceHandler { return a.Create() }

// UpdateArtist method
func UpdateArtist(a Artist) ResourceHandler { return a.Update() }

// GetArtistInfo method
func GetArtistInfo(a Artist) ResourceHandler { return a.GetInfo() }

// GetArtistSongs method
func GetArtistSongs(a Artist) ResourceHandler { return a.GetSongs() }

// GetArtistAlbums method
func GetArtistAlbums(a Artist) ResourceHandler { return a.GetAlbums() }

type artist struct {
	repo store.ArtistRepo

	Resource
}

func (a artist) Create() ResourceHandler    { return a.create }
func (a artist) Update() ResourceHandler    { return a.update }
func (a artist) GetInfo() ResourceHandler   { return a.getInfo }
func (a artist) GetSongs() ResourceHandler  { return a.getSongs }
func (a artist) GetAlbums() ResourceHandler { return a.getAlbums }

func (a artist) create(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, a.Resource, pathCodeArtistCreate)

	if ok == false {
		return
	}

	param, ok := decodeArtist(w, r, a.ResponseWriter)

	if ok == false {
		return
	}

	param.Key = store.ArtistKey(mux.Vars(r)["id"])
	artist, err := store.CreateArtist(a.repo, param)

	if err != nil {
		writeRespErr(w, r, a.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, a.ResponseWriter, artist)
}

func (a artist) update(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, a.Resource, pathCodeArtistUpdate)

	if ok == false {
		return
	}

	param, ok := decodeArtist(w, r, a.ResponseWriter)

	if ok == false {
		return
	}

	param.Key = store.ArtistKey(mux.Vars(r)["id"])
	artist, err := store.UpdateArtist(a.repo, param)

	if err != nil {
		writeRespErr(w, r, a.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, a.ResponseWriter, artist)
}

func (a artist) getInfo(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, a.Resource, pathCodeArtistUpdate)

	if ok == false {
		return
	}

	key := store.ArtistKey(mux.Vars(r)["id"])
	artist, err := store.GetArtistByKey(a.repo, key)

	if err != nil {
		writeRespErr(w, r, a.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, a.ResponseWriter, artist)
}

func (a artist) getSongs(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, a.Resource, pathCodeArtistGetSongs)

	if ok == false {
		return
	}

	key := store.ArtistKey(mux.Vars(r)["id"])
	songs, err := store.GetArtistSongs(a.repo, key)

	if err != nil {
		writeRespErr(w, r, a.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, a.ResponseWriter, songs)
}

func (a artist) getAlbums(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, a.Resource, pathCodeArtistGetAlbums)

	if ok == false {
		return
	}

	key := store.ArtistKey(mux.Vars(r)["id"])
	albums, err := store.GetArtistAlbums(a.repo, key)

	if err != nil {
		writeRespErr(w, r, a.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, a.ResponseWriter, albums)
}

// ArtistResource method
func ArtistResource(res Resource, r store.ArtistRepo) Artist {
	return artist{
		repo:     r,
		Resource: res,
	}
}
