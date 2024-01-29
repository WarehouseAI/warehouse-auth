package register

import (
	"auth-service/internal/app/adapter"
	"auth-service/internal/app/dataservice"
	e "auth-service/internal/pkg/errors/http"
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type RegisterVerifyRequest struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}

type RegisterVerifyResponse struct {
	Verified bool `json:"verified"`
}

func validateVerifyRequest(request RegisterVerifyRequest) error {
	if request.Token == "" || request.UserId == "" {
		return e.NewHttpError(400, "", fmt.Errorf("One of the parameters is empty."))
	}

	return nil
}

func RegisterVerify(
	request RegisterVerifyRequest,
	user adapter.UserGrpcInterface,
	verificationToken dataservice.VerificationTokenInterface,
) (*RegisterVerifyResponse, error) {
	if err := validateVerifyRequest(request); err != nil {
		return nil, err
	}

	existVerificationToken, dbErr := verificationToken.Get(map[string]interface{}{"user_id": request.UserId})

	if dbErr != nil {
		return nil, e.NewHttpErrorByDbStatus(dbErr)
	}

	// TODO: Переделать эту логику
	// Удаляем токен, если он протух. Пользователю нужно отправлять запрос еще раз.
	if time.Now().After(existVerificationToken.ExpiresAt) {
		verificationToken.Delete(map[string]interface{}{"id": existVerificationToken.ID})
		return nil, e.NewHttpError(400, "", fmt.Errorf("Invalid or expired verification token."))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existVerificationToken.Token), []byte(request.Token)); err != nil {
		return nil, e.NewHttpError(400, err.Error(), fmt.Errorf("Invalid register verification key"))
	}

	verified, httpErr := user.UpdateVerificationStatus(context.Background(), request.UserId)

	if httpErr != nil {
		return nil, httpErr
	}

	if err := verificationToken.Delete(map[string]interface{}{"id": existVerificationToken.ID}); err != nil {
		return nil, e.NewHttpErrorByDbStatus(err)
	}

	return &RegisterVerifyResponse{Verified: verified}, nil
}
