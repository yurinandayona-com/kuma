package userdb

import (
	"strings"
	"testing"
)

func TestLoadUserDB(t *testing.T) {
	for _, c := range []struct {
		filename string
		msg      string
	}{
		{filename: "not_found.toml", msg: "failed to load user DB file:"},
		{filename: "invalid.toml", msg: "failed to load user DB TOML:"},
		{filename: "conflicted.toml", msg: "user ID conflicted:"},
	} {
		userDB, err := LoadUserDB("testdata/" + c.filename)
		if err == nil {
			t.Fatalf("unexpected ok on LoadUserDB: %#v", userDB)
		}
		if msg := err.Error(); !strings.HasPrefix(msg, c.msg) {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	}
}

func TestUserDB(t *testing.T) {
	userDB, err := LoadUserDB("testdata/example.toml")
	if err != nil {
		t.Fatalf("failed to load user DB: %+#v", err)
	}

	users := userDB.GetUsers()
	find := func(name string) *User {
		var found *User
		for _, u := range users {
			if u.Name == name {
				found = u
				break
			}
		}
		if found == nil {
			t.Fatalf("%s not found", name)
		}
		return found
	}
	verify := func(user *User) {

	}

	mnj := find("MakeNowJust")
	sh := find("sh4869")
	if len(users) != 2 {
		t.Fatalf("unexpected users size: %d", len(users))
	}

	verify(mnj)
	verify(sh)

	for _, user := range []*User{mnj, sh} {
		u, err := userDB.Verify(user.ID, user.Name)
		if err != nil {
			t.Fatalf("unexpected error on Verify: %+#v", err)
		}
		if *u != *user {
			t.Fatalf("unexpected verify result: %#v", u)
		}
	}

	for _, c := range []struct {
		id, name string
		msg      string
	}{
		{
			id:   "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX",
			name: "404_not_found",
			msg:  "user not found",
		},
		{
			id:   sh.ID,
			name: mnj.Name,
			msg:  "invalid name",
		},
	} {
		u, err := userDB.Verify(c.id, c.name)
		if err == nil {
			t.Fatalf("unexpected ok on Verify: %#v", u)
		}
		if msg := err.Error(); msg != c.msg {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	}
}
