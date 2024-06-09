package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/warehouse/auth-service/internal/config"
	internalGrpc "github.com/warehouse/auth-service/internal/handler/grpc"
	"github.com/warehouse/auth-service/internal/pkg/logger"
	"github.com/warehouse/auth-service/internal/warehousepb"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type grpcServer struct {
	log      logger.Logger
	cfg      config.Grpc
	server   *grpc.Server
	wg       sync.WaitGroup
	listener net.Listener

	authHandler *internalGrpc.AuthHandler
}

func (g *grpcServer) Start() {
	g.log.Zap().Info("Start grpc server", zap.String("host", g.cfg.AuthAddr))

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		if err := g.server.Serve(g.listener); err != nil {
			g.log.Zap().Panic("Error while server grpc server", zap.Error(err))
		}
	}()
}

func (g *grpcServer) Stop() error {
	g.log.Zap().Info("Stop grpc server")

	g.server.Stop()

	g.wg.Wait()
	return nil
}

func NewGrpcServer(
	log logger.Logger,
	cfg config.Config,
	authHandler *internalGrpc.AuthHandler,
) (Server, error) {
	var err error
	listener, err := net.Listen("tcp", cfg.Grpc.UserHost)
	if err != nil {
		return nil, fmt.Errorf("cannot listen grps host: %w", err)
	}
	server := grpc.NewServer()
	warehousepb.RegisterAuthServer(server, authHandler)

	return &grpcServer{
		log:         log,
		cfg:         cfg.Grpc,
		server:      server,
		listener:    listener,
		authHandler: authHandler,
	}, nil
}
