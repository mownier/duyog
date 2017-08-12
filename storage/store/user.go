package store

// UserRepo interface
type UserRepo interface {
	AddAvatar(k string) (string, error)
	Verify(k, fk string) error
}

// AddUserAvatar method
func AddUserAvatar(r UserRepo, k string) (string, error) {
	return r.AddAvatar(k)
}

// VerifyUser method
func VerifyUser(r UserRepo, k, fk string) error {
	return r.Verify(k, fk)
}
