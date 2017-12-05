package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"

	"github.com/pkg/errors"
	"github.com/yurinandayona-com/kuma/api"
	"github.com/yurinandayona-com/kuma/version"
)

type hub struct {
	sync.RWMutex

	Host   string
	Server *Server
	User   User
	Stream api.Hub_ConnectServer

	closed bool

	nextTunnelID int64
	tunnels      map[int64]*tunnel
}

func newHub(host string, server *Server, user User, stream api.Hub_ConnectServer) *hub {
	return &hub{
		Host:   host,
		Server: server,
		User:   user,
		Stream: stream,

		tunnels: make(map[int64]*tunnel, 0),
	}
}

func (hub *hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// add headers for forward proxy
	r.Header.Add("X-Forwarded-For", r.RemoteAddr)
	r.Header.Add("X-Forwarded-Host", r.Host)
	r.Header.Add("X-Forwarded-Proto", "http")
	r.Header.Add("Via", fmt.Sprintf("%s %s", r.Proto, version.Full))

	tunnel, closeTunnel, ok := hub.openTunnel(w, r)
	if !ok {
		return
	}
	defer closeTunnel()

	ctx := r.Context()
	select {
	case <-tunnel.C:
	case <-ctx.Done():
	}
}

func (hub *hub) openTunnel(w http.ResponseWriter, r *http.Request) (*tunnel, func(), bool) {
	hub.Lock()
	defer hub.Unlock()

	if hub.closed {
		http.Error(w, "kuma: hub is closed", 500)
		return nil, nil, false
	}

	tunnelID := hub.nextTunnelID
	hub.nextTunnelID++

	tunnelIDHash, err := encodeTunnelID(hub.Server.HashID, tunnelID)
	if err != nil {
		log.Printf("error: %s", err)
		http.Error(w, "kuma: invalid tunnelID", 400)
		return nil, nil, false
	}

	err = hub.Stream.Send(&api.Request{Host: hub.Host, TunnelID: tunnelIDHash})
	if err != nil {
		log.Printf("error: %s", err)
		http.Error(w, "kuma: failed to open a new tunnel", 500)
		return nil, nil, false
	}

	tunnel := newTunnel(tunnelID, tunnelIDHash, hub, w, r)
	hub.tunnels[tunnelID] = tunnel

	closeTunnel := func() {
		hub.Lock()
		defer hub.Unlock()

		tunnel.Close()
		delete(hub.tunnels, tunnelID)
	}

	return tunnel, closeTunnel, true
}

func (hub *hub) Close() {
	hub.Lock()
	defer hub.Unlock()

	if hub.closed {
		return
	}

	for _, tunnel := range hub.tunnels {
		tunnel.Close()
	}

	hub.closed = true
}

func (hub *hub) ReceiveHeader(tunnelID int64) (*api.RequestHeader, error) {
	tunnel, err := hub.getTunnel(tunnelID)
	if err != nil {
		return nil, err
	}

	headers := make([]*api.Header, 0, len(tunnel.Request.Header))
	for name, values := range tunnel.Request.Header {
		headers = append(headers, &api.Header{
			Name:   name,
			Values: values,
		})
	}

	reqHeader := &api.RequestHeader{
		Method:  tunnel.Request.Method,
		Path:    tunnel.Request.RequestURI,
		Headers: headers,
	}

	log.Printf("debug: Tunnel.ReceiveHeader: %s [%s]: %s %s", hub.Host, tunnel.IDHash, reqHeader.Method, reqHeader.Path)

	return reqHeader, nil
}

func (hub *hub) ReceiveBody(tunnelID int64, stream api.Tunnel_ReceiveBodyServer) error {
	tunnel, err := hub.getTunnel(tunnelID)
	if err != nil {
		return err
	}

	writer := &receiveBodyWriter{stream: stream}
	defer tunnel.Request.Body.Close()

	n, err := io.Copy(writer, tunnel.Request.Body)

	log.Printf("debug: Tunnel.ReceiveBody: %s [%s]: %d bytes", hub.Host, tunnel.IDHash, n)

	return err
}

type receiveBodyWriter struct {
	stream api.Tunnel_ReceiveBodyServer
}

func (rbw *receiveBodyWriter) Write(buf []byte) (int, error) {
	err := rbw.stream.Send(&api.RequestBody{Body: buf})
	if err != nil {
		return 0, err
	}

	return len(buf), nil
}

func (hub *hub) SendHeader(tunnelID int64, resHeader *api.ResponseHeader) error {
	tunnel, err := hub.getTunnel(tunnelID)
	if err != nil {
		return err
	}

	err = tunnel.UpdateStatus(tunnelStatusInitialized, tunnelStatusSentHeader)
	if err != nil {
		return err
	}

	header := tunnel.Response.Header()
	for _, h := range resHeader.Headers {
		header[h.Name] = h.Values
	}
	header.Add("Via", fmt.Sprintf("%s %s", tunnel.Request.Proto, version.Full))

	tunnel.Response.WriteHeader(int(resHeader.Status))

	log.Printf("debug: Tunnel.SendHeader: %s [%s]: %d", hub.Host, tunnel.IDHash, resHeader.Status)

	return nil
}

func (hub *hub) SendBody(tunnelID int64, stream api.Tunnel_SendBodyServer) error {
	tunnel, err := hub.getTunnel(tunnelID)
	if err != nil {
		return err
	}

	err = tunnel.UpdateStatus(tunnelStatusSentHeader, tunnelStatusSentBody)
	if err != nil {
		return err
	}
	defer tunnel.Close()

	n := 0
	var outerErr error
	go func() {
		for {
			resBody, err := stream.Recv()
			if err != nil {
				if err != io.EOF {
					outerErr = err
				}
				tunnel.Close()
				return
			}

			m, err := tunnel.Response.Write(resBody.Body)
			n += m
			if err != nil {
				outerErr = err
				tunnel.Close()
				return
			}
		}
	}()

	ctx := stream.Context()
	select {
	case <-ctx.Done():
	case <-tunnel.C:
	}

	log.Printf("debug: Tunnel.SendBody: %s [%s]: %d bytes", hub.Host, tunnel.IDHash, n)

	return outerErr
}

func (hub *hub) SendError(tunnelID int64, resErr *api.ResponseError) error {
	tunnel, err := hub.getTunnel(tunnelID)
	if err != nil {
		return err
	}

	log.Printf("debug: Tunnel.SendError: %s [%s]: %s", hub.Host, tunnel.IDHash, resErr.Error)

	tunnel.Close()
	http.Error(tunnel.Response, "kuma: internal server error", 500)

	return nil
}

func (hub *hub) getTunnel(tunnelID int64) (*tunnel, error) {
	hub.RLock()
	defer hub.RUnlock()

	tunnel, ok := hub.tunnels[tunnelID]
	if !ok {
		return nil, errors.New("tunnel not found")
	}

	return tunnel, nil
}
