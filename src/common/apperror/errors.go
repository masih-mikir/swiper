package apperror

import "errors"

var (
	StatusBadRequest    = errors.New("Status Bad Request")
	InternalServerError = errors.New("Internal Server Error")
	DecodeError         = errors.New("Wrong request params format, see example in data")
	AccountNotExists    = errors.New("Account not exists")
	RecreationNotExists = errors.New("Recreation not exists")
)

type ErrorCodes struct {
	HTTPcode   int
	StatusCode int
}

var (
	DefaultErrorCode  = ErrorCodes{400, 100000}
	NotFoundErrorCode = ErrorCodes{0, 0}
)

var errorCodes = map[error]ErrorCodes{
	StatusBadRequest:    ErrorCodes{400, 100101},
	InternalServerError: ErrorCodes{500, 100102},
	DecodeError:         ErrorCodes{400, 100201},
	AccountNotExists:    ErrorCodes{400, 200010},
	RecreationNotExists: ErrorCodes{400, 300010},
}

func GetErrorCodes(err error) ErrorCodes {
	error := errorCodes[err]
	if error == NotFoundErrorCode {
		return DefaultErrorCode
	}
	return error
}
