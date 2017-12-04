package server

import (
	"net/http"
	"sync"

	"github.com/pkg/errors"
)

type tunnelStatus int

const (
	tunnelStatusInitialized tunnelStatus = iota
	tunnelStatusSentHeader
	tunnelStatusSentBody
	tunnelStatusClosed
)

type tunnel struct {
	sync.Mutex

	ID       int64
	IDHash   string
	Hub      *hub
	Response http.ResponseWriter
	Request  *http.Request

	C chan struct{}

	status tunnelStatus
}

func newTunnel(id int64, idHash string, hub *hub, w http.ResponseWriter, r *http.Request) *tunnel {
	return &tunnel{
		ID:       id,
		IDHash:   idHash,
		Hub:      hub,
		Response: w,
		Request:  r,

		C: make(chan struct{}),
	}
}

func (t *tunnel) UpdateStatus(from, to tunnelStatus) error {
	t.Lock()
	defer t.Unlock()

	if t.status != from {
		return errors.New("kuma: invalid status")
	}
	t.status = to

	return nil
}

func (t *tunnel) Close() {
	t.Lock()
	defer t.Unlock()

	if t.status == tunnelStatusClosed {
		return
	}

	close(t.C)
	t.status = tunnelStatusClosed
}
