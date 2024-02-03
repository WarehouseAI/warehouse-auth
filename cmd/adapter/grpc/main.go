package grpc

import (
	a "auth-service/internal/app/adapter"
	"auth-service/internal/app/adapter/grpc/server"
	d "auth-service/internal/app/dataservice"
	"auth-service/internal/pkg/protogen"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func Start(host string, sessionDB d.SessionInterface, tokenDB d.VerificationTokenInterface, broker a.BrokerInterface, logger *logrus.Logger) func() {
	grpc := grpc.NewServer()
	server := newAuthGrpcServer(sessionDB, tokenDB, broker)
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

func newAuthGrpcServer(sessionDB d.SessionInterface, tokenDB d.VerificationTokenInterface, broker a.BrokerInterface) *server.AuthGrpcServer {
	return &server.AuthGrpcServer{
		SessionRepo: sessionDB,
		TokenRepo:   tokenDB,
		Broker:      broker,
	}
}
