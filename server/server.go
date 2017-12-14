// Package server provides kuma server implementation.
package server

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/speps/go-hashids"
	"github.com/yurinandayona-com/kuma/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

var (
	validSubdomain = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
)

// Config is configuration to create server.
type Config struct {
	BaseDomain string
	HashID     *hashids.HashID
}

// Server implements kuma server.
//
// Server pointer implements these interfaces:
//
//    - http.Handler
//    - api.HubServer
//    - api.TunnelServer
//
// So this can be used as these handlers.
type Server struct {
	sync.RWMutex
	*Config

	nextHubID int64
	hubHosts  map[string]int64
	hubs      map[int64]*hub
}

func New(cfg *Config) (*Server, error) {
	svr := &Server{
		Config: cfg,

		hubHosts: make(map[string]int64),
		hubs:     make(map[int64]*hub),
	}

	return svr, nil
}

func (svr *Server) getHubFromHost(host string) (*hub, bool) {
	svr.RLock()
	defer svr.RUnlock()

	hubID, ok := svr.hubHosts[host]
	if !ok {
		return nil, ok
	}

	hub, ok := svr.hubs[hubID]
	return hub, ok
}

func (svr *Server) getHubFromID(hubID int64) (*hub, bool) {
	svr.RLock()
	defer svr.RUnlock()

	hub, ok := svr.hubs[hubID]
	return hub, ok
}

// ServeHTTP handles HTTP request.
func (svr *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	if hub, ok := svr.getHubFromHost(host); ok {
		hub.ServeHTTP(w, r)
	} else {
		http.Error(w, "kuma: host not found", http.StatusNotFound)
	}
}

// Prepare handles Hub.Prepare gRPC call.
func (svr *Server) Prepare(ctx context.Context, config *api.HubConfig) (*api.HubInfo, error) {
	log.Printf("debug: Hub.Prepare: start: %s", config.Subdomain)
	defer log.Printf("debug: Hub.Prepare: end: %s", config.Subdomain)

	if !validSubdomain.MatchString(config.Subdomain) {
		return nil, errors.New("invalid subdomain")
	}

	host := config.Subdomain + "." + svr.BaseDomain
	return &api.HubInfo{Host: host}, nil
}

// Connect handles Hub.Connect gRPC call.
func (svr *Server) Connect(info *api.HubInfo, stream api.Hub_ConnectServer) error {
	log.Printf("debug: Hub.Connect: start: %s", info.Host)
	defer log.Printf("debug: Hub.Connect: end: %s", info.Host)

	ctx := stream.Context()

	if !strings.HasSuffix(info.Host, "."+svr.BaseDomain) {
		return errors.New("invalid subdomain")
	}
	subdomain := strings.TrimSuffix(info.Host, "."+svr.BaseDomain)
	if !validSubdomain.MatchString(subdomain) {
		return errors.New("invalid subdomain")
	}

	host := subdomain + "." + svr.BaseDomain
	closeHub, err := svr.newHub(host, stream)
	if err != nil {
		return err
	}
	defer closeHub()

	<-ctx.Done()

	return nil
}

func (svr *Server) newHub(host string, stream api.Hub_ConnectServer) (func(), error) {
	svr.Lock()
	defer svr.Unlock()

	if _, ok := svr.hubHosts[host]; ok {
		return nil, errors.New("hub is already connected")
	}

	hubID := svr.nextHubID
	svr.nextHubID++

	hub := newHub(hubID, getPeerID(stream.Context()), host, svr, stream)
	svr.hubs[hubID] = hub
	svr.hubHosts[host] = hubID

	closeHub := func() {
		svr.Lock()
		defer svr.Unlock()

		hub.Close()
		delete(svr.hubs, hubID)
		delete(svr.hubHosts, host)
	}

	return closeHub, nil
}

// ReceiveHeader handles Tunnel.ReceiveHeader gRPC call.
func (svr *Server) ReceiveHeader(ctx context.Context, _ *google_protobuf.Empty) (*api.RequestHeader, error) {
	hub, tunnelID, err := svr.verifySessionMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return hub.ReceiveHeader(tunnelID)
}

// ReceiveBody handles Tunnel.ReceiveBody gRPC call.
func (svr *Server) ReceiveBody(_ *google_protobuf.Empty, stream api.Tunnel_ReceiveBodyServer) error {
	hub, tunnelID, err := svr.verifySessionMetadata(stream.Context())
	if err != nil {
		return err
	}

	return hub.ReceiveBody(tunnelID, stream)
}

// SendHeader handles Tunnel.SendHeader gRPC call.
func (svr *Server) SendHeader(ctx context.Context, resHeader *api.ResponseHeader) (*google_protobuf.Empty, error) {
	hub, tunnelID, err := svr.verifySessionMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return &google_protobuf.Empty{}, hub.SendHeader(tunnelID, resHeader)
}

// SendBody handles Tunnel.SendBody gRPC call.
func (svr *Server) SendBody(stream api.Tunnel_SendBodyServer) error {
	hub, tunnelID, err := svr.verifySessionMetadata(stream.Context())
	if err != nil {
		return err
	}

	err = hub.SendBody(tunnelID, stream)
	if err != nil {
		return err
	}
	return stream.SendAndClose(&google_protobuf.Empty{})
}

// SendError handles Tunnel.SendError gRPC call.
func (svr *Server) SendError(ctx context.Context, resErr *api.ResponseError) (*google_protobuf.Empty, error) {
	hub, tunnelID, err := svr.verifySessionMetadata(ctx)
	if err != nil {
		return nil, err
	}

	err = hub.SendError(tunnelID, resErr)
	if err != nil {
		return nil, err
	}
	return &google_protobuf.Empty{}, nil
}

func (svr *Server) verifySessionMetadata(ctx context.Context) (*hub, int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, 0, errors.New("failed to retrieve a metadata")
	}

	hubID, tunnelID, err := svr.decodeSessionID(getMD(md, "kuma-session-id"))
	if err != nil {
		return nil, 0, err
	}

	hub, ok := svr.getHubFromID(hubID)
	if !ok {
		return nil, 0, errors.New("invalid session id")
	}

	peerID := getPeerID(ctx)
	if hub.PeerID != peerID {
		return nil, 0, errors.New("invalid client certificate")
	}

	return hub, tunnelID, nil
}

func getMD(md map[string][]string, key string) string {
	if vs, ok := md[key]; ok && len(vs) > 0 {
		return vs[len(vs)-1]
	}

	return ""
}

func getPeerID(ctx context.Context) string {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return ""
	}

	ti, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return ""
	}

	if len(ti.State.PeerCertificates) < 1 {
		return ""
	}

	return ti.State.PeerCertificates[0].Subject.CommonName
}
