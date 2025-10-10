package middleware

type HTTPError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e *HTTPError) Error() string {
	return e.Message
}

func NewHttpError(code int, msg string) *HTTPError {
	return &HTTPError{Code: code, Message: msg}
}
