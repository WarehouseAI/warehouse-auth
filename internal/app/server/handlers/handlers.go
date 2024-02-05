package handlers

import (
	"auth-service/internal/app/adapter/broker"
	"auth-service/internal/app/adapter/grpc/client/user"
	rt "auth-service/internal/app/repository/resetToken"
	s "auth-service/internal/app/repository/session"
	vt "auth-service/internal/app/repository/verificationToken"

	"auth-service/internal/app/server/response"
	"auth-service/internal/app/service"
	"auth-service/internal/app/service/confirm"
	"auth-service/internal/app/service/login"
	"auth-service/internal/app/service/register"
	e "auth-service/internal/pkg/errors/http"
	"encoding/json"
	"fmt"
	"net/http"
)

// func newHttpHandler(
// 	resetTokenDB *tokendata.Database[m.ResetToken],
// 	verificationTokenDB *tokendata.Database[m.VerificationToken],
// 	sessionDB *sessiondata.Database,
// 	pictureStorage *picturedata.Storage,
// 	mailProducer *broker.Broker,
// 	logger *logrus.Logger,
// ) *h.Handler {

// 	userClient := user.NewUserGrpcClient("user:8001")

// 	return &h.Handler{
// 		ResetTokenDB:        resetTokenDB,
// 		VerificationTokenDB: verificationTokenDB,
// 		SessionDB:           sessionDB,
// 		PictureStorage:      pictureStorage,
// 		Broker:              mailProducer,
// 		Logger:              logger,
// 		UserClient:          userClient,
// 	}
// }

type Handler struct {
	ResetRepo        *rt.Database // reset tokens
	VerificationRepo *vt.Database // verification tokens
	SessionRepo      *s.Database  // sessions repo
	UserClient       *user.UserGrpcClient
	Broker           *broker.Broker
}

// TODO: add error & logger middleware

func NewHandlers(
	resetRepo *rt.Database,
	verificationRepo *vt.Database,
	sessionRepo *s.Database,
	userAddr string,
	broker *broker.Broker,
) *Handler {
	userClient := user.NewUserGrpcClient(userAddr)

	return &Handler{
		ResetRepo:        resetRepo,
		VerificationRepo: verificationRepo,
		SessionRepo:      sessionRepo,
		UserClient:       userClient,
		Broker:           broker,
	}
}

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req register.RegisterRequest
	err := r.ParseMultipartForm(5 << 20)

	if err != nil {
		response.WithError(w, e.NewHttpError(http.StatusInternalServerError, err.Error(), fmt.Errorf("Something went wrong.")))
		return
	}

	// TODO: Добавить отправку изображения в user сервис через grpc
	form := r.MultipartForm

	req.Username = form.Value["username"][0]
	req.Firstname = form.Value["firstname"][0]
	req.Lastname = form.Value["lastname"][0]
	req.Password = form.Value["password"][0]
	req.Email = form.Value["email"][0]
	req.ViaGoogle = false

	regResponse, regErr := register.Register(&req, h.UserClient, h.VerificationRepo, h.Broker)

	if regErr != nil {
		response.WithError(w, regErr)
		return
	}

	response.JSON(w, http.StatusCreated, regResponse)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var request login.LoginRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&request); err != nil {
		response.WithError(w, e.NewHttpError(http.StatusInternalServerError, err.Error(), fmt.Errorf("Something went wrong.")))
		return
	}

	loginResp, session, err := login.Login(&request, h.UserClient, h.SessionRepo)

	if err != nil {
		response.WithError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    session.ID,
		Path:     "/",
		MaxAge:   0,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	response.JSON(w, http.StatusOK, loginResp)
}

func (h *Handler) ConfirmHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	user := r.URL.Query().Get("user")

	request := confirm.ConfirmRequest{
		UserId: user,
		Token:  token,
	}

	confirmResp, err := confirm.ConfirmEmail(request, h.UserClient, h.VerificationRepo)

	if err != nil {
		response.WithError(w, err)
		return
	}

	response.JSON(w, http.StatusOK, confirmResp)
}

func (h *Handler) ConfirmPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	resetTokenId := r.URL.Query().Get("token_id")
	var request service.ResetConfirmRequest

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&request); err != nil {
		response.WithError(w, e.NewHttpError(500, err.Error(), fmt.Errorf("Something went wrong.")))
		return
	}

	resetResp, resetErr := service.ConfirmResetToken(&request, resetTokenId, h.UserClient, h.ResetRepo)

	if resetErr != nil {
		response.WithError(w, resetErr)
		return
	}

	response.JSON(w, http.StatusOK, resetResp)
}

func (h *Handler) VerifyPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	resetCode := r.URL.Query().Get("verification")
	resetTokenId := r.URL.Query().Get("token_id")

	resetResp, resetErr := service.VerifyResetCode(resetCode, resetTokenId, h.ResetRepo)

	if resetErr != nil {
		response.WithError(w, resetErr)
		return
	}

	response.JSON(w, http.StatusOK, resetResp)
}

func (h *Handler) RequestPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var request service.ResetAttemptRequest
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&request); err != nil {
		response.WithError(w, e.NewHttpError(http.StatusInternalServerError, err.Error(), fmt.Errorf("Something went wrong.")))
		return
	}

	resetToken, resetErr := service.RequestResetToken(request, h.ResetRepo, h.UserClient, h.Broker)

	if resetErr != nil {
		response.WithError(w, resetErr)
		return
	}

	response.JSON(w, http.StatusCreated, resetToken)
}

func (h *Handler) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("sessionId")

	if err != nil {
		response.WithError(w, e.NewHttpError(
			http.StatusUnauthorized,
			err.Error(),
			fmt.Errorf("Session cookie not found."),
		))
		return
	}

	if err := service.Logout(session.Value, h.SessionRepo); err != nil {
		response.WithError(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    "",
		MaxAge:   -1,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	response.Status(w, http.StatusOK)
}

func (h *Handler) WhoAmIHandler(w http.ResponseWriter, r *http.Request) {
	session, cookieErr := r.Cookie("sessionId")

	if cookieErr != nil {
		response.WithError(w, e.NewHttpError(
			http.StatusUnauthorized,
			cookieErr.Error(),
			fmt.Errorf("Session cookie not found."),
		))
		return
	}

	_, newSession, authErr := service.Authenticate(session.Value, h.SessionRepo)

	if authErr != nil {
		response.WithError(w, authErr)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "sessionId",
		Value:    newSession.ID,
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	})

	response.Status(w, http.StatusOK)
}
