package service

import (
	"duyog/progerr"
	"duyog/storage/store"

	"net/http"
)

// Artist interface
type Artist interface {
	UploadPhoto() ResourceHandler
	DownloadPhoto() ResourceHandler
}

// UploadArtistPhoto method
func UploadArtistPhoto(a Artist) ResourceHandler {
	return a.UploadPhoto()
}

// DownloadArtistPhoto method
func DownloadArtistPhoto(a Artist) ResourceHandler {
	return a.DownloadPhoto()
}

type artist struct {
	dir                  string
	repo                 store.ArtistRepo
	imageSizeLimit       int64
	validImageExtensions []string

	Resource
}

func (a artist) UploadPhoto() ResourceHandler   { return a.uploadPhoto }
func (a artist) DownloadPhoto() ResourceHandler { return a.downloadPhoto }

func (a artist) uploadPhoto(w http.ResponseWriter, r *http.Request) {
	if validate(w, r, a.Resource, pathCodeFileUploadArtistPhoto, progerr.DataArtistNotVerified) == false {
		return
	}

	input := uploadInput{
		dir: a.dir,

		formKey:         a.FormKey,
		sizeLimit:       a.imageSizeLimit,
		multipartMaxMem: a.MaxMem,
		validExtensions: a.validImageExtensions,

		keyGen: func(id string) (string, error) {
			return store.AddArtistPhoto(a.repo, id)
		},

		errExtension:    progerr.FileImageSizeExceeded(a.imageSizeLimit),
		errExceededSize: progerr.FileInvalidImageExtension(a.validImageExtensions),

		writer: a.ResponseWriter,
	}

	upload(w, r, input)
}

func (a artist) downloadPhoto(w http.ResponseWriter, r *http.Request) {
	if validate(w, r, a.Resource, pathCodeFileGetArtistPhoto, progerr.DataArtistNotVerified) == false {
		return
	}

	download(w, r, a.ResponseWriter, a.dir, func(id, name string) error {
		return store.VerifyArtist(a.repo, id, name)
	})
}

// ArtistResource method
func ArtistResource(res Resource, repo store.ArtistRepo, dir string, imgSizeLimit int64, validImgExt []string) Artist {
	return artist{
		dir:                  dir,
		repo:                 repo,
		imageSizeLimit:       imgSizeLimit,
		validImageExtensions: validImgExt,

		Resource: res,
	}
}
