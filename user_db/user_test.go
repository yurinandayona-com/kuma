package user_db

import (
	"testing"
)

func TestUser(t *testing.T) {
	u := &User{
		Name: "MakeNowJust",
		ID: "00000000-0000-0000-0000-000000000000",
	}

	if name := u.GetName(); name != "MakeNowJust" {
		t.Errorf("unexpected u.GetName() result: %#v", name)
	}

	if id := u.GetID(); id != "00000000-0000-0000-0000-000000000000" {
		t.Errorf("unexpected u.GetID() result: %#v", id)
	}
}
