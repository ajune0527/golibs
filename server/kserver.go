package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	netHttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/encoding/json"
	"github.com/go-kratos/kratos/v2/transport/http"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
)

type ServOption func(*KServer)

func AddressK(addr string) ServOption {
	return func(s *KServer) {
		s.address = addr
	}
}

// Timeout with server timeout.
func Timeout(timeout time.Duration) ServOption {
	return func(s *KServer) {
		s.timeout = timeout
	}
}

func UseSnakeCase() ServOption {
	return func(server *KServer) {
		json.MarshalOptions = protojson.MarshalOptions{
			EmitUnpopulated: true,
			UseProtoNames:   true,
		}
	}
}

type KServer struct {
	Server      *http.Server
	address     string
	stopTimeout time.Duration
	timeout     time.Duration
}

func NewKServer(opts ...ServOption) *KServer {
	srv := &KServer{
		stopTimeout: 3 * time.Second,
		timeout:     10 * time.Second,
	}

	for _, o := range opts {
		o(srv)
	}

	op := []http.ServerOption{
		http.Timeout(srv.timeout),
		http.Address(srv.address),
	}
	srv.Server = http.NewServer(op...)

	return srv
}

func (s *KServer) Start(ctx context.Context) error {
	log.Println("[HTTP] Server listen:" + s.address)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.Server.Start(ctx); err != nil && !errors.Is(err, netHttp.ErrServerClosed) {
			return fmt.Errorf("HTTP listen: %s", err)
		}
		return nil
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			_ = s.Stop(ctx)
			return ctx.Err()
		case <-c:
			log.Printf("[HTTP] Shutdown waiting %s\n", s.stopTimeout.String())
			time.Sleep(s.stopTimeout)
			err := s.Stop(ctx)
			if err != nil {
				return fmt.Errorf("HTTP Server Shutdown:%s", err)
			}

			log.Println("[HTTP] server shutdown success")
			return nil
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (s *KServer) Stop(ctx context.Context) error {
	log.Println("[HTTP] server stopping")
	_ = s.Server.Close()
	return s.Server.Shutdown(ctx)
}

func (s *KServer) RegisterGinRouter(ginRouter *gin.Engine) *KServer {
	for _, info := range ginRouter.Routes() {
		s.Server.HandleFunc(info.Path, func(w http.ResponseWriter, req *http.Request) {
			ginRouter.ServeHTTP(w, req)
		})
	}

	return s
}
