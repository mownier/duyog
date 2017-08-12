package store

// AlbumRepo interface
type AlbumRepo interface {
	Create(a Album) (Album, error)

	Update(a Album) (Album, error)

	GetByKey(k AlbumKey) (Album, error)
	GetSongs(k AlbumKey) (Songs, error)
	GetArtists(k AlbumKey) (Artists, error)
}

// CreateAlbum method
func CreateAlbum(r AlbumRepo, a Album) (Album, error) {
	return r.Create(a)
}

// UpdateAlbum method
func UpdateAlbum(r AlbumRepo, a Album) (Album, error) {
	return r.Update(a)
}

// GetAlbumByKey method
func GetAlbumByKey(r AlbumRepo, k AlbumKey) (Album, error) {
	return r.GetByKey(k)
}

// GetAlbumSongs method
func GetAlbumSongs(r AlbumRepo, k AlbumKey) (Songs, error) {
	return r.GetSongs(k)
}

// GetAlbumArtists method
func GetAlbumArtists(r AlbumRepo, k AlbumKey) (Artists, error) {
	return r.GetArtists(k)
}
