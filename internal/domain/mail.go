package domain

type EmailType string

const (
	VerificationType EmailType = "verification_email"
	ResetType        EmailType = "reset_type"
)

type (
	EmailMessage struct {
		To      string    `json:"to"`
		Type    EmailType `json:"type"`
		Payload Payload
	}

	Payload struct {
		Firstname     string        `json:"firstname"`
		ResetPayload  ResetPayload  `json:"reset_payload"`
		VerifyPayload VerifyPayload `json:"verify_payload"`
	}

	ResetPayload struct {
		TokenId string `json:"token_id"`
		Token   string `json:"token"`
		AccId   string `json:"acc_id"`
	}

	VerifyPayload struct {
		Token string `json:"token"`
	}
)
