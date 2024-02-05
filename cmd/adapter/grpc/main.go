package grpc

import (
	a "auth-service/internal/app/adapter"
	"auth-service/internal/app/adapter/grpc/server"
	s "auth-service/internal/app/repository/session"
	vt "auth-service/internal/app/repository/verificationToken"
	"auth-service/internal/pkg/protogen"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start(host string, sessionDB s.Repository, verificationTokenDB vt.Repository, broker a.BrokerInterface, logger *logrus.Logger) func() {
	grpc := grpc.NewServer()
	server := newAuthGrpcServer(sessionDB, verificationTokenDB, broker)
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

func newAuthGrpcServer(sessionDB s.Repository, tokenDB vt.Repository, broker a.BrokerInterface) *server.AuthGrpcServer {
	return &server.AuthGrpcServer{
		SessionRepo: sessionDB,
		TokenRepo:   tokenDB,
		Broker:      broker,
	}
}
