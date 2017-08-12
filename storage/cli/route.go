package main

import (
	"duyog/storage/service"
	"net/http"

	"github.com/gorilla/mux"
)

func setupAlbumRoutes(r *mux.Router, s service.Album) {
	r.HandleFunc("/"+config.Version+"/file/album/{id}/photo/upload", service.UploadAlbumPhoto(s)).Methods(http.MethodPost).Queries("access_token", "{token}")
	r.HandleFunc("/"+config.Version+"/file/album/{id}/photo/{name}", service.DownloadAlbumPhoto(s)).Methods(http.MethodGet).Queries("access_token", "{token}")
}

func setupArtistRoutes(r *mux.Router, s service.Artist) {
	r.HandleFunc("/"+config.Version+"/file/artist/{id}/photo/upload", service.UploadArtistPhoto(s)).Methods(http.MethodPost).Queries("access_token", "{token}")
	r.HandleFunc("/"+config.Version+"/file/artist/{id}/photo/{name}", service.DownloadArtistPhoto(s)).Methods(http.MethodGet).Queries("access_token", "{token}")
}

func setupPlaylistRoutes(r *mux.Router, s service.Playlist) {
	r.HandleFunc("/"+config.Version+"/file/playlist/{id}/photo/upload", service.UploadPlaylistPhoto(s)).Methods(http.MethodPost).Queries("access_token", "{token}")
	r.HandleFunc("/"+config.Version+"/file/playlist/{id}/photo/{name}", service.DownloadPlaylistPhoto(s)).Methods(http.MethodGet).Queries("access_token", "{token}")
}

func setupSongRoutes(r *mux.Router, s service.Song) {
	r.HandleFunc("/"+config.Version+"/file/song/{id}/audio/upload", service.UploadSongAudio(s)).Methods(http.MethodPost).Queries("access_token", "{token}")
	r.HandleFunc("/"+config.Version+"/file/song/{id}/audio/{name}", service.DownloadSongAudio(s)).Methods(http.MethodGet).Queries("access_token", "{token}")
}

func setupUserRoutes(r *mux.Router, s service.User) {
	r.HandleFunc("/"+config.Version+"/file/user/{id}/avatar/upload", service.UploadUserAvatar(s)).Methods(http.MethodPost).Queries("access_token", "{token}")
	r.HandleFunc("/"+config.Version+"/file/user/{id}/avatar/{name}", service.DownloadUserAvatar(s)).Methods(http.MethodGet).Queries("access_token", "{token}")
}
