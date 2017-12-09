package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	google_protobuf "github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/yurinandayona-com/kuma/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

var (
	httpClient = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
)

type tunnel struct {
	Client    *Client
	SessionID string

	tunnel api.TunnelClient
	ctx    context.Context
}

func (t *tunnel) Run(ctx context.Context) error {
	start := time.Now()

	t.tunnel = api.NewTunnelClient(t.Client.Conn)

	md := metadata.Pairs(
		"kuma-session-id", t.SessionID,
	)
	t.ctx = metadata.NewOutgoingContext(ctx, md)

	ss, err := t.runInternal(ctx)
	ss.Log(time.Since(start))
	if err != nil {
		t.sendError(err)
		return err
	}

	return nil
}

func (t *tunnel) runInternal(ctx context.Context) (ss *stats, err error) {
	ss = &stats{SessionID: t.SessionID}

	httpRequest, w, err := t.receiveHeader(ctx, ss)
	if err != nil {
		return
	}

	errCh := make(chan error, 1)
	go t.receiveBody(ss, errCh, w)

	httpResponse, err := httpClient.Do(httpRequest)
	if err != nil {
		err = errors.Wrap(err, "kuma: failed to do HTTP reuqest")
		return
	}
	defer httpResponse.Body.Close()

	select {
	case err = <-errCh:
		return
	default:
		err = t.sendHeader(ss, httpResponse)
		if err != nil {
			return
		}
	}

	select {
	case err = <-errCh:
		return
	default:
		err = t.sendBody(ss, httpResponse)
		if err != nil {
			return
		}
	}

	close(errCh)

	err = <-errCh
	return
}

func (t *tunnel) receiveHeader(ctx context.Context, ss *stats) (*http.Request, io.WriteCloser, error) {
	reqHeader, err := t.tunnel.ReceiveHeader(t.ctx, &google_protobuf.Empty{})
	if err != nil {
		return nil, nil, errors.Wrap(err, "kuma: failed to receive header")
	}

	r, w := io.Pipe()
	url := fmt.Sprintf("http://localhost:%d%s", t.Client.Port, reqHeader.Path)
	httpReq, err := http.NewRequest(reqHeader.Method, url, r)
	if err != nil {
		return nil, nil, errors.Wrap(err, "kuma: failed to create HTTP request")
	}

	httpReq = httpReq.WithContext(ctx)

	for _, h := range reqHeader.Headers {
		httpReq.Header[h.Name] = h.Values
	}

	ss.Method = reqHeader.Method
	ss.Path = reqHeader.Path

	return httpReq, w, nil
}

func (t *tunnel) receiveBody(ss *stats, errCh chan error, w io.WriteCloser) {
	defer w.Close()

	bodyStream, err := t.tunnel.ReceiveBody(t.ctx, &google_protobuf.Empty{})
	if err != nil {
		errCh <- errors.Wrap(err, "kuma: failed to receive body")
	}

	for {
		body, err := bodyStream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			errCh <- errors.Wrap(err, "kuma: failed to receive body from stream")
			return
		}

		n, err := w.Write(body.Body)
		ss.RequestBodySize += int64(n)
		if err != nil {
			errCh <- errors.Wrap(err, "kuma: failed to write received body")
			return
		}
	}
}

func (t *tunnel) sendHeader(ss *stats, httpRes *http.Response) error {
	headers := make([]*api.Header, 0, len(httpRes.Header))
	for name, values := range httpRes.Header {
		headers = append(headers, &api.Header{Name: name, Values: values})
	}
	resHeader := &api.ResponseHeader{
		Status:  int32(httpRes.StatusCode),
		Headers: headers,
	}

	_, err := t.tunnel.SendHeader(t.ctx, resHeader)
	if err != nil {
		return errors.Wrap(err, "kuma: failed to send header")
	}

	ss.Status = httpRes.StatusCode

	return err
}

func (t *tunnel) sendBody(ss *stats, httpRes *http.Response) error {
	resBody, err := t.tunnel.SendBody(t.ctx)
	if err != nil {
		return errors.Wrap(err, "kuma: failed to send body")
	}

	ss.ResponseBodySize, err = io.Copy(&sendBodyWriter{stream: resBody}, httpRes.Body)
	if err != nil {
		return errors.Wrap(err, "kuma: failed to send body to response stream")
	}

	_, err = resBody.CloseAndRecv()
	if err != nil {
		return errors.Wrap(err, "kuma: failed to close response stream")
	}

	return nil
}

type sendBodyWriter struct {
	stream api.Tunnel_SendBodyClient
}

func (sbw *sendBodyWriter) Write(buf []byte) (int, error) {
	err := sbw.stream.Send(&api.ResponseBody{Body: buf})
	if err != nil {
		return 0, err
	}

	return len(buf), nil
}

func (t *tunnel) sendError(err error) {
	_, err = t.tunnel.SendError(t.ctx, &api.ResponseError{Error: err.Error()})
	if err != nil {
		log.Printf("alert: tunnel [%s]: kuma: failed to send an error to tunnel: %s", t.SessionID, err)
	}
}
