package models

type (
	Token struct {
		UserId    string `db:"user_id"`
		Number    int64  `db:"number"`
		Purpose   int    `db:"purpose"`
		Secret    string `db:"secret"`
		ExpiresAt int64  `db:"expires_at"`
	}
)
