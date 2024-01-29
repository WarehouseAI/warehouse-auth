package register

import (
	m "auth-service/internal/app/model"
	"auth-service/internal/pkg/email"
	httpe "auth-service/internal/pkg/errors/http"
	mock "auth-service/internal/pkg/mocks"
	gen "auth-service/internal/pkg/protogen"
	"context"
	"fmt"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

// Valid request
func TestRegisterValidate(t *testing.T) {
	request := &RegisterRequest{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Email:     "validmail@mail.com",
		Username:  "Username",
		Password:  "12345678",
		// Image:     "",
		ViaGoogle: false,
	}

	err := validateRegisterRequest(request)
	require.Nil(t, err)
}

func TestValidateError(t *testing.T) {
	cases := []struct {
		name          string
		request       *RegisterRequest
		expectedError error
	}{
		{
			name: "long_password",
			request: &RegisterRequest{
				Firstname: "Firstname",
				Lastname:  "Lastname",
				Email:     "validmail@mail.com",
				Username:  "Username",
				Password:  "rqrZBhrHzy9tnNTbL9HzPaAYdtnMqVJ4qEQBkrY77bP5GiaceM5op8642FB3DRMGRA9kSsvaa",
				// Image:     "",
				ViaGoogle: false,
			},
			expectedError: httpe.NewHttpError(400, "", fmt.Errorf("Password is too long")),
		},
		{
			name: "short_password",
			request: &RegisterRequest{
				Firstname: "Firstname",
				Lastname:  "Lastname",
				Email:     "validmail@mail.com",
				Username:  "Username",
				Password:  "1234567",
				// Image:     "",
				ViaGoogle: false,
			},
			expectedError: httpe.NewHttpError(400, "", fmt.Errorf("Password is too short")),
		},
		{
			name: "bad_email",
			request: &RegisterRequest{
				Firstname: "Firstname",
				Lastname:  "Lastname",
				Email:     "validmail",
				Username:  "Username",
				Password:  "12345678",
				// Image:     "",
				ViaGoogle: false,
			},
			expectedError: httpe.NewHttpError(400, "", fmt.Errorf("The provided string is not email")),
		},
	}

	for _, tCase := range cases {
		t.Run(tCase.name, func(t *testing.T) {
			err := validateRegisterRequest(tCase.request)

			require.NotNil(t, err)
			require.Equal(t, tCase.expectedError.Error(), err.Error())
		})
	}
}

func TestRegister(t *testing.T) {
	ctl := gomock.NewController(t)

	grpcMock := mock.NewMockUserGrpcInterface(ctl)
	dbMock := mock.NewMockVerificationTokenInterface(ctl)
	brokerMock := mock.NewMockBrokerInterface(ctl)

	request := &RegisterRequest{
		Firstname: "Firstname",
		Lastname:  "Lastname",
		Username:  "Username",
		Password:  "12345678",
		// Image:     "",
		Email:     "validmail@mail.com",
		ViaGoogle: false,
	}

	userId := uuid.Must(uuid.NewV4()).String()

	grpcMock.EXPECT().Create(context.Background(), gomock.AssignableToTypeOf(&gen.CreateUserMsg{})).Return(userId, nil).Times(1)
	dbMock.EXPECT().Create(gomock.AssignableToTypeOf(&m.VerificationToken{})).Return(nil).Times(1)
	brokerMock.EXPECT().SendMessage("mail", gomock.AssignableToTypeOf(email.Email{})).Return(nil).Times(1)

	resp, err := Register(request, grpcMock, dbMock, brokerMock)

	require.Nil(t, err)
	require.Equal(t, &RegisterResponse{UserId: userId}, resp)
}

// TODO: переделать тест ошибочных запросов под новые ошибки
// func TestRegisterGrpcError(t *testing.T) {
// 	cases := []struct {
// 		name   string
// 		req    *RegisterRequest
// 	}{
// 		{
// 			name: "already_exist",
// 			req: &RegisterRequest{
// 				Firstname: "Firstname",
// 				Lastname:  "Lastname",
// 				Username:  "Username",
// 				Password:  "12345678",
// 				Image:     "",
// 				Email:     "validmail@mail.com",
// 				ViaGoogle: false,
// 			},
// 		},
// 		{
// 			name: "internal_error",
// 			req: &RegisterRequest{
// 				Firstname: "Firstname",
// 				Lastname:  "Lastname",
// 				Username:  "Username",
// 				Password:  "12345678",
// 				Image:     "",
// 				Email:     "validmail@mail.com",
// 				ViaGoogle: false,
// 			},
// 		},
// 	}

// 	ctl := gomock.NewController(t)

// 	grpcMock := mock.NewMockUserGrpcInterface(ctl)
// 	dbMock := mock.NewMockVerificationTokenInterface(ctl)
// 	brokerMock := mock.NewMockBrokerInterface(ctl)
// 	logger := logrus.New()

// 	for _, tCase := range cases {
// 		t.Run(tCase.name, func(t *testing.T) {
// 			grpcMock.EXPECT().Create(context.Background(), gomock.AssignableToTypeOf(&gen.CreateUserMsg{})).Return("", error).Times(1)

// 			resp, err := Register(tCase.req, grpcMock, dbMock, brokerMock, logger)

// 			require.NotNil(t, err)
// 			require.Equal(t, tCase.expErr, err)
// 			require.Nil(t, resp)
// 		})
// 	}
// }

// func TestRegisterDbError(t *testing.T) {
// 	cases := []struct {
// 		name   string
// 		req    *RegisterRequest
// 		expErr *ehttpeDBError
// 	}{
// 		{
// 			name: "already_exist",
// 			req: &RegisterRequest{
// 				Firstname: "Firstname",
// 				Lastname:  "Lastname",
// 				Username:  "Username",
// 				Password:  "12345678",
// 				Image:     "",
// 				Email:     "validmail@mail.com",
// 				ViaGoogle: false,
// 			},
// 			expErr: &ehttpeDBError{
// 				ErrorType: ehttpeDbExist,
// 				Message:   "Token with this key/keys already exists.",
// 				Payload:   "token already exists payload",
// 			},
// 		},
// 		{
// 			name: "internal_error",
// 			req: &RegisterRequest{
// 				Firstname: "Firstname",
// 				Lastname:  "Lastname",
// 				Username:  "Username",
// 				Password:  "12345678",
// 				Image:     "",
// 				Email:     "validmail@mail.com",
// 				ViaGoogle: false,
// 			},
// 			expErr: &ehttpeDBError{
// 				ErrorType: ehttpeDbSystem,
// 				Message:   "Something went wrong.",
// 				Payload:   "internal error payload",
// 			},
// 		},
// 	}

// 	ctl := gomock.NewController(t)

// 	grpcMock := aMock.NewMockUserGrpcInterface(ctl)
// 	dbMock := dMock.NewMockVerificationTokenInterface(ctl)
// 	brokerMock := aMock.NewMockBrokerInterface(ctl)
// 	logger := logrus.New()

// 	for _, tCase := range cases {
// 		t.Run(tCase.name, func(t *testing.T) {
// 			grpcMock.EXPECT().Create(context.Background(), gomock.AssignableToTypeOf(&gen.CreateUserMsg{})).Return("id", nil).Times(1)
// 			dbMock.EXPECT().Create(gomock.AssignableToTypeOf(&m.VerificationToken{})).Return(tCase.expErr).Times(1)
// 			brokerMock.EXPECT().SendUserReject("id").Return(nil).Times(1)

// 			resp, err := Register(tCase.req, grpcMock, dbMock, brokerMock, logger)

// 			require.NotNil(t, err)
// 			require.Equal(t, ehttpeNewErrorResponseFromDBError(tCase.expErr.ErrorType, tCase.expErr.Message), err)
// 			require.Nil(t, resp)
// 		})
// 	}
// }
