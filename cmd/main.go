package main

import (
	"auth-service/cmd/adapter/broker"
	"auth-service/cmd/adapter/grpc"
	"auth-service/cmd/dataservice"
	"auth-service/cmd/server"
	"auth-service/configs"
	"auth-service/internal/pkg/logger"

	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	config, err := configs.ReadConfig()

	if err != nil {
		log.Fatalln(err)
	}
	log.Infoln("✅ Environment successfully connected.")

	log = logger.Setup(config.Server.Env, log)
	log.Infoln("✅ Logger successfully connected.")

	sessionDB := dataservice.NewSessionDatabase(*config)
	resetTokenDB := dataservice.NewResetTokenDatabase(*config)
	verificationTokenDB := dataservice.NewVerificationTokenDatabase(*config)
	broker := broker.NewBroker(*config)

	log.Infoln("✅ Database successfully connected.")

	grpcServer := grpc.Start("auth:8041", sessionDB, log)
	go grpcServer()

	if err := server.Start(
		config.Server.Port,
		config.Server.UserAddr,
		config.Server.AllowedOrigins,
		resetTokenDB,
		verificationTokenDB,
		sessionDB,
		broker,
		log,
	); err != nil {
		log.Fatalln("❌ Failed to start the HTTP Handler.", err)
	}

	defer func() {
		broker.Channel.Close()
		broker.Connection.Close()
	}()
}
