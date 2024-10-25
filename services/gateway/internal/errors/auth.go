package errors

type AuthError struct {
	Message string
	Cause   error
}

func (e AuthError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}

	return e.Message
}

func NewAuthError(msg string) *AuthError {
	return &AuthError{Message: msg}
}

func NewAuthErrorWithCause(cause error, msg string) *AuthError {
	return &AuthError{Cause: cause, Message: msg}
}

const (
	AuthErrInvalidLoginOrPasswordMsg = "invalid login or password"
	AuthErrorInvalidTokenFormatMsg   = "invalid token format"
	AuthErrorTokenExpiredMsg         = "token expired or invalid"
)
