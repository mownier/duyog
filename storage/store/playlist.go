package store

// PlaylistRepo interface
type PlaylistRepo interface {
	AddPhoto(k string) (string, error)
	Verify(k, fk string) error
}

// AddPlaylistPhoto method
func AddPlaylistPhoto(r PlaylistRepo, k string) (string, error) {
	return r.AddPhoto(k)
}

// VerifyPlaylist method
func VerifyPlaylist(r PlaylistRepo, k, fk string) error {
	return r.Verify(k, fk)
}
