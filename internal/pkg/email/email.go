package email

import (
	"auth-service/internal/app/adapter"
	"fmt"
	"os"
)

type EmailType string

const (
	AccountVerification EmailType = "account_verification"
	AccountReset        EmailType = "account_reset"
)

type EmailPayload struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type Email struct {
	To string `json:"to"`
	EmailPayload
}

var Messages map[EmailType]EmailPayload = map[EmailType]EmailPayload{
	AccountVerification: {
		Subject: "Подтверждение электронной почты",
		Message: `
		Здравствуйте, %s!
		
		Для завершения регистрации перейдите, пожалуйста, по ссылке:
		%s
		
		Если вы не указывали эту электронную почту - проигнорируйте данное письмо.
		
		WarehouseAI Team
		`,
	},

	AccountReset: {
		Subject: "Восстановление пароля",
		Message: `
		Здравствуйте, %s!
		
		Мы получили запрос на восстановление пароля от аккаунта, связанного с почтой %s.
		Ваш код верификации: %s
		
		Если это не вы - проигнорируйте данное письмо.
		
		WarehouseAI Team
		`,
	},
}

func SendVerification(userId, token, firstname, email string, broker adapter.BrokerInterface) error {
	emailPayload := Messages[AccountVerification]
	emailPayload.Message = fmt.Sprintf(emailPayload.Message, firstname, fmt.Sprintf("%s/confirm?user=%s&token=%s", os.Getenv("MAIL_DOMAIN"), userId, token))
	message := Email{
		To:           email,
		EmailPayload: emailPayload,
	}

	if err := broker.SendMessage("mail", message); err != nil {
		return err
	}

	return nil
}
