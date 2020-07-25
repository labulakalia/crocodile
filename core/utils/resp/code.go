package resp

const (
	// Success Success
	Success = 0
	// ErrBadRequest 请求参数错误
	ErrBadRequest = 10400
	// ErrUnauthorized 非法请求
	ErrUnauthorized = 10401
	// ErrUserPassword 用户密码错误
	ErrUserPassword = 10402
	// ErrUserForbid 禁止登陆
	ErrUserForbid = 10403
	// ErrUserNameExist 邮箱已经存在
	ErrUserNameExist = 10413
	// ErrEmailExist 用户名已存在
	ErrEmailExist = 10414
	// ErrUserNotExist 用户不存在
	ErrUserNotExist = 10415

	// ErrTaskExist 任务名已存在
	ErrTaskExist = 10416
	// ErrTaskNotExist 任务不存在
	ErrTaskNotExist = 10417

	// ErrHostgroupExist 主机组已存在
	ErrHostgroupExist = 10418
	// ErrHostgroupNotExist 主机组不存在
	ErrHostgroupNotExist = 10419
	// ErrDelHostUseByOtherHG 正在被其他的主机组使用，不能删除
	ErrDelHostUseByOtherHG = 10420
	//ErrHostNotExist 主机不存在
	ErrHostNotExist = 10421

	// ErrCronExpr CronExpr表达式不规范
	ErrCronExpr = 10422

	// ErrTaskUseByOtherTask 别的任务依赖此任务，请先在其他的任务的父子任务中移除此任务
	ErrTaskUseByOtherTask = 10423

	// ErrDelHostGroupUseByTask 正在被其他的任务使用，不能删除
	ErrDelHostGroupUseByTask = 10424
	// ErrDelUserUseByOther // 请先删除此用户创建的主机组或者任务后再删除
	ErrDelUserUseByOther = 10425

	// ErrInternalServer 服务端错误
	ErrInternalServer = 10500
	// ErrCtxDeadlineExceeded 调用超时
	ErrCtxDeadlineExceeded = 10600
	// ErrCtxCanceled 取消调用
	ErrCtxCanceled = 10601

	// ErrRPCUnauthenticated  密钥认证失败
	ErrRPCUnauthenticated = 10602
	// ErrRPCUnavailable 调用对端不可用
	ErrRPCUnavailable = 10603
	// ErrRPCUnknow 调用未知错误
	ErrRPCUnknow = 10604
	// ErrRPCNotValidHost  未发现worker
	ErrRPCNotValidHost = 10605
	// ErrRPCNotConnHost 未找到存活的worker
	ErrRPCNotConnHost = 10606

	// NeedInstall 系统还未安装，请等待安装后再进行操作
	NeedInstall = 10700
	// IsInstall 系统已经安装完成，请勿再次执行
	IsInstall = 10701
	// ErrInstall 安装失败
	ErrInstall = 10702
	// ErrDBConnFail 数据库连接失败
	ErrDBConnFail = 10703
)
