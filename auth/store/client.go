package store

// ClientRepo interface
type ClientRepo interface {
	Create(c Client) (Client, error)
	GetByKey(k string) (Client, error)
	GetByAPIKey(k string) (Client, error)
	GetBySecretToken(t string) (Client, error)
	GetByAccessToken(t string) (Client, error)
}

// CreateClient method
func CreateClient(r ClientRepo, c Client) (Client, error) {
	return r.Create(c)
}

// GetClientByKey method
func GetClientByKey(r ClientRepo, k string) (Client, error) {
	return r.GetByKey(k)
}

// GetClientByAPIKey method
func GetClientByAPIKey(r ClientRepo, k string) (Client, error) {
	return r.GetByAPIKey(k)
}

// GetClientBySecretToken method
func GetClientBySecretToken(r ClientRepo, t string) (Client, error) {
	return r.GetBySecretToken(t)
}

// GetClientByAccessToken method
func GetClientByAccessToken(r ClientRepo, t string) (Client, error) {
	return r.GetByAccessToken(t)
}
