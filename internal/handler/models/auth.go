package models

import "auth-service/internal/domain"

type (
	LoginRequestData struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	Tokens struct {
		AccessToken  domain.JwtTokenInfo `json:"access_token"`
		RefreshToken domain.JwtTokenInfo `json:"refresh_token"`
	}

	CreateRequestData struct {
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		Email     string `json:"email"`
	}

	CreateResponsedata struct {
		VerificationTokenId string `json:"verification_token_id"`
	}

	PasswordResetConfirmRequest struct {
		Token       string `json:"token"`
		TokenId     string `json:"token_id"`
		AccId       string `json:"acc_id"`
		NewPassword string `json:"new_password"`
	}
)
