// Package user_db provides user DB and verifier implementation.
//
// And it also provides JWT manager to verify user token.
package user_db

import (
	"io/ioutil"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// UserDB means user DB.
//
// It is the User map indexed by user IDs internally.
type UserDB struct {
	users map[string]*User
}

// LoadUserDB loads user DB from TOML.
//
// TOML format looks like this:
//
//     # User DB TOML should contain users array which has two keys, name
//     # and id.
//
//     # 1st user
//     [[users]]
//       name = "MakeNowJust"
//       id   = "00000000-0000-0000-0000-000000000000"
//
//     # 2nd user
//     [[users]]
//       name = "sh4869"
//       id   = "11111111-1111-1111-1111-111111111111"
//
//     # and more...
func LoadUserDB(filename string) (*UserDB, error) {
	p, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load user DB file")
	}

	var internalUserDB struct {
		Users []*User `toml:"users"`
	}

	err = toml.Unmarshal(p, &internalUserDB)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load user DB TOML")
	}

	users := make(map[string]*User, len(internalUserDB.Users))

	for _, u := range internalUserDB.Users {
		if _, ok := users[u.ID]; ok {
			return nil, errors.Errorf("user ID conflicted: %v", u.ID)
		}

		users[u.ID] = u
	}

	return &UserDB{users: users}, nil
}

// Verify finds the user having the given id and checks its name is the
// same as the given name.
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

// GetUsers returns users as slice of this user DB.
func (udb *UserDB) GetUsers() []*User {
	users := make([]*User, 0, len(udb.users))
	for _, u := range udb.users {
		users = append(users, u)
	}
	return users
}
