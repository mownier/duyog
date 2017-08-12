package store

// SongRepo interface
type SongRepo interface {
	Create(s Song, ar []ArtistKey, al []AlbumKey) (Songs, error)

	Update(s Song) (Song, error)
	UpdateAlbums(k SongKey, a []AlbumKey) ([]AlbumKey, error)
	UpdateArtists(k SongKey, a []ArtistKey) ([]ArtistKey, error)

	GetByKey(k SongKey) (Songs, error)
}

// CreateSong method
func CreateSong(r SongRepo, s Song, ar []ArtistKey, al []AlbumKey) (Songs, error) {
	return r.Create(s, ar, al)
}

// UpdateSong method
func UpdateSong(r SongRepo, s Song) (Song, error) {
	return r.Update(s)
}

// UpdateSongAlbums method
func UpdateSongAlbums(r SongRepo, k SongKey, a []AlbumKey) ([]AlbumKey, error) {
	return r.UpdateAlbums(k, a)
}

// UpdateSongArtists method
func UpdateSongArtists(r SongRepo, k SongKey, a []ArtistKey) ([]ArtistKey, error) {
	return r.UpdateArtists(k, a)
}

// GetSongByKey method
func GetSongByKey(r SongRepo, k SongKey) (Songs, error) {
	return r.GetByKey(k)
}
