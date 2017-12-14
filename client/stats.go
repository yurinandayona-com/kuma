package client

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
)

type stats struct {
	SessionID string

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

	fmt.Fprintf(line, "info: [%s]", ss.SessionID)
	fmt.Fprintf(line, " %s (%6s) %10s", color.New(statusColor(ss.Status)).Sprintf(" %d ", ss.Status), humanize.Bytes(uint64(ss.ResponseBodySize)), d)
	fmt.Fprintf(line, " | %s %s", color.New(methodColor(ss.Method), color.Bold).Sprintf("%s", ss.Method), color.New(color.Bold).Sprintf("%s", ss.Path))
	if ss.RequestBodySize > 0 {
		fmt.Fprintf(line, " (%s)", humanize.Bytes(uint64(ss.RequestBodySize)))
	}

	log.Print(line)
}

func statusColor(code int) color.Attribute {
	switch {
	case 200 <= code && code < 300:
		return color.BgHiGreen
	case 300 <= code && code < 400:
		return color.BgHiWhite
	case 400 <= code && code < 500:
		return color.BgHiRed
	default:
		return color.BgHiYellow
	}
}

func methodColor(method string) color.Attribute {
	switch method {
	case "GET":
		return color.FgBlue
	case "POST":
		return color.FgCyan
	case "PUT":
		return color.FgYellow
	case "DELETE":
		return color.FgRed
	case "PATCH":
		return color.FgGreen
	case "HEAD":
		return color.FgMagenta
	case "OPTIONS":
		return color.FgWhite
	default:
		return color.Reset
	}
}
