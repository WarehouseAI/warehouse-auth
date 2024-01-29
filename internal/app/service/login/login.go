package login

import (
	"auth-service/internal/app/adapter"
	"auth-service/internal/app/dataservice"
	"auth-service/internal/app/model"
	e "auth-service/internal/pkg/errors/http"
	"context"
	"fmt"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserId string `json:"user_id"`
}

func validateLoginRequest(req *LoginRequest) error {
	if _, err := mail.ParseAddress(req.Email); err != nil {
		return e.NewHttpError(400, err.Error(), fmt.Errorf("Invalid email address"))
	}

	return nil
}

func Login(req *LoginRequest, user adapter.UserGrpcInterface, session dataservice.SessionInterface) (*LoginResponse, *model.Session, error) {
	if err := validateLoginRequest(req); err != nil {
		return nil, nil, err
	}

	existUser, httpErr := user.GetByEmail(context.Background(), req.Email)

	if httpErr != nil {
		return nil, nil, httpErr
	}

	if !existUser.Verified {
		return nil, nil, e.NewHttpError(403, "", fmt.Errorf("Verify your email first"))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(req.Password)); err != nil {
		return nil, nil, e.NewHttpError(400, err.Error(), fmt.Errorf("Invalid credentials"))
	}

	newSession, dbErr := session.Create(context.Background(), existUser.Id)

	if dbErr != nil {
		return nil, nil, e.NewHttpErrorByDbStatus(dbErr)
	}

	return &LoginResponse{UserId: existUser.Id}, newSession, nil
}
