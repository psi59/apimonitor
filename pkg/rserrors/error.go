package rserrors

const (
	ErrInvalidParameter Error = "invalid parameter"
	ErrUnexpected       Error = "unexpected error"
)

type Error string

func (err Error) String() string {
	return string(err)
}

func (err Error) Error() string {
	return err.String()
}
