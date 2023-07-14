// Package grpc implements and describes the grpc server of the application.
package grpc

import (
	"context"
	"net"
	"os"

	"golang.org/x/exp/slog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sreway/shorturl/internal/config"
	"github.com/sreway/shorturl/internal/usecases"
	pb "github.com/sreway/shorturl/proto/shorturl/v1"
)

type (
	delivery struct {
		shortener usecases.Shortener
		pb.UnimplementedShortURLServiceServer
		logger *slog.Logger
	}
)

// New implements grpc server initialization.
func New(uc usecases.Shortener) (*delivery, error) {
	log := slog.New(slog.NewJSONHandler(os.Stdout).
		WithAttrs([]slog.Attr{slog.String("service", "grpc")}))

	d := &delivery{
		shortener: uc,
		logger:    log,
	}

	return d, nil
}

// Run implements run grpc server.
func (d *delivery) Run(ctx context.Context, config config.GRPC) error {
	var serverOptions []grpc.ServerOption

	if config.UseTLS() {
		tls, err := credentials.NewServerTLSFromFile(config.GetTLS().CertPath, config.GetTLS().KeyPath)
		if err != nil {
			return err
		}
		serverOptions = append(serverOptions, grpc.Creds(tls))
	}

	server := grpc.NewServer(serverOptions...)

	pb.RegisterShortURLServiceServer(server, d)

	ctxServer, stopServer := context.WithCancel(context.Background())
	defer stopServer()

	listen, err := net.Listen("tcp", config.GetAddress())
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		d.logger.Info("trigger graceful shutdown grpc server")
		server.GracefulStop()
		stopServer()
	}()

	d.logger.Info("grpc server is ready to listen and serv")

	err = server.Serve(listen)
	if err != nil {
		return err
	}

	<-ctxServer.Done()
	return nil
}
