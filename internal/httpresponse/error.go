package httpresponse

type Error struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  string `json:"errors,omitempty"`
}

func Fail(message, errors string) Error {
	return Error{
		Success: false,
		Message: message,
		Errors:  errors,
	}
}
