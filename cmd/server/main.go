package server

import (
	"auth-service/internal/app/adapter/broker"
	d "auth-service/internal/app/dataservice/operations"
	m "auth-service/internal/app/model"
	h "auth-service/internal/app/server/handlers"
	mw "auth-service/internal/app/server/middleware"
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
	r.HandleFunc("/register", mw.WriteLog(logger, h.RegisterHandler)).Methods("POST")
	r.HandleFunc("/register/confirm", mw.WriteLog(logger, h.RegisterVerifyHandler)).Methods("GET")
	r.HandleFunc("/reset/request", mw.WriteLog(logger, h.RequestPasswordResetHandler)).Methods("POST")
	r.HandleFunc("/reset/verify", mw.WriteLog(logger, h.VerifyPasswordResetHandler)).Methods("GET")
	r.HandleFunc("/reset/confirm", mw.WriteLog(logger, h.ConfirmPasswordResetHandler)).Methods("POST")
	r.HandleFunc("/logout", mw.WriteLog(logger, h.LogoutHandler)).Methods("DELETE")
	r.HandleFunc("/whoami", mw.WriteLog(logger, h.WhoAmIHandler)).Methods("GET")
}

func Start(
	port string,
	userServiceAddr string,
	allowedOrigins []string,
	resetDB *d.TokenDatabase[m.ResetToken],
	verificationDB *d.TokenDatabase[m.VerificationToken],
	sessionDB *d.SessionDatabase,
	broker *broker.Broker,
	logger *logrus.Logger,
) error {
	hands := h.NewHandlers(resetDB, verificationDB, sessionDB, userServiceAddr, broker, logger)
	router := mux.NewRouter()

	router.PathPrefix("/auth")
	setupRoutes(router, logger, hands)

	return http.ListenAndServe(":"+port, setupCors(allowedOrigins)(router))
}
