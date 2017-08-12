package store

// ArtistRepo interface
type ArtistRepo interface {
	AddPhoto(k string) (string, error)
	Verify(k, fk string) error
}

// AddArtistPhoto method
func AddArtistPhoto(r ArtistRepo, k string) (string, error) {
	return r.AddPhoto(k)
}

// VerifyArtist method
func VerifyArtist(r ArtistRepo, k, fk string) error {
	return r.Verify(k, fk)
}
