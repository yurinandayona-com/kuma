package user_db

type User struct {
	ID   string `toml:"id"`
	Name string `toml:"name"`
}

func (u *User) GetID() string {
	return u.ID
}

func (u *User) GetName() string {
	return u.Name
}
