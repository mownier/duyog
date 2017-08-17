// Package service provides the authentication server with
// the HTTP handler functions for token and user resources.
// Also, it provides a microservice using JSON RPC to verify
// whether the access token is valid or not.
package service

// pathCode identifies an HTTP handler function
type pathCode int

const (
	// pathCodeUserSignIn ties with the HTTP
	// handler function on signing in a user
	pathCodeUserSignIn pathCode = 1000

	// pathCodeUserSignUp ties with the HTTP
	// handler function on signing up a user
	pathCodeUserSignUp pathCode = 1001

	// pathCodeUserChangePass ties with the HTTP
	// handler function on changing user password
	pathCodeUserChangePass pathCode = 1002

	// pathCodeTokenRefresh ties with the HTTP
	// handler function on refreshing an access token
	pathCodeTokenRefresh pathCode = 2000
)

// forbiddenPath holds the information whether a certain
// HTTP handler function forbidden or not.
// This will be overruled if the role of a client is admin.
// Client's role will be retrieved by decoding the
// access token which is a JSON Web Token.
var forbiddenPath = map[pathCode]bool{
	pathCodeUserSignIn:     false,
	pathCodeUserSignUp:     false,
	pathCodeUserChangePass: false,

	pathCodeTokenRefresh: false,
}

// isForbidden determines whether a particular pathCode
// is forbidden or not.
func isForbidden(c pathCode) bool {
	return forbiddenPath[c]
}
