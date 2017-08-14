package service

import (
	"github.com/mownier/duyog/progerr"
	"github.com/mownier/duyog/storage/store"

	"net/http"
)

// Song interface
type Song interface {
	UploadAudio() ResourceHandler
	DownloadAudio() ResourceHandler
}

// UploadSongAudio method
func UploadSongAudio(a Song) ResourceHandler {
	return a.UploadAudio()
}

// DownloadSongAudio method
func DownloadSongAudio(a Song) ResourceHandler {
	return a.DownloadAudio()
}

type song struct {
	dir                  string
	repo                 store.SongRepo
	audioSizeLimit       int64
	validAudioExtensions []string

	Resource
}

func (a song) UploadAudio() ResourceHandler   { return a.uploadAudio }
func (a song) DownloadAudio() ResourceHandler { return a.downloadAudio }

func (a song) uploadAudio(w http.ResponseWriter, r *http.Request) {
	if validate(w, r, a.Resource, pathCodeFileUploadMusic, progerr.DataSongNotVerified) == false {
		return
	}

	input := uploadInput{
		dir: a.dir,

		formKey:         a.FormKey,
		sizeLimit:       a.audioSizeLimit,
		multipartMaxMem: a.MaxMem,
		validExtensions: a.validAudioExtensions,

		keyGen: func(id string) (string, error) {
			return store.AddSongAudio(a.repo, id)
		},

		errExtension:    progerr.FileAudioSizeExceeded(a.audioSizeLimit),
		errExceededSize: progerr.FileInvalidAudioExtension(a.validAudioExtensions),

		writer: a.ResponseWriter,
	}

	upload(w, r, input)
}

func (a song) downloadAudio(w http.ResponseWriter, r *http.Request) {
	if validate(w, r, a.Resource, pathCodeFileGetMusic, progerr.DataSongNotVerified) == false {
		return
	}

	download(w, r, a.ResponseWriter, a.dir, func(id, name string) error {
		return store.VerifySong(a.repo, id, name)
	})
}

// SongResource method
func SongResource(res Resource, repo store.SongRepo, dir string, imgSizeLimit int64, validImgExt []string) Song {
	return song{
		dir:                  dir,
		repo:                 repo,
		audioSizeLimit:       imgSizeLimit,
		validAudioExtensions: validImgExt,

		Resource: res,
	}
}
