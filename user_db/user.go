package user_db

// User means a user of user DB.
//
// It implements server.User.
type User struct {
	ID   string `toml:"id"`
	Name string `toml:"name"`
}

// It is the same as u.ID.
func (u *User) GetID() string {
	return u.ID
}

// It is the same as u.Name.
func (u *User) GetName() string {
	return u.Name
}
