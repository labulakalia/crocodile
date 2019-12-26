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

	ErrTaskExist    = 10416
	ErrTaskNotExist = 10417

	ErrHostgroupExist    = 10418
	ErrHostgroupNotExist = 10419

	ErrExecPlanExist    = 10420
	ErrExecPlanNotExist = 10421

	ErrInternalServer = 10500

	ErrRpcDeadlineExceeded = 10600
	ErrRpcCanceled         = 10601
	ErrRpcUnauthenticated  = 10602
	ErrRpcUnavailable      = 10603
	ErrRpcUnknow           = 10604
	ErrRpcNotValidHost     = 10605
	ErrRpcNotConn          = 10606
)
