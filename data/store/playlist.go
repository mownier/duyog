package store

// PlaylistRepo interface
type PlaylistRepo interface {
	Create(k UserKey, p Playlist) (Playlist, error)

	Update(p Playlist) (Playlist, error)

	AddSongs(pk PlaylistKey, sk []SongKey) ([]SongKey, error)

	GetSongs(k PlaylistKey) (Songs, error)
	GetByKey(k PlaylistKey) (Playlist, error)
	GetByUser(k UserKey) (Playlists, error)
}

// CreatePlaylist method
func CreatePlaylist(r PlaylistRepo, k UserKey, p Playlist) (Playlist, error) {
	return r.Create(k, p)
}

// UpdatePlaylist method
func UpdatePlaylist(r PlaylistRepo, p Playlist) (Playlist, error) {
	return r.Update(p)
}

// AddPlaylistSongs method
func AddPlaylistSongs(r PlaylistRepo, pk PlaylistKey, sk []SongKey) ([]SongKey, error) {
	return r.AddSongs(pk, sk)
}

// GetPlaylistByKey method
func GetPlaylistByKey(r PlaylistRepo, k PlaylistKey) (Playlist, error) {
	return r.GetByKey(k)
}

// GetPlaylistsByUser method
func GetPlaylistsByUser(r PlaylistRepo, k UserKey) (Playlists, error) {
	return r.GetByUser(k)
}

// GetPlaylistSongs method
func GetPlaylistSongs(r PlaylistRepo, k PlaylistKey) (Songs, error) {
	return r.GetSongs(k)
}
