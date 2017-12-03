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

type JWTManager struct {
	UserDB  *UserDB
	HMACKey []byte
}

type jwtUserClaims struct {
	jwt.StandardClaims

	ID   string `json:"https://github.com/yurinandayona-com/kuma/claim-types/user-id"`
	Name string `json:"https://github.com/yurinandayona-com/kuma/claim-types/user-name"`
}

func (jm *JWTManager) Verify(t string) (server.User, error) {
	token, err := jwt.ParseWithClaims(t, &jwtUserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid JWT algorithm")
		}
		return jm.HMACKey, nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "kuma: invalid JWT token")
	}

	if claims, ok := token.Claims.(*jwtUserClaims); ok && token.Valid {
		return jm.UserDB.Verify(claims.ID, claims.Name)
	} else {
		return nil, errors.New("kuma: invalid JWT")
	}
}

func (jm *JWTManager) Sign(u *User) (string, error) {
	u, err := jm.UserDB.Verify(u.ID, u.Name)
	if err != nil {
		return "", err
	}

	claims := &jwtUserClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenTimeout).Unix(),
		},
		ID:   u.ID,
		Name: u.Name,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jm.HMACKey)
}
