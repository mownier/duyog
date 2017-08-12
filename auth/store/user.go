package store

// UserRepo interface
type UserRepo interface {
	Create(u User) (User, error)
	ChangePass(i ChangePassInput) error
	ValidateCredential(u User) (User, error)
	GetByKey(k string) (User, error)
}

// CreateUser method
func CreateUser(r UserRepo, u User) (User, error) {
	return r.Create(u)
}

// ChangeUserPass method
func ChangeUserPass(r UserRepo, i ChangePassInput) error {
	return r.ChangePass(i)
}

// ValidateUserCredential method
func ValidateUserCredential(r UserRepo, u User) (User, error) {
	return r.ValidateCredential(u)
}

// GetUserByKey method
func GetUserByKey(r UserRepo, k string) (User, error) {
	return r.GetByKey(k)
}
