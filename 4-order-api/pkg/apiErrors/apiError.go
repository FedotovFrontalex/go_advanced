package apierrors

type Error struct {
	status int
	text   string
}

func NewError(status int, text string) *Error {
	return &Error{
		status: status,
		text:   text,
	}
}

func (err *Error) Error() string {
	return err.text
}

func (err *Error) GetStatus() int {
	return err.status
}
