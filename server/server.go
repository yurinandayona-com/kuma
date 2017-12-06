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
	"google.golang.org/grpc/metadata"
)

var (
	validSubdomain = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
)

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

	BaseDomain   string
	HashID       *hashids.HashID
	UserVerifier UserVerifier

	hubs map[string]*hub
}

func (svr *Server) init() {
	if svr.hubs == nil {
		svr.hubs = make(map[string]*hub, 0)
	}
}

func (svr *Server) getHub(host string) (*hub, bool) {
	svr.RLock()
	defer svr.RUnlock()
	svr.init()

	hub, ok := svr.hubs[host]
	return hub, ok
}

// ServeHTTP handles HTTP request.
func (svr *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	if hub, ok := svr.getHub(host); ok {
		hub.ServeHTTP(w, r)
	} else {
		http.Error(w, "kuma: host not found", http.StatusNotFound)
	}
}

// Prepare handles Hub.Prepare gRPC call.
func (svr *Server) Prepare(ctx context.Context, config *api.HubConfig) (*api.HubInfo, error) {
	log.Printf("debug: Hub.Prepare: start: %s", config.Subdomain)
	defer log.Printf("debug: Hub.Prepare: end: %s", config.Subdomain)

	_, err := svr.verifyHubMetadata(ctx)
	if err != nil {
		return nil, err
	}

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
	user, err := svr.verifyHubMetadata(ctx)
	if err != nil {
		return err
	}

	if !strings.HasSuffix(info.Host, "."+svr.BaseDomain) {
		return errors.New("invalid subdomain")
	}
	subdomain := strings.TrimSuffix(info.Host, "."+svr.BaseDomain)
	if !validSubdomain.MatchString(subdomain) {
		return errors.New("invalid subdomain")
	}

	host := subdomain + "." + svr.BaseDomain
	_, closeHub, err := svr.newHub(user, host, stream)
	if err != nil {
		return err
	}
	defer closeHub()

	<-ctx.Done()

	return nil
}

func (svr *Server) newHub(user User, host string, stream api.Hub_ConnectServer) (*hub, func(), error) {
	svr.Lock()
	defer svr.Unlock()
	svr.init()

	if _, ok := svr.hubs[host]; ok {
		return nil, nil, errors.New("failed to create a new hub")
	}

	hub := newHub(host, svr, user, stream)
	svr.hubs[host] = hub

	closeHub := func() {
		svr.Lock()
		defer svr.Unlock()

		hub.Close()
		delete(svr.hubs, host)
	}

	return hub, closeHub, nil
}

func (svr *Server) verifyHubMetadata(ctx context.Context) (User, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("failed to retrieve a metadata")
	}

	user, err := svr.UserVerifier.Verify(mdGet(md, "token"))
	if err != nil {
		return nil, err
	}

	return user, nil
}

// ReceiveHeader handles Tunnel.ReceiveHeader gRPC call.
func (svr *Server) ReceiveHeader(ctx context.Context, _ *google_protobuf.Empty) (*api.RequestHeader, error) {
	hub, tunnelID, err := svr.verifyTunnelMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return hub.ReceiveHeader(tunnelID)
}

// ReceiveBody handles Tunnel.ReceiveBody gRPC call.
func (svr *Server) ReceiveBody(_ *google_protobuf.Empty, stream api.Tunnel_ReceiveBodyServer) error {
	hub, tunnelID, err := svr.verifyTunnelMetadata(stream.Context())
	if err != nil {
		return err
	}

	return hub.ReceiveBody(tunnelID, stream)
}

// SendHeader handles Tunnel.SendHeader gRPC call.
func (svr *Server) SendHeader(ctx context.Context, resHeader *api.ResponseHeader) (*google_protobuf.Empty, error) {
	hub, tunnelID, err := svr.verifyTunnelMetadata(ctx)
	if err != nil {
		return nil, err
	}

	return &google_protobuf.Empty{}, hub.SendHeader(tunnelID, resHeader)
}

// SendBody handles Tunnel.SendBody gRPC call.
func (svr *Server) SendBody(stream api.Tunnel_SendBodyServer) error {
	hub, tunnelID, err := svr.verifyTunnelMetadata(stream.Context())
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
	hub, tunnelID, err := svr.verifyTunnelMetadata(ctx)
	if err != nil {
		return nil, err
	}

	err = hub.SendError(tunnelID, resErr)
	if err != nil {
		return nil, err
	}
	return &google_protobuf.Empty{}, nil
}

func (svr *Server) verifyTunnelMetadata(ctx context.Context) (*hub, int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, 0, errors.New("failed to retrieve a metadata")
	}

	hub, ok := svr.getHub(mdGet(md, "host"))
	if !ok {
		return nil, 0, errors.New("host not found")
	}

	user, err := svr.UserVerifier.Verify(mdGet(md, "token"))
	if err != nil {
		return nil, 0, err
	}
	if user.GetID() != hub.User.GetID() {
		return nil, 0, errors.New("invalid user")
	}

	tunnelID, err := decodeTunnelID(svr.HashID, mdGet(md, "tunnel-id"))
	if err != nil {
		return nil, 0, err
	}

	return hub, tunnelID, nil
}

func mdGet(md map[string][]string, key string) string {
	if vs, ok := md[key]; ok && len(vs) > 0 {
		return vs[len(vs)-1]
	}

	return ""
}
