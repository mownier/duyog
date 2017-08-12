package store

// UserRepo interface
type UserRepo interface {
	Register(u User) (User, error)

	Update(u User) (User, error)

	AddPlaylists(k UserKey, p []PlaylistKey) ([]PlaylistKey, error)

	GetByKey(k UserKey) (User, error)

	HasPlaylist(uk UserKey, pk PlaylistKey) error
}

// RegisterUser method
func RegisterUser(r UserRepo, u User) (User, error) {
	return r.Register(u)
}

// UpdateUser method
func UpdateUser(r UserRepo, u User) (User, error) {
	return r.Update(u)
}

// AddUserPlaylists method
func AddUserPlaylists(r UserRepo, k UserKey, p []PlaylistKey) ([]PlaylistKey, error) {
	return r.AddPlaylists(k, p)
}

// GetUserByKey method
func GetUserByKey(r UserRepo, k UserKey) (User, error) {
	return r.GetByKey(k)
}

// UserHasPlaylist method
func UserHasPlaylist(r UserRepo, uk UserKey, pk PlaylistKey) error {
	return r.HasPlaylist(uk, pk)
}
