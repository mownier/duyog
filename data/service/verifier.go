package service

import (
	"net/http"

	"github.com/mownier/duyog/data/store"
	"github.com/mownier/duyog/logger"
	"github.com/mownier/duyog/progerr"
	"github.com/mownier/duyog/validator"

	"github.com/gorilla/rpc/v2/json2"
)

// Verifier interface
type Verifier interface {
	ValidateSong(*http.Request, *validator.DataArgs, *validator.DataReply) error
	ValidateUser(*http.Request, *validator.DataArgs, *validator.DataReply) error
	ValidateAlbum(*http.Request, *validator.DataArgs, *validator.DataReply) error
	ValidateArtist(*http.Request, *validator.DataArgs, *validator.DataReply) error
	ValidatePlaylist(*http.Request, *validator.DataArgs, *validator.DataReply) error
}

type verifier struct {
	userRepo     store.UserRepo
	songRepo     store.SongRepo
	albumRepo    store.AlbumRepo
	artistRepo   store.ArtistRepo
	playlistRepo store.PlaylistRepo

	log logger.Request
}

func (v verifier) ValidateAlbum(r *http.Request, args *validator.DataArgs, reply *validator.DataReply) error {
	return v.validate(r, args, reply, progerr.AlbumNotFound, func() error {
		_, err := store.GetAlbumByKey(v.albumRepo, store.AlbumKey(args.Key))
		return err
	})
}

func (v verifier) ValidateArtist(r *http.Request, args *validator.DataArgs, reply *validator.DataReply) error {
	return v.validate(r, args, reply, progerr.ArtistNotFound, func() error {
		_, err := store.GetArtistByKey(v.artistRepo, store.ArtistKey(args.Key))
		return err
	})
}

func (v verifier) ValidatePlaylist(r *http.Request, args *validator.DataArgs, reply *validator.DataReply) error {
	return v.validate(r, args, reply, progerr.PlaylistNotFound, func() error {
		_, err := store.GetPlaylistByKey(v.playlistRepo, store.PlaylistKey(args.Key))
		return err
	})
}

func (v verifier) ValidateSong(r *http.Request, args *validator.DataArgs, reply *validator.DataReply) error {
	return v.validate(r, args, reply, progerr.SongNotFound, func() error {
		_, err := store.GetSongByKey(v.songRepo, store.SongKey(args.Key))
		return err
	})
}

func (v verifier) ValidateUser(r *http.Request, args *validator.DataArgs, reply *validator.DataReply) error {
	return v.validate(r, args, reply, progerr.UserNotFound, func() error {
		_, err := store.GetUserByKey(v.userRepo, store.UserKey(args.Key))
		return err
	})
}

func (v verifier) validate(r *http.Request, args *validator.DataArgs, reply *validator.DataReply, e progerr.Err, storeFunc func() error) error {
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

// NewVerifier method
func NewVerifier(u store.UserRepo, s store.SongRepo, al store.AlbumRepo, ar store.ArtistRepo, p store.PlaylistRepo, l logger.Request) Verifier {
	return verifier{
		userRepo:     u,
		songRepo:     s,
		albumRepo:    al,
		artistRepo:   ar,
		playlistRepo: p,

		log: l,
	}
}
