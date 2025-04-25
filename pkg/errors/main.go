package errors

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

type ErrorInfo struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
	Code       string `json:"code"`
	Help       string `json:"help,omitempty"`
	Err        error  `json:"-"`
	Type       string `json:"type,omitempty"`
}

type SystemError struct {
	*ErrorInfo
	StackTrace errors.StackTrace `json:"stack,omitempty" swaggertype:"object"`
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (e SystemError) Error() string {
	b, _ := json.Marshal(e)
	return string(b)

}

func (e SystemError) Cause() error {
	return errors.Cause(e.Err)
}

const (
	databasePrefix = "00"
	authPrefix     = "01"
	productPrefix  = "02"
	unknownPrefix  = "99"
)

func getStackTrace(e error) errors.StackTrace {
	err, ok := errors.Cause(e).(stackTracer)
	if !ok {
		return nil
	}
	return err.StackTrace()
}

func GetMessage(e error) string {
	if e == nil {
		return "unknown error"
	}
	if h, ok := e.(SystemError); ok {
		return h.Message
	} else {
		return e.Error()
	}
}

func GetCode(e error) string {
	if h, ok := e.(SystemError); ok {
		return h.Code
	} else {
		return "99999"
	}
}

func GetHelp(e error) string {
	if h, ok := e.(SystemError); ok {
		return h.Help
	} else {
		return ""
	}
}

func getStatusCode(e error) int {
	if h, ok := e.(SystemError); ok {
		return h.StatusCode
	} else {
		return http.StatusInternalServerError
	}
}

func wrap(e error, message string) error {
	if h, ok := e.(SystemError); ok {
		return errors.WithMessage(h.Err, message)
	}
	if e == nil {
		return errors.New(message)
	}
	_, ok := e.(stackTracer)
	if ok {
		return errors.WithStack(e)
	}
	return errors.New(e.Error())
}

func New(errorInfo ErrorInfo) SystemError {
	message := errorInfo.Message
	err := errorInfo.Err
	if message == "" {
		message = GetMessage(err)
	}

	code := errorInfo.Code
	if code == "" {
		code = GetCode(err)
	}

	statusCode := errorInfo.StatusCode
	if statusCode == 0 {
		statusCode = getStatusCode(err)
	}
	help := errorInfo.Help
	t := errorInfo.Type
	e := wrap(err, message)
	stackTrace := getStackTrace(e)
	return SystemError{
		ErrorInfo: &ErrorInfo{
			Message:    message,
			Code:       code,
			StatusCode: statusCode,
			Help:       help,
			Err:        e,
			Type:       t,
		},
		StackTrace: stackTrace,
	}
}

func ErrorUnknown(message string) SystemError {
	return New(ErrorInfo{
		Message: message,
	})
}
