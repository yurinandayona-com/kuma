package client

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/dustin/go-humanize"
)

type stats struct {
	TunnelID string

	Method string
	Path   string

	RequestBodySize int64

	Status int

	ResponseBodySize int64
}

func (ss *stats) Log(d time.Duration) {
	if ss.Method == "" {
		ss.Method = "???"
	}

	if ss.Path == "" {
		ss.Path = "???"
	}

	if ss.Status == 0 {
		ss.Status = 500
	}

	line := &bytes.Buffer{}
	fmt.Fprintf(line, "info: [%s] %s %s", ss.TunnelID, ss.Method, ss.Path)
	if ss.RequestBodySize > 0 {
		fmt.Fprintf(line, " (%s)", humanize.Bytes(uint64(ss.RequestBodySize)))
	}
	fmt.Fprintf(line, " ==> %d (%s) in %s", ss.Status, humanize.Bytes(uint64(ss.ResponseBodySize)), d)

	log.Print(line)
}
