package serve

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
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
		cert, err := tls.LoadX509KeyPair(r.Config.GRPC.TLSCert, r.Config.GRPC.TLSKey)
		if err != nil {
			errCh <- err
			return
		}
		cfg := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		if r.Config.GRPC.TLSClientCA != "" {
			cfg.ClientAuth = tls.RequireAndVerifyClientCert
			clientCAs := x509.NewCertPool()
			pem, err := ioutil.ReadFile(r.Config.GRPC.TLSClientCA)
			if err != nil {
				errCh <- err
				return
			}
			if !clientCAs.AppendCertsFromPEM(pem) {
				errCh <- errors.New("failed to append certificates")
				return
			}
			cfg.ClientCAs = clientCAs
		}
		cfg.BuildNameToCertificate()

		creds := credentials.NewTLS(cfg)
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
