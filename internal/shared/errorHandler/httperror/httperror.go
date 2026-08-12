package httperror

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/khaingminhtun/production-go-api/internal/shared/errorHandler/apperror"
)

type Error struct {
	Status  int
	Code    string
	Message string
	Err     error
}

func New(status int, code, message string, err error) *Error {
	return &Error{
		Status:  status,
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

func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	// Already an HTTP error.
	var httpErr *Error
	if errors.As(err, &httpErr) {
		return httpErr
	}

	// Application/business error.
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		return fromAppError(appErr)
	}

	// GORM "not found".
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return New(
			http.StatusNotFound,
			"NOT_FOUND",
			"resource not found",
			err,
		)
	}

	// Unknown error.
	return New(
		http.StatusInternalServerError,
		"INTERNAL_SERVER_ERROR",
		"internal server error",
		err,
	)
}

func fromAppError(err *apperror.Error) *Error {
	switch err.Code {

	case apperror.CodeUserNotFound:
		return New(
			http.StatusNotFound,
			string(err.Code),
			err.Message,
			err,
		)

	case apperror.CodeUserAlreadyExists:
		return New(
			http.StatusConflict,
			string(err.Code),
			err.Message,
			err,
		)

	case apperror.CodeInvalidCredentials:
		return New(
			http.StatusUnauthorized,
			string(err.Code),
			err.Message,
			err,
		)

	case apperror.CodeEmailNotVerified:
		return New(
			http.StatusForbidden,
			string(err.Code),
			err.Message,
			err,
		)

	case apperror.CodeInvalidVerifyCode:
		return New(
			http.StatusBadRequest,
			string(err.Code),
			err.Message,
			err,
		)

	case apperror.CodeVerifyCodeExpired:
		return New(
			http.StatusBadRequest,
			string(err.Code),
			err.Message,
			err,
		)

	default:
		return New(
			http.StatusInternalServerError,
			"INTERNAL_SERVER_ERROR",
			"internal server error",
			err,
		)
	}
}
