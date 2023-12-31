package servers

import (
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"net"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpcctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

func (s *server) NewGrpcServer() (grpcServer *grpc.Server, err error) {
	tcpListener, err := net.Listen("tcp", s.cfg.Grpc.PORT)
	if err != nil {
		return nil, errors.Wrap(err, "net.Listen")
	}
	//defer tcpListener.Close()

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		s.logger.Fatal("failed to load key pair: %s", err, "app")
	}

	grpcServer = grpc.NewServer(
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpcctxtags.UnaryServerInterceptor(),
			grpcrecovery.UnaryServerInterceptor(),
		)),
	)

	if s.cfg.Server.APP_ENV == "development" {
		reflection.Register(grpcServer)
	}

	go func() {
		s.logger.Info(fmt.Sprintf("%s gRPC server is listening on port: {%s}", s.cfg.Server.PROJECT_NAME, s.cfg.Grpc.PORT), nil, "app")
		s.logger.Error("gRPC server error", grpcServer.Serve(tcpListener), "app")
	}()

	return
}
