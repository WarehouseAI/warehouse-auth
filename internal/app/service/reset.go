package service

import (
	"auth-service/internal/app/adapter"
	"auth-service/internal/app/dataservice"
	"auth-service/internal/app/model"
	"auth-service/internal/pkg/email"
	e "auth-service/internal/pkg/errors/http"
	gen "auth-service/internal/pkg/protogen"
	"context"
	"fmt"
	"math/rand"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ResetAttemptRequest struct {
	Email string `json:"email"`
}

type ResetAttemptResponse struct {
	TokenId string `json:"token_id"`
}

type ResetConfirmRequest struct {
	UserId   string `json:"user_id"`
	Password string `json:"password"`
}

type ResetConfirmResponse struct {
	UserId string `json:"user_id"`
}

type ResetVerifyResponse struct {
	UserId string `json:"user_id"`
}

func ConfirmResetToken(request *ResetConfirmRequest, resetTokenId string, user adapter.UserGrpcInterface, resetToken dataservice.ResetTokenInterface) (*ResetConfirmResponse, error) {
	existResetToken, dbErr := resetToken.Get(map[string]interface{}{"id": resetTokenId})

	if dbErr != nil {
		return nil, e.NewHttpErrorByDbStatus(dbErr)
	}

	if err := resetToken.Delete(map[string]interface{}{"id": existResetToken.ID}); err != nil {
		return nil, e.NewHttpErrorByDbStatus(err)
	}

	hash, hashErr := bcrypt.GenerateFromPassword([]byte(request.Password), 12)

	if hashErr != nil {
		return nil, e.NewHttpError(500, hashErr.Error(), fmt.Errorf("Something went wrong."))
	}

	resp, httpErr := user.ResetPassword(context.Background(), &gen.ResetPasswordRequest{UserId: request.UserId, Password: string(hash)})

	if httpErr != nil {
		return nil, httpErr
	}

	return &ResetConfirmResponse{UserId: resp}, nil
}

func VerifyResetCode(verificationCode string, resetTokenId string, resetToken dataservice.ResetTokenInterface) (*ResetVerifyResponse, error) {
	existResetToken, err := resetToken.Get(map[string]interface{}{"id": resetTokenId})

	if err != nil {
		return nil, e.NewHttpErrorByDbStatus(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existResetToken.Token), []byte(verificationCode)); err != nil {
		return nil, e.NewHttpError(400, err.Error(), fmt.Errorf("Invalid reset token"))
	}

	return &ResetVerifyResponse{UserId: existResetToken.UserId.String()}, nil
}

func RequestResetToken(req ResetAttemptRequest, resetToken dataservice.ResetTokenInterface, user adapter.UserGrpcInterface, broker adapter.BrokerInterface) (*ResetAttemptResponse, error) {
	existUser, httpErr := user.GetByEmail(context.Background(), req.Email)

	if httpErr != nil {
		return nil, httpErr
	}

	resetCode := generateCode(8)
	hash, bcryptErr := bcrypt.GenerateFromPassword([]byte(resetCode), 12)

	if bcryptErr != nil {
		return nil, e.NewHttpError(500, bcryptErr.Error(), fmt.Errorf("Something went wrong."))
	}

	newResetToken := &model.ResetToken{
		UserId: uuid.FromStringOrNil(existUser.Id),
		Token:  string(hash),
	}

	if err := resetToken.Create(newResetToken); err != nil {
		return nil, e.NewHttpErrorByDbStatus(err)
	}

	emailPayload := email.Messages[email.AccountReset]
	emailPayload.Message = fmt.Sprintf(emailPayload.Message, existUser.Firstname, existUser.Email, resetCode)
	message := email.Email{
		To:           req.Email,
		EmailPayload: emailPayload,
	}

	if err := broker.SendMessage("mail", message); err != nil {
		return nil, e.NewHttpError(500, err.Error(), fmt.Errorf("Something went wrong."))
	}

	return &ResetAttemptResponse{TokenId: newResetToken.ID.String()}, nil
}

func generateCode(length int) string {
	charset := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	batch := make([]byte, length)

	for i := range batch {
		batch[i] = charset[rand.Intn(len(charset))]
	}

	return string(batch)
}
