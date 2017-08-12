package store

// AlbumRepo interface
type AlbumRepo interface {
	AddPhoto(k string) (string, error)
	Verify(k, fk string) error
}

// AddAlbumPhoto method
func AddAlbumPhoto(r AlbumRepo, k string) (string, error) {
	return r.AddPhoto(k)
}

// VerifyAlbum method
func VerifyAlbum(r AlbumRepo, k, fk string) error {
	return r.Verify(k, fk)
}
