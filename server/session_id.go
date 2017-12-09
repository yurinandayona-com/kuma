package server

import (
	"github.com/pkg/errors"
)

const (
	sessionIDSign = 8342 // BEAR
)

// encodeSessionID encodes hubID and tunnelID by using hashID.
func (svr *Server) encodeSessionID(hubID, tunnelID int64) (string, error) {
	encoded, err := svr.HashID.EncodeInt64([]int64{sessionIDSign, svr.uniq, hubID, tunnelID})
	if err != nil {
		return "", errors.Wrap(err, "failed to encode session ID")
	}

	return encoded, nil
}

// decodeSessionID decodes tunnelIDHash by using hashID.
func (svr *Server) decodeSessionID(sessionIDHash string) (hubID int64, tunnelID int64, err error) {
	decoded := svr.HashID.DecodeInt64(sessionIDHash)

	if len(decoded) == 0 || len(decoded) != 4 || decoded[0] != sessionIDSign || decoded[1] != svr.uniq {
		err = errors.New("invalid session ID")
	} else {
		hubID = decoded[2]
		tunnelID = decoded[3]
	}

	return
}
