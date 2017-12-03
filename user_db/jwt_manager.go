package user_db

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"github.com/yurinandayona-com/kuma/server"
	"time"
)

const (
	TokenTimeout = 24 * time.Hour * 100
)

// JWTManager is token manager of user DB.
//
// It implements server.UserVerifier.
type JWTManager struct {
	UserDB  *UserDB
	HMACKey []byte
}

// JWTUserClaims is JWT claims containing user information.
type JWTUserClaims struct {
	jwt.StandardClaims

	ID   string `json:"https://github.com/yurinandayona-com/kuma/claim-types/user-id"`
	Name string `json:"https://github.com/yurinandayona-com/kuma/claim-types/user-name"`
}

// Verify verifies t as JWT and then returns a user bound this JWT or error.
func (jm *JWTManager) Verify(t string) (server.User, error) {
	claims, valid, err := jm.Parse(t)
	if err != nil {
		return nil, err
	}

	if valid {
		return jm.UserDB.Verify(claims.ID, claims.Name)
	} else {
		return nil, errors.New("kuma: invalid JWT")
	}
}

// Parse parses t as JWT and then returns JWTUserClaims bound this JWT and
// flag which is token validation status and error.
func (jm *JWTManager) Parse(t string) (*JWTUserClaims, bool, error) {
	token, err := jwt.ParseWithClaims(t, &JWTUserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid JWT algorithm")
		}
		return jm.HMACKey, nil
	})
	if err != nil {
		return nil, false, errors.Wrap(err, "kuma: invalid JWT token")
	}

	if claims, ok := token.Claims.(*JWTUserClaims); ok {
		return claims, token.Valid, nil
	} else {
		return nil, false, errors.New("kuma: invalid JWT")
	}
}

func (jm *JWTManager) Sign(u *User) (string, error) {
	u, err := jm.UserDB.Verify(u.ID, u.Name)
	if err != nil {
		return "", err
	}

	claims := &JWTUserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenTimeout).Unix(),
		},
		ID:   u.ID,
		Name: u.Name,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jm.HMACKey)
}
