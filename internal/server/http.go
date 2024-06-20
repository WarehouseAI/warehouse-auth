package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/warehouse/auth-service/internal/config"
	internalHttp "github.com/warehouse/auth-service/internal/handler/http"
	"github.com/warehouse/auth-service/internal/handler/middlewares"
	"github.com/warehouse/auth-service/internal/pkg/logger"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type appServer struct {
	log      logger.Logger
	cfg      config.Server
	server   *http.Server
	wg       sync.WaitGroup
	listener net.Listener

	middleware    middlewares.Middleware
	authEndpoints internalHttp.Handler
}

func (a *appServer) Start() {
	a.log.Zap().Info("Start app server", zap.Int("port", a.cfg.Port))

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()

		if err := a.server.Serve(a.listener); err != nil && err != http.ErrServerClosed {
			a.log.Zap().Panic("Error while serve app server", zap.Error(err))
		}
	}()
}

func (a *appServer) Stop() error {
	a.log.Zap().Info("Stop app server")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return err
	}

	a.wg.Wait()
	return nil
}

func NewHttpServer(
	log logger.Logger,
	cfg config.Server,

	middleware middlewares.Middleware,
	authEndpoints internalHttp.Handler,
) (Server, error) {
	var err error
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("cannot listen app port: %w", err)
	}

	router := mux.NewRouter()
	server := &appServer{
		log: log.Named("app_server"),
		cfg: cfg,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Handler: router,
		},
		listener:      listener,
		middleware:    middleware,
		authEndpoints: authEndpoints,
	}
	server.initRoutes(router)
	return server, nil
}

func (s *appServer) initRoutes(router *mux.Router) {
	router.Use(s.middleware.QueueMiddleware)

	s.authEndpoints.FillHandlers(router)
}
