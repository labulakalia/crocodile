package rbac

type Casbin struct {
	Id       int    `json:"id"`
	Ptype    string `json:"ptype"`
	RoleName string `json:"rolename"`
	Path     string `json:"path"`
	Method   string `json:"method"`
}

// 添加权限
