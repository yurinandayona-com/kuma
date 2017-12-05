package server

// UserVerifier is to use verifing the token of request.
type UserVerifier interface {
	Verify(token string) (User, error)
}

// User is user.
type User interface {
	GetID() string
	GetName() string
}
