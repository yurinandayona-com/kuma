package userdb

import (
	"strings"
	"testing"
)

func TestLoadUserDB(t *testing.T) {
	{
		userDB, err := LoadUserDB("testdata/not_found.toml")
		if err == nil {
			t.Fatalf("unexpected ok on LoadUserDB: %#v", userDB)
		}
		if msg := err.Error(); !strings.HasPrefix(msg, "failed to load user DB file:") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	}

	{
		userDB, err := LoadUserDB("testdata/invalid.toml")
		if err == nil {
			t.Fatalf("unexpected ok on LoadUserDB: %#v", userDB)
		}
		if msg := err.Error(); !strings.HasPrefix(msg, "failed to load user DB TOML:") {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	}

	{
		userDB, err := LoadUserDB("testdata/conflicted.toml")
		if err == nil {
			t.Fatalf("unexpected ok on LoadUserDB: %#v", userDB)
		}
		if msg := err.Error(); !strings.HasPrefix(msg, "user ID conflicted:") {
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
		u, err := userDB.Verify(user.ID, user.Name)
		if err != nil {
			t.Fatalf("unexpected error on Verify: %+#v", err)
		}
		if *u != *user {
			t.Fatalf("unexpected verify result: %#v", u)
		}
	}

	mnj := find("MakeNowJust")
	sh := find("sh4869")
	if len(users) != 2 {
		t.Fatalf("unexpected users size: %d", len(users))
	}

	verify(mnj)
	verify(sh)

	{
		u, err := userDB.Verify("XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX", "404_not_found")
		if err == nil {
			t.Fatalf("unexpected ok on Verify: %#v", u)
		}
		if msg := err.Error(); msg != "user not found" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	}

	{
		u, err := userDB.Verify(sh.ID, mnj.Name)
		if err == nil {
			t.Fatalf("unexpected ok on Verify: %#v", u)
		}
		if msg := err.Error(); msg != "invalid name" {
			t.Fatalf("unexpected error message: %#v", msg)
		}
	}
}
