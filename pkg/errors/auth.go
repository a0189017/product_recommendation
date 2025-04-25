package errors

import (
	"fmt"
	"net/http"
)

func ErrorTokenNotFound() SystemError {
	return SystemError{
		ErrorInfo: &ErrorInfo{
			Message:    "Token not found",
			StatusCode: http.StatusUnauthorized,
			Code:       fmt.Sprintf("%s%s", authPrefix, "000"),
			Help:       "Add Authorization with Bearer {access_token} obtained from the oauth API in the header",
		},
	}
}

func ErrorUserPasswordFormatInvalid() SystemError {
	return SystemError{
		ErrorInfo: &ErrorInfo{
			Message:    "User password format invalid",
			StatusCode: http.StatusBadRequest,
			Code:       fmt.Sprintf("%s%s", authPrefix, "001"),
			Help:       "Password must be 6-16 characters long, contain at least one uppercase letter, one lowercase letter, one number, and one special character",
		},
	}
}

func ErrorUserAlreadyExists() SystemError {
	return SystemError{
		ErrorInfo: &ErrorInfo{
			Message:    "User already exists",
			StatusCode: http.StatusBadRequest,
			Code:       fmt.Sprintf("%s%s", authPrefix, "002"),
			Help:       "Please use another email",
		},
	}
}

func ErrorInvalidPassword() SystemError {
	return SystemError{
		ErrorInfo: &ErrorInfo{
			Message:    "Invalid password",
			StatusCode: http.StatusBadRequest,
			Code:       fmt.Sprintf("%s%s", authPrefix, "003"),
			Help:       "Please check your password",
		},
	}
}

func ErrorMissingOTP() SystemError {
	return SystemError{
		ErrorInfo: &ErrorInfo{
			Message:    "Missing OTP",
			StatusCode: http.StatusBadRequest,
			Code:       fmt.Sprintf("%s%s", authPrefix, "004"),
			Help:       "Please fill in the OTP",
		},
	}
}

func ErrorInvalidOTP() SystemError {
	return SystemError{
		ErrorInfo: &ErrorInfo{
			Message:    "Invalid OTP",
			StatusCode: http.StatusBadRequest,
			Code:       fmt.Sprintf("%s%s", authPrefix, "005"),
			Help:       "Please check your OTP",
		},
	}
}

func ErrorUserNotFound() SystemError {
	return SystemError{
		ErrorInfo: &ErrorInfo{
			Message:    "User not found",
			StatusCode: http.StatusNotFound,
			Code:       fmt.Sprintf("%s%s", authPrefix, "006"),
			Help:       "Please check your account",
		},
	}
}

func ErrorTokenMalformed() SystemError {
	return New(ErrorInfo{
		Message:    "That's not a token",
		StatusCode: http.StatusUnauthorized,
		Code:       fmt.Sprintf("%s%s", authPrefix, "007"),
		Help:       "Add Authorization with Bearer {access_token} obtained from the oauth API in the header",
	})
}

func ErrorTokenExpired() SystemError {
	return New(ErrorInfo{
		Message:    "Token is expired",
		StatusCode: http.StatusUnauthorized,
		Code:       fmt.Sprintf("%s%s", authPrefix, "008"),
		Help:       "Please log in again",
	})
}

func ErrorTokenSignatureInvalid() SystemError {
	return New(ErrorInfo{
		Message:    "Token signature is invalid",
		StatusCode: http.StatusUnauthorized,
		Code:       fmt.Sprintf("%s%s", authPrefix, "009"),
		Help:       "Add Authorization with Bearer {access_token} obtained from the oauth API in the header",
	})
}

func ErrorTokenNotValid() SystemError {
	return New(ErrorInfo{
		Message:    "Token is not valid",
		StatusCode: http.StatusUnauthorized,
		Code:       fmt.Sprintf("%s%s", authPrefix, "010"),
		Help:       "Please log in again",
	})
}

func ErrorOTPExpired() SystemError {
	return SystemError{
		ErrorInfo: &ErrorInfo{
			Message:    "OTP expired",
			StatusCode: http.StatusBadRequest,
			Code:       fmt.Sprintf("%s%s", authPrefix, "011"),
			Help:       "Please request a new OTP",
		},
	}
}

func ErrorUserAccountFormatInvalid() SystemError {
	return SystemError{
		ErrorInfo: &ErrorInfo{
			Message:    "User account format invalid",
			StatusCode: http.StatusBadRequest,
			Code:       fmt.Sprintf("%s%s", authPrefix, "012"),
			Help:       "Please check your account is valid email",
		},
	}
}

func ErrorOTPAlreadyUsed() SystemError {
	return SystemError{
		ErrorInfo: &ErrorInfo{
			Message:    "OTP already used",
			StatusCode: http.StatusBadRequest,
			Code:       fmt.Sprintf("%s%s", authPrefix, "013"),
			Help:       "Please request a new OTP",
		},
	}
}
