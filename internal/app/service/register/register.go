package register

import (
	"auth-service/internal/app/adapter"
	"auth-service/internal/app/dataservice"
	"auth-service/internal/app/service/confirm"
	"auth-service/internal/pkg/email"
	e "auth-service/internal/pkg/errors/http"
	gen "auth-service/internal/pkg/protogen"
	"context"
	"fmt"
	"net/http"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	// Image     string `json:"image"`
	Email     string `json:"email"`
	ViaGoogle bool   `json:"via_google"`
}

type RegisterResponse struct {
	UserId string `json:"user_id"`
}

func hashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 12)

	return string(hash)
}

func validateRegisterRequest(req *RegisterRequest) error {
	if len(req.Password) > 72 {
		return e.NewHttpError(400, "", fmt.Errorf("Password is too long"))
	}

	if len(req.Password) < 8 {
		return e.NewHttpError(400, "", fmt.Errorf("Password is too short"))
	}

	if _, err := mail.ParseAddress(req.Email); err != nil {
		return e.NewHttpError(400, err.Error(), fmt.Errorf("The provided string is not email"))
	}

	return nil
}

func Register(
	req *RegisterRequest,
	user adapter.UserGrpcInterface,
	tokenRepository dataservice.VerificationTokenInterface,
	broker adapter.BrokerInterface,
) (*RegisterResponse, error) {
	if err := validateRegisterRequest(req); err != nil {
		return nil, err
	}

	// Create user
	userId, httpErr := user.Create(context.Background(), &gen.CreateUserMsg{Firstname: req.Firstname, Lastname: req.Lastname, Username: req.Username, Password: hashPassword(req.Password), Email: req.Email, ViaGoogle: req.ViaGoogle})

	if httpErr != nil {
		return nil, httpErr
	}

	verificationToken, createTokenErr := confirm.CreateAndStoreVerificationToken(userId, req.Email, tokenRepository)

	if createTokenErr != nil {
		if brokerErr := broker.SendMessage("user_saga", userId); brokerErr != nil {
			return nil, e.NewHttpError(500, brokerErr.Error(), fmt.Errorf("Something went wrong."))
		}

		return nil, createTokenErr
	}

	if err := email.SendVerification(userId, verificationToken, req.Firstname, req.Email, broker); err != nil {
		if err := broker.SendMessage("user_saga", userId); err != nil {
			return nil, e.NewHttpError(http.StatusInternalServerError, err.Error(), fmt.Errorf("Something went wrong."))
		}

		return nil, e.NewHttpError(http.StatusInternalServerError, err.Error(), fmt.Errorf("Something went wrong."))
	}

	return &RegisterResponse{UserId: userId}, nil
}
