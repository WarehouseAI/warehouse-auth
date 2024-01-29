package register

import (
	"auth-service/internal/app/adapter"
	"auth-service/internal/app/dataservice"
	m "auth-service/internal/app/model"
	"auth-service/internal/pkg/email"
	e "auth-service/internal/pkg/errors/http"
	gen "auth-service/internal/pkg/protogen"
	"context"
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/mail"
	"os"

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

func generateToken(length int) (string, error) {
	randomBytes := make([]byte, length)

	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	key := base64.URLEncoding.EncodeToString(randomBytes)
	key = key[:length]

	return key, nil
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

	token, err := generateToken(12)

	if err != nil {
		return nil, e.NewHttpError(500, err.Error(), fmt.Errorf("Something went wrong."))
	}

	tokenHash, err := bcrypt.GenerateFromPassword([]byte(token), 12)

	if err != nil {
		return nil, e.NewHttpError(500, err.Error(), fmt.Errorf("Something went wrong."))
	}

	// Store verification token
	verificationTokenItem := m.VerificationToken{
		UserId: userId,
		Token:  string(tokenHash),
	}

	if err := tokenRepository.Create(&verificationTokenItem); err != nil {
		if err := broker.SendMessage("user_saga", userId); err != nil {
			return nil, e.NewHttpError(500, err.Error(), fmt.Errorf("Something went wrong."))
		}

		return nil, e.NewHttpErrorByDbStatus(err)
	}

	emailPayload := email.Messages[email.AccountVerification]
	emailPayload.Message = fmt.Sprintf(emailPayload.Message, req.Firstname, fmt.Sprintf("%s/register/confirm?user=%s&token=%s", os.Getenv("DOMAIN"), userId, token))
	message := email.Email{
		To:           req.Email,
		EmailPayload: emailPayload,
	}

	if err := broker.SendMessage("mail", message); err != nil {
		if err := broker.SendMessage("user_saga", userId); err != nil {
			return nil, e.NewHttpError(500, err.Error(), fmt.Errorf("Something went wrong."))
		}

		return nil, e.NewHttpErrorByDbStatus(err)
	}

	return &RegisterResponse{UserId: userId}, nil
}
