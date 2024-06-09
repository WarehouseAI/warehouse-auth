package jwt

import "github.com/warehouse/auth-service/internal/domain"

var UserTokenMap = map[domain.Role]string{
	domain.RoleAdmin: "admin_tokens",
	domain.RoleUser:  "user_tokens",
}
