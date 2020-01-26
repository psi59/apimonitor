package amerr

import (
	"fmt"
	"net/http"
)

const (
	ErrBadRequest        = 400
	ErrUnsupportedMethod = 40001

	ErrNotFound           = 404
	ErrWebServiceNotFound = 4041
	ErrTestNotFound       = 4042

	ErrConflict             = 409
	ErrDuplicatedWebService = 4091
	ErrDuplicatedTest       = 4092

	ErrInternalServer = 500
)

var (
	errors = map[int]*ErrorWithLanguage{
		ErrBadRequest: newErrorWithLanguage(
			newError(http.StatusBadRequest, ErrBadRequest, "잘못된 요청으로 인해 작업을 완료할 수 없습니다."),
		),
		ErrInternalServer: newErrorWithLanguage(
			newError(http.StatusInternalServerError, ErrInternalServer, "서버 내부에러로 인해 작업을 완료할 수 없습니다."),
		),
		ErrConflict: newErrorWithLanguage(
			newError(http.StatusConflict, ErrConflict, "중복된 요청입니다."),
		),

		ErrUnsupportedMethod: newErrorWithLanguage(
			newError(http.StatusBadRequest, ErrUnsupportedMethod, "지원하지 않는 메소드 입니다."),
		),

		ErrNotFound: newErrorWithLanguage(
			newError(http.StatusNotFound, ErrNotFound, "해당 요청을 찾을 수 없습니다."),
		),
		ErrWebServiceNotFound: newErrorWithLanguage(
			newError(http.StatusNotFound, ErrWebServiceNotFound, "해당 웹서비스를 찾을 수 없습니다."),
		),
		ErrTestNotFound: newErrorWithLanguage(
			newError(http.StatusNotFound, ErrTestNotFound, "해당 테스트를 찾을 수 없습니다."),
		),

		ErrDuplicatedWebService: newErrorWithLanguage(
			newError(http.StatusConflict, ErrDuplicatedWebService, "이미 같은 호스트의 웹서비스가 존재합니다."),
		),
		ErrDuplicatedTest: newErrorWithLanguage(
			newError(http.StatusConflict, ErrDuplicatedTest, "이미 같은 테스트가 존재합니다."),
		),
	}
)

type Error struct {
	StatusCode int
	ErrorCode  int
	Message    string
}

func (e Error) Error() string {
	return fmt.Sprintf("status_code='%d', error_code='%d', msg='%s'", e.StatusCode, e.ErrorCode, e.Message)
}

func newError(statusCode, errorCode int, msg string) error {
	return &Error{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    msg,
	}
}

type ErrorWithLanguage map[string]error

func (e ErrorWithLanguage) GetErrFromLanguage(lang string) error {
	return e[lang]
}

func newErrorWithLanguage(koErr error) *ErrorWithLanguage {
	return &ErrorWithLanguage{
		"ko": koErr,
	}
}

func GetErrorsFromCode(code int) *ErrorWithLanguage {
	return errors[code]
}

func GetErrInternalServer() *ErrorWithLanguage {
	return GetErrorsFromCode(ErrInternalServer)
}
