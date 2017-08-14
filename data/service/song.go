package service

import (
	"github.com/mownier/duyog/data/store"
	"net/http"

	"github.com/gorilla/mux"
)

// Song interface
type Song interface {
	Create() ResourceHandler
	Update() ResourceHandler
	GetInfo() ResourceHandler
	UpdateAlbums() ResourceHandler
	UpdateArtists() ResourceHandler
}

// CreateSong method
func CreateSong(s Song) ResourceHandler { return s.Create() }

// UpdateSong method
func UpdateSong(s Song) ResourceHandler { return s.Update() }

// GetSongInfo method
func GetSongInfo(s Song) ResourceHandler { return s.GetInfo() }

// UpdateSongAlbums method
func UpdateSongAlbums(s Song) ResourceHandler { return s.UpdateAlbums() }

// UpdateSongArtists method
func UpdateSongArtists(s Song) ResourceHandler { return s.UpdateArtists() }

type song struct {
	repo store.SongRepo

	Resource
}

func (s song) Create() ResourceHandler        { return s.create }
func (s song) Update() ResourceHandler        { return s.update }
func (s song) GetInfo() ResourceHandler       { return s.getInfo }
func (s song) UpdateAlbums() ResourceHandler  { return s.updateAlbums }
func (s song) UpdateArtists() ResourceHandler { return s.updateArtists }

func (s song) create(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, s.Resource, pathCodeSongCreate)

	if ok == false {
		return
	}

	param, ok := decodeSongInput(w, r, s.ResponseWriter)

	if ok == false {
		return
	}

	songs, err := store.CreateSong(s.repo, param.Song, param.Artists, param.Albums)

	if err != nil {
		writeRespErr(w, r, s.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, s.ResponseWriter, songs)
}

func (s song) update(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, s.Resource, pathCodeSongUpdate)

	if ok == false {
		return
	}

	param, ok := decodeSong(w, r, s.ResponseWriter)

	if ok == false {
		return
	}

	param.Key = store.SongKey(mux.Vars(r)["id"])
	song, err := store.UpdateSong(s.repo, param)

	if err != nil {
		writeRespErr(w, r, s.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, s.ResponseWriter, song)
}

func (s song) getInfo(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, s.Resource, pathCodeSongGetInfo)

	if ok == false {
		return
	}

	key := store.SongKey(mux.Vars(r)["id"])
	song, err := store.GetSongByKey(s.repo, key)

	if err != nil {
		writeRespErr(w, r, s.ResponseWriter, err)
		return
	}

	marshalResponse(w, r, s.ResponseWriter, song)
}

func (s song) updateAlbums(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, s.Resource, pathCodeSongUpdateAlbums)

	if ok == false {
		return
	}

	param, ok := decodeAlbumKeys(w, r, s.ResponseWriter)

	if ok == false {
		return
	}

	key := store.SongKey(mux.Vars(r)["id"])
	albums, err := store.UpdateSongAlbums(s.repo, key, param.Albums)

	if err != nil {
		writeRespErr(w, r, s.ResponseWriter, err)
		return
	}

	param.Albums = albums
	marshalResponse(w, r, s.ResponseWriter, param)
}

func (s song) updateArtists(w http.ResponseWriter, r *http.Request) {
	_, ok := validate(w, r, s.Resource, pathCodeSongUpdateArtists)

	if ok == false {
		return
	}

	param, ok := decodeArtistKeys(w, r, s.ResponseWriter)

	if ok == false {
		return
	}

	key := store.SongKey(mux.Vars(r)["id"])
	artists, err := store.UpdateSongArtists(s.repo, key, param.Artists)

	if err != nil {
		writeRespErr(w, r, s.ResponseWriter, err)
		return
	}

	param.Artists = artists
	marshalResponse(w, r, s.ResponseWriter, param)
}

// SongResource method
func SongResource(res Resource, r store.SongRepo) Song {
	return song{
		repo:     r,
		Resource: res,
	}
}
