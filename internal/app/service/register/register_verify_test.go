package register

import (
	m "auth-service/internal/app/model"
	e "auth-service/internal/pkg/errors/http"
	mock "auth-service/internal/pkg/mocks"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestVerifyValidate(t *testing.T) {
	request := RegisterVerifyRequest{
		Token:  "someToken",
		UserId: "some-uuid",
	}

	err := validateVerifyRequest(request)
	require.Nil(t, err)
}

func TestVerifyValidateError(t *testing.T) {
	request := RegisterVerifyRequest{
		Token:  "",
		UserId: "",
	}

	expErr := e.NewHttpError(400, "", fmt.Errorf("One of the parameters is empty."))

	err := validateVerifyRequest(request)
	require.NotNil(t, err)
	require.Equal(t, expErr, err)
}

func TestRegisterVerify(t *testing.T) {
	ctl := gomock.NewController(t)
	repositoryMock := mock.NewMockVerificationTokenInterface(ctl)
	grpcMock := mock.NewMockUserGrpcInterface(ctl)

	plainTokenPayload := "some-token"
	hashTokenPayload, _ := bcrypt.GenerateFromPassword([]byte(plainTokenPayload), 12)

	existToken := &m.VerificationToken{
		ID:        uuid.Must(uuid.NewV4()),
		UserId:    uuid.Must(uuid.NewV4()).String(),
		Token:     string(hashTokenPayload),
		ExpiresAt: time.Now().Add(time.Minute * 10),
		CreatedAt: time.Now(),
	}

	request := RegisterVerifyRequest{
		Token:  plainTokenPayload,
		UserId: existToken.UserId,
	}

	repositoryMock.EXPECT().Get(map[string]interface{}{"user_id": existToken.UserId}).Return(existToken, nil).Times(1)
	grpcMock.EXPECT().UpdateVerificationStatus(context.Background(), request.UserId).Return(true, nil).Times(1)
	repositoryMock.EXPECT().Delete(map[string]interface{}{"id": existToken.ID}).Return(nil).Times(1)

	resp, err := RegisterVerify(request, grpcMock, repositoryMock)

	require.Nil(t, err)
	require.Equal(t, &RegisterVerifyResponse{Verified: true}, resp)
}
