package store

// ArtistRepo interface
type ArtistRepo interface {
	Create(a Artist) (Artist, error)

	Update(a Artist) (Artist, error)

	GetByKey(k ArtistKey) (Artist, error)
	GetSongs(k ArtistKey) (Songs, error)
	GetAlbums(k ArtistKey) (Albums, error)
}

// CreateArtist method
func CreateArtist(r ArtistRepo, a Artist) (Artist, error) {
	return r.Create(a)
}

// UpdateArtist method
func UpdateArtist(r ArtistRepo, a Artist) (Artist, error) {
	return r.Update(a)
}

// GetArtistByKey method
func GetArtistByKey(r ArtistRepo, k ArtistKey) (Artist, error) {
	return r.GetByKey(k)
}

// GetArtistSongs method
func GetArtistSongs(r ArtistRepo, k ArtistKey) (Songs, error) {
	return r.GetSongs(k)
}

// GetArtistAlbums method
func GetArtistAlbums(r ArtistRepo, k ArtistKey) (Albums, error) {
	return r.GetAlbums(k)
}
