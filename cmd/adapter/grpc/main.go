package grpc

import (
	"auth-service/internal/app/adapter/grpc/server"
	"auth-service/internal/app/dataservice"
	"auth-service/internal/pkg/protogen"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start(host string, db dataservice.SessionInterface, logger *logrus.Logger) func() {
	grpc := grpc.NewServer()
	server := newAuthGrpcServer(db, logger)
	listener, err := net.Listen("tcp", host)

	if err != nil {
		fmt.Println("❌Failed to listen the GRPC host.")
		logger.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("Auth Microservice")
		panic(err)
	}

	return func() {
		protogen.RegisterAuthServiceServer(grpc, server)

		if err := grpc.Serve(listener); err != nil {
			fmt.Println("❌Failed to start the GRPC server.")
			logger.WithFields(logrus.Fields{"time": time.Now().String(), "error": err.Error()}).Info("Auth Microservice")
			panic(err)
		}
	}
}

func newAuthGrpcServer(database dataservice.SessionInterface, logger *logrus.Logger) *server.AuthGrpcServer {
	return &server.AuthGrpcServer{
		DB:     database,
		Logger: logger,
	}
}
