package serve

import (
	"github.com/speps/go-hashids"
	"github.com/yurinandayona-com/kuma/api"
	"github.com/yurinandayona-com/kuma/server"
	"github.com/yurinandayona-com/kuma/user_db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"net/http"
)

const (
	hashIDsMinLength = 17
)

type runner struct {
	Config *Config

	httpServer *http.Server
	grpcServer *grpc.Server

	server *server.Server
}

func (r *runner) Run() error {
	log.Printf("info: load user DB: %s", r.Config.UserDB)
	userDB, err := user_db.LoadUserDB(r.Config.UserDB)
	if err != nil {
		return err
	}

	hd := hashids.NewData()
	hd.Salt = r.Config.HashIDsSalt
	hd.MinLength = hashIDsMinLength
	hashID := hashids.NewWithData(hd)

	r.server = &server.Server{
		BaseDomain: r.Config.BaseDomain,
		HashID:     hashID,
		UserVerifier: &user_db.JWTManager{
			UserDB:  userDB,
			HMACKey: []byte(r.Config.HMACKey),
		},
	}

	errCh := make(chan error, 2)

	go r.RunHTTPServer(errCh)
	go r.RunGRPCServer(errCh)

	err = <-errCh
	return err
}

func (r *runner) RunHTTPServer(errCh chan<- error) {
	lis, err := net.Listen("tcp", r.Config.HTTP.Listen)
	if err != nil {
		errCh <- err
		return
	}

	httpServer := &http.Server{
		Handler: r.server,
	}

	r.httpServer = httpServer
	log.Printf("info: start HTTP server: http://%s", lis.Addr())
	errCh <- httpServer.Serve(lis)
	r.Stop()
}

func (r *runner) RunGRPCServer(errCh chan<- error) {
	lis, err := net.Listen("tcp", r.Config.GRPC.Listen)
	if err != nil {
		errCh <- err
		return
	}

	serverOpt := make([]grpc.ServerOption, 0)
	if r.Config.GRPC.UseTLS {
		creds, err := credentials.NewServerTLSFromFile(r.Config.GRPC.TLSCert, r.Config.GRPC.TLSKey)
		if err != nil {
			errCh <- err
			return
		}
		serverOpt = append(serverOpt, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(serverOpt...)
	api.RegisterHubServer(grpcServer, r.server)
	api.RegisterTunnelServer(grpcServer, r.server)

	r.grpcServer = grpcServer
	log.Printf("info: start gRPC server: %s", lis.Addr())
	errCh <- grpcServer.Serve(lis)
	r.Stop()
}

func (r *runner) Stop() {
	if r.httpServer != nil {
		r.httpServer.Close()
	}

	if r.grpcServer != nil {
		r.grpcServer.Stop()
	}
}
