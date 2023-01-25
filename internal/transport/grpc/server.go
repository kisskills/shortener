package grpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"
	"shortener/internal/entities"
	"shortener/internal/proto_gen/shortener"
)

type Server struct {
	httpPort   int
	grpcPort   int
	log        *zap.SugaredLogger
	svc        shortener.ShortenerServer
	grpsServer *grpc.Server
	httpServer *http.Server
}

func NewServer(log *zap.SugaredLogger, svc shortener.ShortenerServer, httpPort int, grpcPort int) (*Server, error) {
	if log == nil {
		return nil, errors.WithMessage(entities.ErrInvalidParam, "empty logger")
	}

	if svc == nil {
		return nil, errors.WithMessage(entities.ErrInvalidParam, "empty service")
	}

	if httpPort == 0 {
		return nil, errors.WithMessage(entities.ErrInvalidParam, "empty http port")
	}

	if grpcPort == 0 {
		return nil, errors.WithMessage(entities.ErrInvalidParam, "empty grpc port")
	}

	return &Server{
		httpPort:   httpPort,
		grpcPort:   grpcPort,
		log:        log,
		svc:        svc,
		grpsServer: grpc.NewServer(),
	}, nil
}

func (s *Server) Run(ctx context.Context) {
	//позволяет не прикреплять прото файл в постмане
	reflection.Register(s.grpsServer)

	grpcAddr := fmt.Sprintf(":%d", s.grpcPort)
	httpAddr := fmt.Sprintf(":%d", s.httpPort)

	go s.stopProcess(ctx)

	// Setup gRPC servers
	shortener.RegisterShortenerServer(s.grpsServer, s.svc)

	// Setup gRPC gateway.
	rmux := runtime.NewServeMux()
	mux := http.NewServeMux()
	mux.Handle("/", rmux)
	err := shortener.RegisterShortenerHandlerServer(ctx, rmux, s.svc)
	if err != nil {
		s.log.Fatal(err)
	}

	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		s.log.Fatal(err)
	}
	s.log.Infof("server listen grpc %s", grpcAddr)

	go func() {
		if err := s.grpsServer.Serve(grpcListener); err != nil {
			s.log.Fatal(err)
		}
	}()

	s.httpServer = &http.Server{Addr: httpAddr, Handler: mux}
	s.log.Infof("server listen http %s", httpAddr)
	err = s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		s.log.Fatal(err)
	}
}

func (s *Server) stopProcess(ctx context.Context) {
	<-ctx.Done()

	s.grpsServer.GracefulStop()
	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.log.Error(err)
	}
}
