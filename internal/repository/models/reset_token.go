package models

import (
	"github.com/rs/xid"
)

type ResetToken struct {
	ID        xid.ID `db:"id"`
	UserId    xid.ID `db:"user_id"`
	Token     string `db:"token"`
	ExpiresAt int64  `db:"expires_at"`
	CreatedAt int64  `db:"created_at"`
}
