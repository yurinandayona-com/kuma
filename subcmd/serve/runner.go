package serve

import (
	"log"
	"net"
	"net/http"

	"github.com/speps/go-hashids"
	"github.com/yurinandayona-com/kuma/api"
	"github.com/yurinandayona-com/kuma/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	hashIDMinLength = 17
)

type runner struct {
	Config *Config

	httpServer *http.Server
	grpcServer *grpc.Server

	server *server.Server
}

func (r *runner) Run() error {
	hd := hashids.NewData()
	hd.Salt = r.Config.HashIDSecret
	hd.MinLength = hashIDMinLength
	hashID := hashids.NewWithData(hd)

	svr, err := server.New(&server.Config{
		BaseDomain: r.Config.BaseDomain,
		HashID:     hashID,
	})
	if err != nil {
		return err
	}
	r.server = svr

	errCh := make(chan error, 2)

	go r.runHTTPServer(errCh)
	go r.runGRPCServer(errCh)

	err = <-errCh
	return err
}

func (r *runner) runHTTPServer(errCh chan<- error) {
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

func (r *runner) runGRPCServer(errCh chan<- error) {
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
