package domain

import "errors"

var (
	// general
	ErrInternalServerError = errors.New("internal server error")
	ErrInvalidContextKey   = errors.New("invalid context key")

	// transport layer error
	ErrInvalidRequest = errors.New("invalid request")
	ErrForbiden       = errors.New("forbiden request")
	ErrToomanyrequest = errors.New("too many request")

	// service error
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrMerchantAlreadyExists    = errors.New("merchant already exists")
	ErrMerchantKYCAlreadyExists = errors.New("merchant already KYC, please wait for approval")

	// jwt error
	ErrInvalidToken = errors.New("invalid token")

	// auth error
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrTokenRequired      = errors.New("token required")

	// user error
	ErrEmailAlreadyExists    = errors.New("email already exist")
	ErrPhoneAlreadyExists    = errors.New("phone already exist")
	ErrNicknameAlreadyExists = errors.New("nickname already exist")
)
