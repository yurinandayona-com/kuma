package server

type UserVerifier interface {
	Verify(token string) (User, error)
}

type User interface {
	GetID() string
	GetName() string
}
