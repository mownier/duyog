package service

import (
	"github.com/mownier/duyog/data/store"
	"github.com/mownier/duyog/logger"
	"github.com/mownier/duyog/progerr"
	"github.com/mownier/duyog/validator"
	"net/http"

	"github.com/gorilla/rpc/v2/json2"
)

// Validator interface
type Validator interface {
	ValidateSong() RPCHandler
	ValidateUser() RPCHandler
	ValidateAlbum() RPCHandler
	ValidateArtist() RPCHandler
	ValidatePlaylist() RPCHandler
}

// ValidateAlbum method
func ValidateAlbum(v Validator) RPCHandler { return v.ValidateAlbum() }

// ValidateArtist method
func ValidateArtist(v Validator) RPCHandler { return v.ValidateArtist() }

// ValidatePlaylist method
func ValidatePlaylist(v Validator) RPCHandler { return v.ValidatePlaylist() }

// ValidateSong method
func ValidateSong(v Validator) RPCHandler { return v.ValidateSong() }

// ValidateUser method
func ValidateUser(v Validator) RPCHandler { return v.ValidateUser() }

type dataValidator struct {
	userRepo     store.UserRepo
	songRepo     store.SongRepo
	albumRepo    store.AlbumRepo
	artistRepo   store.ArtistRepo
	playlistRepo store.PlaylistRepo

	log logger.Request
}

func (v dataValidator) ValidateSong() RPCHandler     { return v.song }
func (v dataValidator) ValidateUser() RPCHandler     { return v.user }
func (v dataValidator) ValidateAlbum() RPCHandler    { return v.album }
func (v dataValidator) ValidateArtist() RPCHandler   { return v.artist }
func (v dataValidator) ValidatePlaylist() RPCHandler { return v.playlist }

func (v dataValidator) album(r *http.Request, args *validator.DataArgs, reply *validator.DataReply) error {
	return v.validate(r, args, reply, progerr.AlbumNotFound,
		func() error {
			_, err := store.GetAlbumByKey(v.albumRepo, store.AlbumKey(args.Key))
			return err
		})
}

func (v dataValidator) artist(r *http.Request, args *validator.DataArgs, reply *validator.DataReply) error {
	return v.validate(r, args, reply, progerr.ArtistNotFound,
		func() error {
			_, err := store.GetArtistByKey(v.artistRepo, store.ArtistKey(args.Key))
			return err
		})
}

func (v dataValidator) playlist(r *http.Request, args *validator.DataArgs, reply *validator.DataReply) error {
	return v.validate(r, args, reply, progerr.PlaylistNotFound,
		func() error {
			_, err := store.GetPlaylistByKey(v.playlistRepo, store.PlaylistKey(args.Key))
			return err
		})
}

func (v dataValidator) song(r *http.Request, args *validator.DataArgs, reply *validator.DataReply) error {
	return v.validate(r, args, reply, progerr.SongNotFound,
		func() error {
			_, err := store.GetSongByKey(v.songRepo, store.SongKey(args.Key))
			return err
		})
}

func (v dataValidator) user(r *http.Request, args *validator.DataArgs, reply *validator.DataReply) error {
	return v.validate(r, args, reply, progerr.UserNotFound,
		func() error {
			_, err := store.GetUserByKey(v.userRepo, store.UserKey(args.Key))
			return err
		})
}

func (v dataValidator) validate(r *http.Request, args *validator.DataArgs, reply *validator.DataReply, e progerr.Err, storeFunc func() error) error {
	logger.LogRequest(v.log, r)

	if r.Method != http.MethodPost {
		return &json2.Error{
			Code:    json2.E_NO_METHOD,
			Message: progerr.MethodNotAllowed.Message,
		}
	}

	if r.Header.Get("Content-Type") != "application/json" {
		return &json2.Error{
			Code:    json2.E_INVALID_REQ,
			Message: progerr.ContentTypeNotJSON.Message,
		}
	}

	if args.Key == "" {
		return &json2.Error{
			Code:    json2.E_INVALID_REQ,
			Message: e.Message,
		}
	}

	err := storeFunc()

	if err != nil {
		return &json2.Error{
			Code:    json2.E_INVALID_REQ,
			Message: e.Message,
		}
	}

	reply.OK = true

	return nil
}

// DataValidator method
func DataValidator(u store.UserRepo, s store.SongRepo, al store.AlbumRepo, ar store.ArtistRepo, p store.PlaylistRepo, l logger.Request) Validator {
	return dataValidator{
		userRepo:     u,
		songRepo:     s,
		albumRepo:    al,
		artistRepo:   ar,
		playlistRepo: p,

		log: l,
	}
}
