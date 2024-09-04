package errs

var (
	ErrNotFound     = New("not found")
	ErrValidation   = New("validation error")
	ErrInternal     = New("internal error")
	ErrInvalidParam = New("invalid parameter")

	ErrInsufficientBalance = New("insufficient balance")
	ErrOverLimit           = New("over limit")
)
