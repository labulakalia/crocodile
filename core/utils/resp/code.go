package resp

const (
	Success = 0

	ErrBadRequest   = 10400
	ErrUnauthorized = 10401

	ErrUserPassword  = 10402
	ErrUserForbid    = 10403
	ErrUserNameExist = 10413
	ErrEmailExist    = 10414
	ErrUserNotExist  = 10415

	ErrInternalServer = 10500
)
