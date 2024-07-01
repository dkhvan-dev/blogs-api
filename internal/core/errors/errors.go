package errors

type Error struct {
	Msg  string
	Code int
}

func (e Error) Error() string {
	return e.Msg
}

func NewNotFound(message string) Error {
	return Error{message, 404}
}

func NewBadRequest(message string) Error {
	return Error{message, 400}
}

func NewUnauthorized(message string) Error {
	return Error{message, 401}
}

func NewForbidden(message string) Error {
	return Error{message, 403}
}

func NewInternal(message string) Error {
	return Error{message, 500}
}
