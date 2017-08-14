package service

import (
	"github.com/mownier/duyog/progerr"
	"github.com/mownier/duyog/storage/store"

	"net/http"
)

// User interface
type User interface {
	UploadAvatar() ResourceHandler
	DownloadAvatar() ResourceHandler
}

// UploadUserAvatar method
func UploadUserAvatar(a User) ResourceHandler {
	return a.UploadAvatar()
}

// DownloadUserAvatar method
func DownloadUserAvatar(a User) ResourceHandler {
	return a.DownloadAvatar()
}

type user struct {
	dir                  string
	repo                 store.UserRepo
	imageSizeLimit       int64
	validImageExtensions []string

	Resource
}

func (a user) UploadAvatar() ResourceHandler   { return a.uploadAvatar }
func (a user) DownloadAvatar() ResourceHandler { return a.downloadAvatar }

func (a user) uploadAvatar(w http.ResponseWriter, r *http.Request) {
	if validate(w, r, a.Resource, pathCodeFileUploadUserAvatar, progerr.DataUserNotVerified) == false {
		return
	}

	input := uploadInput{
		dir: a.dir,

		formKey:         a.FormKey,
		sizeLimit:       a.imageSizeLimit,
		multipartMaxMem: a.MaxMem,
		validExtensions: a.validImageExtensions,

		keyGen: func(id string) (string, error) {
			return store.AddUserAvatar(a.repo, id)
		},

		errExtension:    progerr.FileImageSizeExceeded(a.imageSizeLimit),
		errExceededSize: progerr.FileInvalidImageExtension(a.validImageExtensions),

		writer: a.ResponseWriter,
	}

	upload(w, r, input)
}

func (a user) downloadAvatar(w http.ResponseWriter, r *http.Request) {
	if validate(w, r, a.Resource, pathCodeFileGetUserAvatar, progerr.DataUserNotVerified) == false {
		return
	}

	download(w, r, a.ResponseWriter, a.dir, func(id, name string) error {
		return store.VerifyUser(a.repo, id, name)
	})
}

// UserResource method
func UserResource(res Resource, repo store.UserRepo, dir string, imgSizeLimit int64, validImgExt []string) User {
	return user{
		dir:                  dir,
		repo:                 repo,
		imageSizeLimit:       imgSizeLimit,
		validImageExtensions: validImgExt,

		Resource: res,
	}
}
