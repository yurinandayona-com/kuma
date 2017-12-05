package server

import (
	"github.com/pkg/errors"
	"github.com/speps/go-hashids"
)

const (
	tunnelIDSign = 8342 // BEAR
)

// encodeTunnelID encodes tunnelID by using hashids.
func encodeTunnelID(hashID *hashids.HashID, tunnelID int64) (string, error) {
	encoded, err := hashID.EncodeInt64([]int64{tunnelIDSign, tunnelID})
	if err != nil {
		return "", errors.Wrap(err, "failed to encode tunnel ID")
	}

	return encoded, nil
}

// decodeTunnelID decodes tunnelIDHash by using hashIDs.
func decodeTunnelID(hashID *hashids.HashID, tunnelIDHash string) (int64, error) {
	decoded := hashID.DecodeInt64(tunnelIDHash)

	if len(decoded) == 0 || len(decoded) != 2 || decoded[0] != tunnelIDSign {
		return 0, errors.New("invalid tunnel ID")
	}

	return decoded[1], nil
}
