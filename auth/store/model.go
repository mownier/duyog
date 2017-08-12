package store

// Client struct
type Client struct {
	Key         string `json:"id,omitempty" redis:"id"`
	Role        string `json:"role,omitempty" redis:"role"`
	Name        string `json:"name,omitempty" redis:"name"`
	Email       string `json:"email,omitempty" redis:"email"`
	APIKey      string `json:"api_key,omitempty" redis:"api_key"`
	SecretToken string `json:"secret_token,omitempty" redis:"secret_token"`
}

// User struct
type User struct {
	Key      string `json:"id,omitempty" redis:"id"`
	Email    string `json:"email,omitempty" redis:"email"`
	Password string `json:"password,omitempty" redis:"password"`
}

// ChangePassInput struct
type ChangePassInput struct {
	Email   string
	NewPass string
	CurPass string
}

// Token struct
type Token struct {
	Key     string `json:"id,omitempty" redis:"id"`
	Access  string `json:"access_token" redis:"access_token"`
	Refresh string `json:"refresh_token" redis:"refresh_token"`
	Expiry  int64  `json:"expiry" redis:"expiry"`
}

// RefreshTokenInput struct
type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token"`
	Expiry       int64  `json:"expiry,omitempty"`
}
