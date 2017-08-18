package progerr

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var (
	// ContentTypeNotJSON error
	ContentTypeNotJSON = New("content type not 'application/json'", http.StatusBadRequest)
	// InvalidAuthHeaderValue error
	InvalidAuthHeaderValue = New("invalid authorization header", http.StatusUnauthorized)
	// InvalidRequestParameter error
	InvalidRequestParameter = New("invalid request parameters", http.StatusBadRequest)
	// MethodNotAllowed error
	MethodNotAllowed = New("method not allowed", http.StatusMethodNotAllowed)
	// NoAuthHeader error
	NoAuthHeader = New("no authorization header", http.StatusUnauthorized)
	// RequestBodyNotJSON error
	RequestBodyNotJSON = New("can not parse request body to JSON", http.StatusBadRequest)
	// RequestPathForbidden error
	RequestPathForbidden = New("forbidden", http.StatusForbidden)
	// Unauthorized error
	Unauthorized = New("unauthorized", http.StatusUnauthorized)

	// FileNotFound error
	FileNotFound = New("file not found", http.StatusNotFound)

	// DataAlbumNotVerified error
	DataAlbumNotVerified = New("album not verified", http.StatusBadRequest)
	// DataArtistNotVerified error
	DataArtistNotVerified = New("artist not verified", http.StatusBadRequest)
	// DataPlaylistNotVerified error
	DataPlaylistNotVerified = New("playlist not verified", http.StatusBadRequest)
	// DataSongNotVerified error
	DataSongNotVerified = New("song not verified", http.StatusBadRequest)
	// DataUserNotVerified error
	DataUserNotVerified = New("user not verified", http.StatusBadRequest)

	// PlaylistHasNoCreator error
	PlaylistHasNoCreator = New("playlist has no creator", http.StatusBadRequest)
	// PlaylistHasNoSongs error
	PlaylistHasNoSongs = New("playlist has no songs", http.StatusBadRequest)
	// PlaylistInvalidKey error
	PlaylistInvalidKey = New("invalid playlist key", http.StatusBadRequest)
	// PlaylistInvalidName error
	PlaylistInvalidName = New("invalid playlist name", http.StatusBadRequest)
	// PlaylistNoSongsAdded error
	PlaylistNoSongsAdded = New("no songs added in the playlist", http.StatusBadRequest)
	// PlaylistNotAdded error
	PlaylistNotAdded = New("there are no playlists that are added", http.StatusBadRequest)
	// PlaylistNotFound error
	PlaylistNotFound = New("playlist not found", http.StatusNotFound)
	// PlaylistNothingToUpdate error
	PlaylistNothingToUpdate = New("there is nothing to update on the playlist", http.StatusBadRequest)

	// UserAlreadyExists error
	UserAlreadyExists = New("user already exists", http.StatusBadRequest)
	// UserChangePassFailed error
	UserChangePassFailed = New("unable to change password", http.StatusBadRequest)
	// UserEmptyPass error
	UserEmptyPass = New("password is empty", http.StatusBadRequest)
	// UserHasNoPlaylists error
	UserHasNoPlaylists = New("user has no playlists", http.StatusBadRequest)
	// UserInvalidCredentials error
	UserInvalidCredentials = New("invalid user credentials", http.StatusBadRequest)
	// UserInvalidEmail error
	UserInvalidEmail = New("invalid user email", http.StatusBadRequest)
	// UserInvalidInfo error
	UserInvalidInfo = New("invalid user info", http.StatusBadRequest)
	// UserInvalidKey error
	UserInvalidKey = New("invalid user key", http.StatusBadRequest)
	// UserInvalidPassword error
	UserInvalidPassword = New("invalid password", http.StatusBadRequest)
	// UserMismatchedCurrentPass error
	UserMismatchedCurrentPass = New("current password not the same", http.StatusBadRequest)
	// UserMismatchedPass error
	UserMismatchedPass = New("mismatched password", http.StatusBadRequest)
	// UserNotFound error
	UserNotFound = New("user not found", http.StatusNotFound)
	// UserNothingToUpdate error
	UserNothingToUpdate = New("nothing to update on user", http.StatusBadRequest)
	// UserNotRegistered error
	UserNotRegistered = New("user not registered", http.StatusBadRequest)

	// SongAlbumsNotUpdated error
	SongAlbumsNotUpdated = New("albums not updated", http.StatusBadRequest)
	// SongArtistsNotUpdated error
	SongArtistsNotUpdated = New("artists not updated", http.StatusBadRequest)
	// SongHasNoArtist error
	SongHasNoArtist = New("song should have at least 1 artist", http.StatusBadRequest)
	// SongInvalidInfo error
	SongInvalidInfo = New("song has invalid info", http.StatusBadRequest)
	// SongInvalidKey error
	SongInvalidKey = New("invalid song key", http.StatusBadRequest)
	// SongNothingToUpdate error
	SongNothingToUpdate = New("nothing to update on song", http.StatusBadRequest)
	// SongNotFound error
	SongNotFound = New("song not found", http.StatusBadRequest)

	// ArtistHasNoAlbums error
	ArtistHasNoAlbums = New("artist has no albums", http.StatusBadRequest)
	// ArtistHasNoSongs error
	ArtistHasNoSongs = New("artist has no songs", http.StatusBadRequest)
	// ArtistInvalidInfo error
	ArtistInvalidInfo = New("artist has invalid info", http.StatusBadRequest)
	// ArtistInvalidKey error
	ArtistInvalidKey = New("invalid artist key", http.StatusBadRequest)
	// ArtistNotFound error
	ArtistNotFound = New("artist not found", http.StatusNotFound)
	// ArtistNothingToUpdate error
	ArtistNothingToUpdate = New("nothing to update on artist", http.StatusBadRequest)

	// AlbumHasNoArtists error
	AlbumHasNoArtists = New("Album has no artists", http.StatusBadRequest)
	// AlbumHasNoSongs error
	AlbumHasNoSongs = New("Album has no songs", http.StatusBadRequest)
	// AlbumInvalidInfo error
	AlbumInvalidInfo = New("Album has invalid info", http.StatusBadRequest)
	// AlbumInvalidKey error
	AlbumInvalidKey = New("Invalid album key", http.StatusBadRequest)
	// AlbumNotFound error
	AlbumNotFound = New("Album not found", http.StatusNotFound)
	// AlbumNothingToUpdate error
	AlbumNothingToUpdate = New("Nothing to update on album", http.StatusBadRequest)

	// TokenAccessNotGenerated error
	TokenAccessNotGenerated = New("access token not generated", http.StatusUnauthorized)
	// TokenExpiredAccess error
	TokenExpiredAccess = New("expired access token", http.StatusUnauthorized)
	// TokenInvalidAccess error
	TokenInvalidAccess = New("invalid access token", http.StatusUnauthorized)
	// TokenInvalidExpiry error
	TokenInvalidExpiry = New("expiry should be greater than 0", http.StatusUnauthorized)
	// TokenInvalidRefresh error
	TokenInvalidRefresh = New("invalid refresh token", http.StatusUnauthorized)
	// TokenRefreshNotGenerated error
	TokenRefreshNotGenerated = New("refresh token not generated", http.StatusUnauthorized)

	// ClientAlreadyExists error
	ClientAlreadyExists = New("client already exists", http.StatusBadRequest)
	// ClientAPIKeyNotGenerated error
	ClientAPIKeyNotGenerated = New("API key not generated", http.StatusBadRequest)
	// ClientInvalidAPIKey error
	ClientInvalidAPIKey = New("invalid API key", http.StatusUnauthorized)
	// ClientInvalidInfo error
	ClientInvalidInfo = New("invalid client info", http.StatusBadRequest)
	// ClientInvalidKey error
	ClientInvalidKey = New("invalid client id", http.StatusUnauthorized)
	// ClientInvalidSecretToken error
	ClientInvalidSecretToken = New("invalid secret token", http.StatusUnauthorized)
	// ClientNotFound error
	ClientNotFound = New("client not found", http.StatusNotFound)
	// ClientSecretTokenNotGenerated error
	ClientSecretTokenNotGenerated = New("secret token not generated", http.StatusBadRequest)

	// DatabaseNotImplemented error
	DatabaseNotImplemented = New("database not implemented", http.StatusBadRequest)
)

// FileAudioSizeExceeded error
func FileAudioSizeExceeded(limit int64) error {
	msg := fmt.Sprintf("audio file size limit is %v MB", limit/1048576)
	return New(msg, http.StatusBadRequest)
}

// FileInvalidAudioExtension error
func FileInvalidAudioExtension(ext []string) error {
	msg := fmt.Sprintf("only accept audio file with extension %v", ext)
	return New(msg, http.StatusBadRequest)
}

// FileInvalidImageExtension error
func FileInvalidImageExtension(ext []string) error {
	msg := fmt.Sprintf("only accept audio file with extension %v", ext)
	return New(msg, http.StatusBadRequest)
}

// FileImageSizeExceeded error
func FileImageSizeExceeded(limit int64) error {
	msg := fmt.Sprintf("image file size limit is %v MB", limit/1048576)
	return New(msg, http.StatusBadRequest)
}

// Err struct
type Err struct {
	Message    string `json:"error"`
	HTTPStatus int    `json:"-"`
}

func (e Err) Error() string {
	return e.Message
}

// Data method
func (e Err) Data() []byte {
	data, err := json.Marshal(e)

	if err != nil {
		return []byte(`{"error":"unknown"}`)
	}

	return data
}

// New method
func New(m string, c int) Err {
	return Err{
		Message:    m,
		HTTPStatus: c,
	}
}

// Internal method
func Internal(err error) Err {
	log.Println("internal:", err)
	return New("something went wrong internally", http.StatusInternalServerError)
}
