package main

import (
	"github.com/mownier/duyog/data/service"
	"net/http"

	"github.com/gorilla/mux"
)

func setupAlbumRoutes(r *mux.Router, s service.Album) {
	r.HandleFunc("/"+config.Version+"/album/{id}", service.GetAlbumInfo(s)).Methods(http.MethodGet)
	r.HandleFunc("/"+config.Version+"/album/create", service.CreateAlbum(s)).Methods(http.MethodPost)
	r.HandleFunc("/"+config.Version+"/album/{id}/update", service.UpdateAlbum(s)).Methods(http.MethodPut)
	r.HandleFunc("/"+config.Version+"/album/{id}/songs", service.GetAlbumSongs(s)).Methods(http.MethodGet)
	r.HandleFunc("/"+config.Version+"/album/{id}/artists", service.GetAlbumArtists(s)).Methods(http.MethodGet)
}

func setupArtistRoutes(r *mux.Router, s service.Artist) {
	r.HandleFunc("/"+config.Version+"/artist/{id}", service.GetArtistInfo(s)).Methods(http.MethodGet)
	r.HandleFunc("/"+config.Version+"/artist/create", service.CreateArtist(s)).Methods(http.MethodPost)
	r.HandleFunc("/"+config.Version+"/artist/{id}/update", service.UpdateArtist(s)).Methods(http.MethodPut)
	r.HandleFunc("/"+config.Version+"/artist/{id}/songs", service.GetArtistSongs(s)).Methods(http.MethodGet)
	r.HandleFunc("/"+config.Version+"/artist/{id}/albums", service.GetArtistAlbums(s)).Methods(http.MethodGet)
}

func setupMeRoutes(r *mux.Router, s service.Me) {
	r.HandleFunc("/"+config.Version+"/me/update", service.UpdateMe(s)).Methods(http.MethodPut)
	r.HandleFunc("/"+config.Version+"/me/playlists", service.GetMyPlaylists(s)).Methods(http.MethodGet)
	r.HandleFunc("/"+config.Version+"/me/playlists/add", service.AddToMyPlaylists(s)).Methods(http.MethodPut)
	r.HandleFunc("/"+config.Version+"/me/playlist/{id}/update", service.UpdateMyPlaylist(s)).Methods(http.MethodPut)
	r.HandleFunc("/"+config.Version+"/me/playlist/{id}/songs/add", service.AddSongsToMyPlaylist(s)).Methods(http.MethodPut)
}

func setupPlaylistRoutes(r *mux.Router, s service.Playlist) {
	r.HandleFunc("/"+config.Version+"/playlist/{id}", service.GetPlaylistInfo(s)).Methods(http.MethodGet)
	r.HandleFunc("/"+config.Version+"/playlist/create", service.CreatePlaylist(s)).Methods(http.MethodPost)
	r.HandleFunc("/"+config.Version+"/playlist/{id}/update", service.UpdatePlaylist(s)).Methods(http.MethodPut)
	r.HandleFunc("/"+config.Version+"/playlist/{id}/songs", service.GetPlaylistSongs(s)).Methods(http.MethodGet)
	r.HandleFunc("/"+config.Version+"/playlist/{id}/songs/add", service.AddPlaylistSongs(s)).Methods(http.MethodPut)
}

func setupSongRoutes(r *mux.Router, s service.Song) {
	r.HandleFunc("/"+config.Version+"/song/{id}", service.GetSongInfo(s)).Methods(http.MethodGet)
	r.HandleFunc("/"+config.Version+"/song/create", service.CreateSong(s)).Methods(http.MethodPost)
	r.HandleFunc("/"+config.Version+"/song/{id}/update", service.UpdateSong(s)).Methods(http.MethodPut)
	r.HandleFunc("/"+config.Version+"/song/{id}/albums/update", service.UpdateSongAlbums(s)).Methods(http.MethodPut)
	r.HandleFunc("/"+config.Version+"/song/{id}/artists/update", service.UpdateSongArtists(s)).Methods(http.MethodPut)
}

func setupUserRoutes(r *mux.Router, s service.User) {
	r.HandleFunc("/"+config.Version+"/user/{id}", service.GetUserInfo(s)).Methods(http.MethodGet)
	r.HandleFunc("/"+config.Version+"/user/register", service.RegisterUser(s)).Methods(http.MethodPost)
	r.HandleFunc("/"+config.Version+"/user/{id}/playlists", service.GetUserPlaylists(s)).Methods(http.MethodGet)
}
