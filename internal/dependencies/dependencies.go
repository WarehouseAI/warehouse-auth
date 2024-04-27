package dependencies

import (
	mailAdpt "auth-service/internal/adapter/mail"
	randomAdpt "auth-service/internal/adapter/random"
	timeAdpt "auth-service/internal/adapter/time"
	userAdpt "auth-service/internal/adapter/user"
	"auth-service/internal/broker"
	"auth-service/internal/config"
	"auth-service/internal/db"
	"auth-service/internal/handler/http"
	"auth-service/internal/handler/middlewares"
	"auth-service/internal/pkg/logger"
	jwtRepo "auth-service/internal/repository/operations/jwt"
	"auth-service/internal/repository/operations/reset_token"
	transactionsRepo "auth-service/internal/repository/operations/transactions"
	"auth-service/internal/repository/operations/verification_token"
	"auth-service/internal/server"
	authSvc "auth-service/internal/service/auth"
	jwtSvc "auth-service/internal/service/jwt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Dependencies interface {
		Close()
		Cfg() *config.Config
		Internal() dependencies
		WaitForInterrupr()

		AppServer() server.Server
	}

	dependencies struct {
		cfg                     *config.Config
		log                     logger.Logger
		warehouseRequestHandler http.WarehouseRequestHandler
		handlerMiddleware       middlewares.Middleware

		psqlClient   *db.PostgresClient
		rabbitClient *broker.RabbitClient

		authHandler http.Handler

		authService authSvc.Service
		jwtService  jwtSvc.Service

		pgxTransactionRepo    transactionsRepo.Repository
		jwtRepo               jwtRepo.Repository
		verificationTokenRepo verification_token.Repository
		resetTokenRepo        reset_token.Repository

		timeAdapter   timeAdpt.Adapter
		randomAdapter randomAdpt.Adapter
		userAdapter   userAdpt.Adapter
		mailAdapter   mailAdpt.Adapter

		appServer server.Server

		shutdownChannel chan os.Signal
		closeCallbacks  []func()
	}
)

func NewDependencies(cfgPath string) (Dependencies, error) {
	cfg, err := config.NewConfig(cfgPath)
	if err != nil && err.Error() == "Config File \"config\" Not Found in \"[]\"" {
		cfg, err = config.NewConfig("./configs/local")
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.LevelKey = "1"
	encoderCfg.TimeKey = "t"

	z := zap.New(
		&logger.WarehouseZapCore{
			Core: zapcore.NewCore(
				zapcore.NewJSONEncoder(encoderCfg),
				zapcore.Lock(os.Stdout),
				zap.NewAtomicLevel(),
			),
		},
		zap.AddCaller(),
	)

	return &dependencies{
		cfg:             cfg,
		log:             logger.NewLogger(z),
		shutdownChannel: make(chan os.Signal),
	}, nil
}

func (d *dependencies) Close() {
	for i := len(d.closeCallbacks) - 1; i >= 0; i-- {
		d.closeCallbacks[i]()
	}
	d.log.Zap().Sync()
}

func (d *dependencies) Internal() dependencies {
	return *d
}

func (d *dependencies) Cfg() *config.Config {
	return d.cfg
}

func (d *dependencies) WarehouseJsonRequestHandler() http.WarehouseRequestHandler {
	if d.warehouseRequestHandler == nil {
		d.warehouseRequestHandler = http.NewWarehouseJsonRequestHandler(d.log, d.cfg.Timeouts.AccCookie)
	}

	return d.warehouseRequestHandler
}

func (d *dependencies) AppServer() server.Server {
	if d.appServer == nil {
		var err error
		msg := "initialize app server"
		if d.appServer, err = server.NewAppServer(
			d.log,
			d.cfg.Server,
			d.HandlerMiddleware(),
			d.AuthHandler(),
		); err != nil {
			d.log.Zap().Panic(msg, zap.Error(err))
		}

		d.closeCallbacks = append(d.closeCallbacks, func() {
			msg := "shutting down app server"
			if err := d.appServer.Stop(); err != nil {
				d.log.Zap().Warn(msg, zap.Error(err))
				return
			}
			d.log.Zap().Info(msg)
		})
	}
	return d.appServer
}

func (d *dependencies) WaitForInterrupr() {
	signal.Notify(d.shutdownChannel, syscall.SIGINT, syscall.SIGTERM)
	d.log.Zap().Info("Wait for receive interrupt signal")
	<-d.shutdownChannel // ждем когда сигнал запишется в канал и сразу убираем его, значит, что сигнал получен
	d.log.Zap().Info("Receive interrupt signal")
}
