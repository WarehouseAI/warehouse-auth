package models

import (
	"github.com/rs/xid"
)

type VerificationToken struct {
	ID        xid.ID `db:"id"`
	UserId    xid.ID `db:"user_id"`
	Token     string `db:"token"`
	SendTo    string `db:"send_to"`
	ExpiresAt int64  `db:"expires_at"`
	CreatedAt int64  `db:"created_at"`
}
