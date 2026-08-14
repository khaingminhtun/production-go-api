package apperror

import "errors"

type Code string

const (
	CodeUserNotFound       Code = "USER_NOT_FOUND"
	CodeUserAlreadyExists  Code = "USER_ALREADY_EXISTS"
	CodeInvalidCredentials Code = "INVALID_CREDENTIALS"
	CodeEmailNotVerified   Code = "EMAIL_NOT_VERIFIED"
	CodeInvalidVerifyCode  Code = "INVALID_VERIFICATION_CODE"
	CodeVerifyCodeExpired  Code = "VERIFICATION_CODE_EXPIRED"
	CodeMethodNotAllowed   Code = "METHOD_NOT_ALLOWED"
	CodeRouteNotFound      Code = "ROUTE_NOT_FOUND"
	CodeInvalidRequest  Code = "INVALID_REQUEST"
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func New(code Code, message string, err error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func (e *Error) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}

	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func Is(err error, code Code) bool {
	var appErr *Error

	if !errors.As(err, &appErr) {
		return false
	}

	return appErr.Code == code
}
