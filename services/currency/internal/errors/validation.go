package errors

type ValidationError struct {
	Msg string
}

func (e ValidationError) Error() string {
	return "validation error: " + e.Msg
}

func NewValidationError(msg string) *ValidationError {
	return &ValidationError{Msg: msg}
}
