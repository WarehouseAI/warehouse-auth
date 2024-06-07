package domain

var (
	Roles = []Role{RoleAdmin, RoleUser}
)

type Role int64

const (
	RoleAdmin Role = iota
	RoleUser
)

type AuthPurpose int64

const (
	PurposeAccess = AuthPurpose(iota)
	PurposeRefresh
)

type (
	Account struct {
		Id        string `json:"id"`
		Role      Role   `json:"role"`
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Email     string `json:"email"`
		Verified  bool   `json:"verified"`
	}

	JwtTokenInfo struct {
		Token     string `json:"token"`
		ExpiresAt int64  `json:"expires_at"`
	}

	VerificationTokenInfo struct {
		ID        string
		UserId    string
		Token     string
		SendTo    string
		ExpiresAt int64
		CreatedAt int64
	}

	ResetTokenInfo struct {
		ID        string
		UserId    string
		Token     string
		ExpiresAt int64
		CreatedAt int64
	}
)

//Request

type (
	ResetPasswordRequestData struct {
		Id       string
		Password string
	}

	UpdateVerificationRequestData struct {
		Id    string
		Email string
	}
)
