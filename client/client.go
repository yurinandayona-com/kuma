// Package client provides kuma client implementation.
package client

import (
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"log"

	"github.com/pkg/errors"
	"github.com/yurinandayona-com/kuma/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// client is internal type of Client.
//
// It is only for hiding internal field Conn.
type client struct {
	Conn *grpc.ClientConn
}

// Client is kuma client.
type Client struct {
	GRPCServer string

	UseTLS    bool
	TLSRootCA string
	TLSCert   string
	TLSKey    string

	Port      int
	Subdomain string

	client
}

// Run starts client process on given context.
func (cli *Client) Run(ctx context.Context) error {
	var err error

	log.Print("debug: connect to gRPC server")
	cli.Conn, err = cli.dialConn(ctx)
	if err != nil {
		return err
	}
	defer cli.Conn.Close()

	hub := api.NewHubClient(cli.Conn)

	log.Print("debug: prepare hub information")
	info, err := hub.Prepare(ctx, &api.HubConfig{
		Subdomain: cli.Subdomain,
	})
	if err != nil {
		return errors.Wrap(err, "kuma: failed to prepare hub information")
	}

	log.Printf("debug: hub information: host = %#v", info.Host)

	log.Print("debug: connect to hub")
	reqStream, err := hub.Connect(ctx, info)
	if err != nil {
		return errors.Wrap(err, "kuma: failed to connect to hub")
	}

	log.Print("info: start kuma connection")
	log.Printf("info: http://%s is now available", info.Host)
	for {
		req, err := reqStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return errors.Wrap(err, "kuma: failed to receive request from stream")
		}

		go cli.handleRequest(ctx, req)
	}

	log.Print("info: finish kuma connection")

	return nil
}

func (cli *Client) dialConn(ctx context.Context) (*grpc.ClientConn, error) {
	opt := make([]grpc.DialOption, 0, 1)
	if cli.UseTLS {
		cfg := &tls.Config{}

		if cli.TLSRootCA != "" {
			rootCAs := x509.NewCertPool()
			pem, err := ioutil.ReadFile(cli.TLSRootCA)
			if err != nil {
				return nil, errors.Wrap(err, "failed to read root CA file")
			}
			if !rootCAs.AppendCertsFromPEM(pem) {
				return nil, errors.New("failed to append certificates")
			}
			cfg.RootCAs = rootCAs
		}

		if cli.TLSCert != "" && cli.TLSKey != "" {
			cert, err := tls.LoadX509KeyPair(cli.TLSCert, cli.TLSKey)
			if err != nil {
				return nil, errors.Wrap(err, "failed to load x509 key pair")
			}
			cfg.Certificates = []tls.Certificate{cert}
		}
		cfg.BuildNameToCertificate()

		creds := credentials.NewTLS(cfg)
		opt = append(opt, grpc.WithTransportCredentials(creds))
	} else {
		opt = append(opt, grpc.WithInsecure())
	}

	conn, err := grpc.DialContext(ctx, cli.GRPCServer, opt...)
	if err != nil {
		return nil, errors.Wrap(err, "kuma: failed to dial gRPC server")
	}

	return conn, nil
}

func (cli *Client) handleRequest(ctx context.Context, req *api.Request) {
	t := &tunnel{
		Client:    cli,
		SessionID: req.SessionID,
	}

	if err := t.Run(ctx); err != nil {
		log.Printf("error: error on request handling: %s", err)
	}
}
