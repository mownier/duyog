package service

import (
	"github.com/mownier/duyog/progerr"
	"github.com/mownier/duyog/storage/store"

	"net/http"
)

// Album interface
type Album interface {
	UploadPhoto() ResourceHandler
	DownloadPhoto() ResourceHandler
}

// UploadAlbumPhoto method
func UploadAlbumPhoto(a Album) ResourceHandler {
	return a.UploadPhoto()
}

// DownloadAlbumPhoto method
func DownloadAlbumPhoto(a Album) ResourceHandler {
	return a.DownloadPhoto()
}

type album struct {
	dir                  string
	repo                 store.AlbumRepo
	imageSizeLimit       int64
	validImageExtensions []string

	Resource
}

func (a album) UploadPhoto() ResourceHandler   { return a.uploadPhoto }
func (a album) DownloadPhoto() ResourceHandler { return a.downloadPhoto }

func (a album) uploadPhoto(w http.ResponseWriter, r *http.Request) {
	if validate(w, r, a.Resource, pathCodeFileUploadAlbumPhoto, progerr.DataAlbumNotVerified) == false {
		return
	}

	input := uploadInput{
		dir: a.dir,

		formKey:         a.FormKey,
		sizeLimit:       a.imageSizeLimit,
		multipartMaxMem: a.MaxMem,
		validExtensions: a.validImageExtensions,

		keyGen: func(id string) (string, error) {
			return store.AddAlbumPhoto(a.repo, id)
		},

		errExtension:    progerr.FileImageSizeExceeded(a.imageSizeLimit),
		errExceededSize: progerr.FileInvalidImageExtension(a.validImageExtensions),

		writer: a.ResponseWriter,
	}

	upload(w, r, input)
}

func (a album) downloadPhoto(w http.ResponseWriter, r *http.Request) {
	if validate(w, r, a.Resource, pathCodeFileGetAlbumPhoto, progerr.DataAlbumNotVerified) == false {
		return
	}

	download(w, r, a.ResponseWriter, a.dir, func(id, name string) error {
		return store.VerifyAlbum(a.repo, id, name)
	})
}

// AlbumResource method
func AlbumResource(res Resource, repo store.AlbumRepo, dir string, imgSizeLimit int64, validImgExt []string) Album {
	return album{
		dir:                  dir,
		repo:                 repo,
		imageSizeLimit:       imgSizeLimit,
		validImageExtensions: validImgExt,

		Resource: res,
	}
}
