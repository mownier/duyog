package service

import (
	"github.com/mownier/duyog/data/store"
	"github.com/mownier/duyog/extractor"
	"github.com/mownier/duyog/progerr"
	"github.com/mownier/duyog/validator"
	"github.com/mownier/duyog/writer"
	"encoding/json"
	"net/http"
	"strings"
)

type albumKeyParam struct {
	Albums []store.AlbumKey `json:"albums"`
}

type artistKeyParam struct {
	Artists []store.ArtistKey `json:"artists"`
}

type songKeyParam struct {
	Songs []store.SongKey `json:"songs"`
}

type playlistKeyParam struct {
	Playlists []store.PlaylistKey `json:"playlists"`
}

type songInputParam struct {
	Song    store.Song        `json:"song"`
	Artists []store.ArtistKey `json:"artists"`
	Albums  []store.AlbumKey  `json:"albums"`
}

// Resource struct
type Resource struct {
	AuthExtractor    extractor.Auth
	AuthValidator    validator.Auth
	RequestValidator validator.Request
	ResponseWriter   writer.Response
}

// ResourceHandler function
type ResourceHandler func(http.ResponseWriter, *http.Request)

// RPCHandler function
type RPCHandler func(*http.Request, *validator.DataArgs, *validator.DataReply) error

func validate(w http.ResponseWriter, r *http.Request, res Resource, c pathCode) (validator.AuthReply, bool) {
	var auth validator.AuthReply

	err := validator.ValidateRequest(res.RequestValidator, r)

	if err != nil {
		writeRespErr(w, r, res.ResponseWriter, err)
		return auth, false
	}

	token := extractor.ExtractAuth(res.AuthExtractor, "")
	reply, err := validator.ValidateAuth(res.AuthValidator, token)

	if err != nil {
		writeRespErr(w, r, res.ResponseWriter, err)
		return auth, false
	}

	if isForbidden(c) && strings.ToLower(reply.Role) != "admin" {
		writeRespErr(w, r, res.ResponseWriter, progerr.RequestPathForbidden)
		return auth, false
	}

	auth = reply

	return auth, true
}

func marshalResponse(w http.ResponseWriter, r *http.Request, wr writer.Response, v interface{}) {
	data, err := json.Marshal(v)

	if err != nil {
		writeRespErr(w, r, wr, err)
		return
	}

	writer.WriteResponse(wr, w, r, http.StatusOK, data)
}

func getPlaylistByKey(r store.PlaylistRepo, p store.Playlist) (store.Playlist, error) {
	return store.GetPlaylistByKey(r, p.Key)
}

func decodePlaylist(w http.ResponseWriter, r *http.Request, wr writer.Response) (store.Playlist, bool) {
	var playlist, param store.Playlist

	err := json.NewDecoder(r.Body).Decode(&param)

	if err != nil {
		writeRespErr(w, r, wr, progerr.RequestBodyNotJSON)
		return playlist, false
	}

	playlist = param

	return playlist, true
}

func decodeAlbum(w http.ResponseWriter, r *http.Request, wr writer.Response) (store.Album, bool) {
	var album, param store.Album

	err := json.NewDecoder(r.Body).Decode(&param)

	if err != nil {
		writeRespErr(w, r, wr, progerr.RequestBodyNotJSON)
		return album, false
	}

	album = param

	return album, true
}

func decodeArtist(w http.ResponseWriter, r *http.Request, wr writer.Response) (store.Artist, bool) {
	var artist, param store.Artist

	err := json.NewDecoder(r.Body).Decode(&param)

	if err != nil {
		writeRespErr(w, r, wr, progerr.RequestBodyNotJSON)
		return artist, false
	}

	artist = param

	return artist, true
}

func decodeSong(w http.ResponseWriter, r *http.Request, wr writer.Response) (store.Song, bool) {
	var song, param store.Song

	err := json.NewDecoder(r.Body).Decode(&param)

	if err != nil {
		writeRespErr(w, r, wr, progerr.RequestBodyNotJSON)
		return song, false
	}

	song = param

	return song, true
}

func decodeUser(w http.ResponseWriter, r *http.Request, wr writer.Response) (store.User, bool) {
	var user, param store.User

	err := json.NewDecoder(r.Body).Decode(&param)

	if err != nil {
		writeRespErr(w, r, wr, progerr.RequestBodyNotJSON)
		return user, false
	}

	user = param

	return user, true
}

func decodeSongKeys(w http.ResponseWriter, r *http.Request, wr writer.Response) (songKeyParam, bool) {
	var keys, param songKeyParam

	err := json.NewDecoder(r.Body).Decode(&param)

	if err != nil {
		writeRespErr(w, r, wr, progerr.RequestBodyNotJSON)
		return keys, false
	}

	keys = param

	return keys, true
}

func decodeAlbumKeys(w http.ResponseWriter, r *http.Request, wr writer.Response) (albumKeyParam, bool) {
	var keys, param albumKeyParam

	err := json.NewDecoder(r.Body).Decode(&param)

	if err != nil {
		writeRespErr(w, r, wr, progerr.RequestBodyNotJSON)
		return keys, false
	}

	keys = param

	return keys, true
}

func decodeArtistKeys(w http.ResponseWriter, r *http.Request, wr writer.Response) (artistKeyParam, bool) {
	var keys, param artistKeyParam

	err := json.NewDecoder(r.Body).Decode(&param)

	if err != nil {
		writeRespErr(w, r, wr, progerr.RequestBodyNotJSON)
		return keys, false
	}

	keys = param

	return keys, true
}

func decodePlaylistKeys(w http.ResponseWriter, r *http.Request, wr writer.Response) (playlistKeyParam, bool) {
	var keys, param playlistKeyParam

	err := json.NewDecoder(r.Body).Decode(&param)

	if err != nil {
		writeRespErr(w, r, wr, progerr.RequestBodyNotJSON)
		return keys, false
	}

	keys = param

	return keys, true
}

func decodeSongInput(w http.ResponseWriter, r *http.Request, wr writer.Response) (songInputParam, bool) {
	var input, param songInputParam

	err := json.NewDecoder(r.Body).Decode(&param)

	if err != nil {
		writeRespErr(w, r, wr, progerr.RequestBodyNotJSON)
		return input, false
	}

	input = param

	return input, true
}

func writeRespErr(w http.ResponseWriter, r *http.Request, wr writer.Response, e error) {
	var err progerr.Err

	switch e.(type) {
	case progerr.Err:
		err = e.(progerr.Err)

	default:
		err = progerr.Internal(err)
	}

	writer.WriteResponse(wr, w, r, err.HTTPStatus, err.Data())
}
