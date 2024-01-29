package login

import (
	m "auth-service/internal/app/model"
	e "auth-service/internal/pkg/errors/http"
	mock "auth-service/internal/pkg/mocks"
	gen "auth-service/internal/pkg/protogen"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestValidateLogin(t *testing.T) {
	req := &LoginRequest{
		Email:    "validemail@mail.com",
		Password: "12345678",
	}

	err := validateLoginRequest(req)

	require.Nil(t, err)
}

func TestValidateLoginError(t *testing.T) {
	cases := []struct {
		name   string
		req    *LoginRequest
		expErr error
	}{
		{
			name: "invalid_email",
			req: &LoginRequest{
				Email:    "invalidemail",
				Password: "12345678",
			},
			expErr: e.NewHttpError(400, "", fmt.Errorf("Invalid email address")),
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := validateLoginRequest(tCase.req)

			require.NotNil(t, err)
			require.Equal(t, tCase.expErr.Error(), err.Error())
		})
	}
}

func TestLogin(t *testing.T) {
	ctl := gomock.NewController(t)

	grpcMock := mock.NewMockUserGrpcInterface(ctl)
	dbMock := mock.NewMockSessionInterface(ctl)

	request := &LoginRequest{
		Email:    "validemail@mail.com",
		Password: "12345678",
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte("12345678"), 12)
	expUser := &gen.User{
		Id:          uuid.Must(uuid.NewV4()).String(),
		Firstname:   "Firstname",
		Lastname:    "Lastname",
		Username:    "Username",
		Password:    string(hash),
		Email:       request.Email,
		ViaGoogle:   false,
		Picture:     "",
		Verified:    true,
		IsDeveloper: false,
		CreatedAt:   time.Now().String(),
		UpdatedAt:   time.Now().String(),
	}

	sessionPayload := m.SessionPayload{
		UserId:    expUser.Id,
		CreatedAt: time.Now(),
	}
	expSession := &m.Session{
		ID:      uuid.Must(uuid.NewV4()).String(),
		Payload: sessionPayload,
		TTL:     24 * time.Hour,
	}

	grpcMock.EXPECT().GetByEmail(context.Background(), request.Email).Return(expUser, nil).Times(1)
	dbMock.EXPECT().Create(context.Background(), expUser.Id).Return(expSession, nil).Times(1)

	resp, session, err := Login(request, grpcMock, dbMock)

	require.NotNil(t, resp)
	require.NotNil(t, session)
	require.Nil(t, err)
	require.Equal(t, &LoginResponse{UserId: expUser.Id}, resp)
	require.IsType(t, &m.Session{}, session)
}

func TestLoginError(t *testing.T) {
	cases := []struct {
		name   string
		req    *LoginRequest
		expErr error
	}{
		{
			name: "user_not_found",
			req: &LoginRequest{
				Email:    "notexistemail@mail.com",
				Password: "12345678",
			},
			expErr: e.NewHttpError(404, "", fmt.Errorf("User not found.")),
		},
	}

	ctl := gomock.NewController(t)

	grpcMock := mock.NewMockUserGrpcInterface(ctl)
	dbMock := mock.NewMockSessionInterface(ctl)

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			grpcMock.EXPECT().GetByEmail(context.Background(), tCase.req.Email).Return(nil, tCase.expErr).Times(1)

			resp, session, err := Login(tCase.req, grpcMock, dbMock)

			require.Nil(t, resp)
			require.Nil(t, session)
			require.NotNil(t, err)
			require.Equal(t, tCase.expErr.Error(), err.Error())
		})
	}
}
