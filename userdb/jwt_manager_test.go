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

	rs256Token = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.e30.b-gLBOyB62jzeiETcpDg4wgLa9EcJcEN5Dh4Hna5Uvs6wqGWRco1uIxdsQJRTvsWPq63A_ZM9g7rjs-SEORyty1DqWNeqaK3uaECr5n80dL_oKcWUhzCDJbC2W_v4_2jQz4lz5m12FH-_N19RRymA_GeKuZMyvH0MUlitVfnjlA"
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

	for _, c := range []struct {
		token string
		msg   string
	}{
		{
			token: rs256Token,
			msg:   "invalid JWT token: invalid JWT algorithm",
		},
		{
			token: "",
			msg:   "invalid JWT",
		},
	} {
		user, err := jm.Verify(c.token)
		if err == nil {
			t.Fatalf("unexpected ok on Verify: %#v", user)
		}
		if msg := err.Error(); msg != c.msg {
			t.Errorf("unexpected error message: %#v", msg)
		}
	}
}
