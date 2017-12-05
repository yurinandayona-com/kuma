package userdb

import (
	"testing"
	"time"
)

var (
	hmacKey = []byte("YURIKUMAAAAA")

	name = "MakeNowJust"
	id   = "00000000-0000-0000-0000-000000000000"

	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOi02MjEzNTU5NjgwMCwiaHR0cHM6Ly9naXRodWIuY29tL3l1cmluYW5kYXlvbmEtY29tL2t1bWEvY2xhaW0tdHlwZXMvdXNlci1pZCI6IjAwMDAwMDAwLTAwMDAtMDAwMC0wMDAwLTAwMDAwMDAwMDAwMCIsImh0dHBzOi8vZ2l0aHViLmNvbS95dXJpbmFuZGF5b25hLWNvbS9rdW1hL2NsYWltLXR5cGVzL3VzZXItbmFtZSI6Ik1ha2VOb3dKdXN0In0.HsoHdMZHsJ-6NL8IVp6iTPZINoow_3yT71Ruz5yIKss"
)

func TestJWTManagerSign(t *testing.T) {
	userDB, err := LoadUserDB("testdata/example.toml")
	if err != nil {
		t.Fatalf("failed to load user DB: %+#v", err)
	}

	jm := &JWTManager{
		UserDB:  userDB,
		HMACKey: hmacKey,
	}

	{
		signed, err := jm.Sign(&User{Name: name, ID: id}, time.Time{})
		if err != nil {
			t.Fatalf("unexpected error on Sign: %+#v", err)
		}

		if signed != token {
			t.Error("unexpected token")
		}
	}

	{
		signed, err := jm.Sign(&User{Name: "not_found"}, time.Time{})
		if err == nil {
			t.Fatalf("unexpected ok on Sign: %#v", signed)
		}

		if msg := err.Error(); msg != "user not found" {
			t.Errorf("unexpected error message: %#v", msg)
		}
	}
}

func TestJWTManagerVerify(t *testing.T) {
	userDB, err := LoadUserDB("testdata/example.toml")
	if err != nil {
		t.Fatalf("failed to load user DB: %+#v", err)
	}

	jm := &JWTManager{
		UserDB:  userDB,
		HMACKey: hmacKey,
	}

	signed, err := jm.Sign(&User{Name: name, ID: id}, time.Now().Add(24*time.Hour))
	if err != nil {
		t.Fatalf("unexpected error on Sign: %+#v", err)
	}

	{
		u, err := jm.Verify(signed)
		if err != nil {
			t.Fatalf("unexpected error on Verify: %+#v", err)
		}

		if n := u.GetName(); n != name {
			t.Errorf("unexpected name: %#v", n)
		}

		if i := u.GetID(); i != id {
			t.Errorf("unexpected id: %#v", i)
		}
	}
}
