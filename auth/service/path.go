package service

type pathCode int

const (
	// pathCodeUserSignIn code
	pathCodeUserSignIn pathCode = 1000

	// pathCodeUserSignUp code
	pathCodeUserSignUp pathCode = 1001

	// pathCodeUserChangePass code
	pathCodeUserChangePass pathCode = 1002

	// pathCodeTokenRefresh code
	pathCodeTokenRefresh pathCode = 2000
)

var forbiddenPath = map[pathCode]bool{
	pathCodeUserSignIn:     false,
	pathCodeUserSignUp:     false,
	pathCodeUserChangePass: false,

	pathCodeTokenRefresh: false,
}

func isForbidden(c pathCode) bool {
	return forbiddenPath[c]
}
