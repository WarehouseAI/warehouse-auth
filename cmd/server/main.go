package server

import (
	"auth-service/internal/app/adapter/broker"
	rt "auth-service/internal/app/repository/resetToken"
	s "auth-service/internal/app/repository/session"
	vt "auth-service/internal/app/repository/verificationToken"
	h "auth-service/internal/app/server/handlers"
	mw "auth-service/internal/app/server/middleware"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func setupCors(origins []string) func(http.Handler) http.Handler {
	allowedHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Content-Length", "Accept"})
	allowedOrigins := handlers.AllowedOrigins(origins)
	allowedMethods := handlers.AllowedMethods([]string{"GET", "DELETE", "POST", "PUT", "DELETE"})
	allowCredentials := handlers.AllowCredentials()

	return handlers.CORS(allowedHeaders, allowedOrigins, allowedMethods, allowCredentials)
}

func setupRoutes(r *mux.Router, logger *logrus.Logger, h *h.Handler) {
	r.HandleFunc("/login", mw.LoggerMW(logger, h.LoginHandler)).Methods("POST")
	r.HandleFunc("/register", mw.LoggerMW(logger, h.RegisterHandler)).Methods("POST")
	r.HandleFunc("/confirm", mw.LoggerMW(logger, h.ConfirmHandler)).Methods("GET")
	r.HandleFunc("/reset/request", mw.LoggerMW(logger, h.RequestPasswordResetHandler)).Methods("POST")
	r.HandleFunc("/reset/verify", mw.LoggerMW(logger, h.VerifyPasswordResetHandler)).Methods("GET")
	r.HandleFunc("/reset/confirm", mw.LoggerMW(logger, h.ConfirmPasswordResetHandler)).Methods("POST")
	r.HandleFunc("/logout", mw.LoggerMW(logger, h.LogoutHandler)).Methods("DELETE")
	r.HandleFunc("/whoami", mw.LoggerMW(logger, h.WhoAmIHandler)).Methods("GET")
}

func Start(
	port string,
	userServiceAddr string,
	allowedOrigins []string,
	resetDB *rt.Database,
	verificationDB *vt.Database,
	sessionDB *s.Database,
	broker *broker.Broker,
	logger *logrus.Logger,
) error {
	hands := h.NewHandlers(resetDB, verificationDB, sessionDB, userServiceAddr, broker)
	router := mux.NewRouter()

	authRouter := router.PathPrefix("/auth").Subrouter()
	setupRoutes(authRouter, logger, hands)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), setupCors(allowedOrigins)(router))
}
