package service

import (
	"github.com/mownier/duyog/progerr"
	"github.com/mownier/duyog/storage/store"

	"net/http"
)

// Playlist interface
type Playlist interface {
	UploadPhoto() ResourceHandler
	DownloadPhoto() ResourceHandler
}

// UploadPlaylistPhoto method
func UploadPlaylistPhoto(a Playlist) ResourceHandler {
	return a.UploadPhoto()
}

// DownloadPlaylistPhoto method
func DownloadPlaylistPhoto(a Playlist) ResourceHandler {
	return a.DownloadPhoto()
}

type playlist struct {
	dir                  string
	repo                 store.PlaylistRepo
	imageSizeLimit       int64
	validImageExtensions []string

	Resource
}

func (a playlist) UploadPhoto() ResourceHandler   { return a.uploadPhoto }
func (a playlist) DownloadPhoto() ResourceHandler { return a.downloadPhoto }

func (a playlist) uploadPhoto(w http.ResponseWriter, r *http.Request) {
	if validate(w, r, a.Resource, pathCodeFileUploadPlaylistPhoto, progerr.DataPlaylistNotVerified) == false {
		return
	}

	input := uploadInput{
		dir: a.dir,

		formKey:         a.FormKey,
		sizeLimit:       a.imageSizeLimit,
		multipartMaxMem: a.MaxMem,
		validExtensions: a.validImageExtensions,

		keyGen: func(id string) (string, error) {
			return store.AddPlaylistPhoto(a.repo, id)
		},

		errExtension:    progerr.FileImageSizeExceeded(a.imageSizeLimit),
		errExceededSize: progerr.FileInvalidImageExtension(a.validImageExtensions),

		writer: a.ResponseWriter,
	}

	upload(w, r, input)
}

func (a playlist) downloadPhoto(w http.ResponseWriter, r *http.Request) {
	if validate(w, r, a.Resource, pathCodeFileGetPlaylistPhoto, progerr.DataPlaylistNotVerified) == false {
		return
	}

	download(w, r, a.ResponseWriter, a.dir, func(id, name string) error {
		return store.VerifyPlaylist(a.repo, id, name)
	})
}

// PlaylistResource method
func PlaylistResource(res Resource, repo store.PlaylistRepo, dir string, imgSizeLimit int64, validImgExt []string) Playlist {
	return playlist{
		dir:                  dir,
		repo:                 repo,
		imageSizeLimit:       imgSizeLimit,
		validImageExtensions: validImgExt,

		Resource: res,
	}
}
