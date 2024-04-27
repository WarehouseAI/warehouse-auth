package errors

import "errors"

var (
	AuthAuthFailedRaw = errors.New("auth failed")
	AuthAuthFailed    = &Error{Code: 401, Reason: "auth failed"}

	AuthParseTokenRaw = errors.New("parse token failed")
	AuthParseToken    = &Error{Code: 400, Reason: AuthParseTokenRaw.Error()}

	AuthHashPassword        = &Error{Code: 400, Reason: "hashing password error"}
	AuthExpiredToken        = &Error{Code: 400, Reason: "expired token"}
	AuthInvalidTokenPurpose = &Error{Code: 400, Reason: "invalid token purpose"}
	AuthInvalidToken        = &Error{Code: 400, Reason: "invalid token"}
	AuthCreateTokens        = &Error{Code: 400, Reason: "create tokens error"}
	AuthVerificationFailed  = &Error{Code: 400, Reason: "account was not successfully updated"}
	AuthNotVerifiedAccount  = &Error{Code: 403, Reason: "account not verified yet"}

	AuthUserAlreadyExists = &Error{Code: 409, Reason: "user already exists"}

	AuthUserNotFoundByIdRaw = errors.New("there is no user with such id")
	AuthUserNotFoundById    = &Error{Code: 400, Reason: AuthUserNotFoundByIdRaw.Error()}

	AuthNumberAssignmentFailedRaw = errors.New("number assignment failed")
	AuthNumberAssignmentFailed    = &Error{Code: 400, Reason: AuthNumberAssignmentFailedRaw.Error()}

	AuthGetUserDataFailed = &Error{Code: 400, Reason: "get user data failed"}

	CreateToken       = errors.New("token not created")
	TokenDoesNotExist = errors.New("token does not exist")
)
