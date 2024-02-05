package confirm

import (
	"auth-service/internal/app/adapter"
	m "auth-service/internal/app/model"
	vt "auth-service/internal/app/repository/verificationToken"
	e "auth-service/internal/pkg/errors/http"
	"auth-service/internal/pkg/utils"
	"context"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type ConfirmRequest struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}

type ConfirmResponse struct {
	Verified bool `json:"verified"`
}

func validateConfirmRequest(request ConfirmRequest) error {
	if request.Token == "" || request.UserId == "" {
		return e.NewHttpError(400, "", fmt.Errorf("One of the parameters is empty."))
	}

	return nil
}

func CreateAndStoreVerificationToken(userId, userEmail string, tokenRepo vt.Repository) (string, error) {
	token, tokenHash, err := utils.GenerateAndHashToken(12)

	if err != nil {
		return "", e.NewHttpError(http.StatusInternalServerError, err.Error(), fmt.Errorf("Something went wrong."))
	}

	verificationTokenItem := m.VerificationToken{
		UserId: userId,
		SendTo: userEmail,
		Token:  tokenHash,
	}

	if err := tokenRepo.Create(&verificationTokenItem, userEmail); err != nil {
		return "", e.NewHttpErrorByDbStatus(err)
	}

	return token, nil
}

func ConfirmEmail(
	request ConfirmRequest,
	user adapter.UserGrpcInterface,
	verificationToken vt.Repository,
) (*ConfirmResponse, error) {
	if err := validateConfirmRequest(request); err != nil {
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

	verified, httpErr := user.UpdateVerificationStatus(context.Background(), request.UserId, existVerificationToken.SendTo)

	if httpErr != nil {
		return nil, httpErr
	}

	if err := verificationToken.Delete(map[string]interface{}{"id": existVerificationToken.ID}); err != nil {
		return nil, e.NewHttpErrorByDbStatus(err)
	}

	return &ConfirmResponse{Verified: verified}, nil
}
