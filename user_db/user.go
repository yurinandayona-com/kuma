package user_db

import (
	"github.com/yurinandayona-com/kuma/server"
)

// User means a user of user DB.
//
// It implements server.User.
type User struct {
	ID   string `toml:"id"`
	Name string `toml:"name"`
}

var _ server.User = &User{}

// GetID is the same as u.ID.
func (u *User) GetID() string {
	return u.ID
}

// GetName is the same as u.Name.
func (u *User) GetName() string {
	return u.Name
}
