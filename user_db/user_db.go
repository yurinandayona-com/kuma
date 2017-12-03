package user_db

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"io/ioutil"
)

type UserDB struct {
	users map[string]*User
}

func LoadUserDB(filename string) (*UserDB, error) {
	p, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "kuma: failed to load user DB file")
	}

	var internalUserDB struct {
		Users []*User `toml:"users"`
	}

	err = toml.Unmarshal(p, &internalUserDB)
	if err != nil {
		return nil, errors.Wrapf(err, "kuma: failed to load user DB TOML")
	}

	users := make(map[string]*User, len(internalUserDB.Users))

	for _, u := range internalUserDB.Users {
		if _, ok := users[u.ID]; ok {
			return nil, errors.Errorf("kuma: duplicated user ID: %v", u.ID)
		}

		users[u.ID] = u
	}

	return &UserDB{users: users}, nil
}

func (udb *UserDB) Verify(id, name string) (*User, error) {
	u, ok := udb.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}

	if name != u.Name {
		return nil, errors.New("invalid name")
	}

	return u, nil
}

func (udb *UserDB) GetUsers() []*User {
	users := make([]*User, 0, len(udb.users))
	for _, u := range udb.users {
		users = append(users, u)
	}
	return users
}
