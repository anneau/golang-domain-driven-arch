package user_usecase

type CreateUserErrKind int

const (
	ErrEmailDuplicated CreateUserErrKind = iota
	ErrValidation
)

type CreateUserError struct {
	Kind    CreateUserErrKind
	Message string
}

func (e *CreateUserError) Error() string { return e.Message }
