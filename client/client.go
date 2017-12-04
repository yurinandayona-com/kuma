// Package client provides kuma client implementation.
package client

import (
	"crypto/x509"
	"io"
	"log"

	"github.com/pkg/errors"
	"github.com/yurinandayona-com/kuma/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// client is internal type of Client.
//
// It is only for hiding internal fields: BaseCtx and Conn.
type client struct {
	BaseCtx context.Context
	cancel  func()

	Conn *grpc.ClientConn
}

// Client is kuma client.
type Client struct {
	GRPCServer string
	UseTLS     bool
	Token      string

	Port      int
	Subdomain string

	client
}

func (cli *Client) Start() error {
	var err error

	log.Print("debug: connect to gRPC server")
	cli.Conn, err = cli.dialConn()
	if err != nil {
		return err
	}
	defer cli.Conn.Close()

	hub := api.NewHubClient(cli.Conn)
	md := metadata.Pairs(
		"token", cli.Token,
	)
	cli.BaseCtx, cli.cancel = context.WithCancel(context.Background())
	ctx := metadata.NewOutgoingContext(cli.BaseCtx, md)

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

		go cli.handleRequest(req)
	}

	log.Print("info: finish kuma connection")

	return nil
}

func (cli *Client) dialConn() (*grpc.ClientConn, error) {
	opt := make([]grpc.DialOption, 0, 1)
	if cli.UseTLS {
		cert, err := x509.SystemCertPool()
		if err != nil {
			return nil, errors.Wrap(err, "kuma: failed to get system cert pool")
		}
		creds := credentials.NewClientTLSFromCert(cert, "")

		opt = append(opt, grpc.WithTransportCredentials(creds))
	} else {
		opt = append(opt, grpc.WithInsecure())
	}

	conn, err := grpc.Dial(cli.GRPCServer, opt...)
	if err != nil {
		return nil, errors.Wrap(err, "kuma: failed to dial gRPC server")
	}

	return conn, nil
}

func (cli *Client) handleRequest(req *api.Request) {
	t := &tunnel{
		Client:  cli,
		Request: req,
	}

	if err := t.Run(); err != nil {
		log.Printf("error: error on request handling: %s", err)
	}
}
