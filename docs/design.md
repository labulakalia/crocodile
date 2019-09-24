## 用户权限
```go
type roleId int
type User struct {
    Id int `json:"id""` 
    Name string `json:"name""`

    Avatar string `json:"avatar"`             // 用户图像
    Role string `json:"role"`                   // 用户类型 regular admin

    Forbid int `json:"forbid"`                  // 禁止用户
    Email string `json:"email"`                 // 用户邮箱 日后任务的通知信息会发送给此邮件


    CreateTime int64 `json:"create_time""`     // 创建时间
    UpdateTime int64 `json:"update_time""`   // 最后一次更新时间
    Remark  string `json:"remark"`              // 备注
}
```
角色 一个角色关联多个权限
```go
type permissionId int
type Role struct {
    Id int `json:"id"`
    Name string `json:"name"`
    PermissionIds []permissionId `json:"permission_ids"` // 角色所包含的权限
    CreateTime int64 `json:"create_time""`     // 创建时间
    UpdateTime int64 `json:"update_time""`   // 最后一次更新时间
}
```
权限 一条权限由 [read|editor]+[path] 组成
```go
type Permission struct {
    Id int `json:"id"`
    Name string `json:"name"`
    Operate string `json:"operate"'` // 操作 读 修改 关联操作如果对一个path 具有修改操作那么必须有
    Path string `json:"path"`   // 分组路径
    CreateTime int64 `json:"create_time"`     // 创建时间
    UpdateTime int64 `json:"update_time"`   // 最后一次更新时间
}
```

## 主机组
```go
type HostGroup {
    Id int `json:"id"`
    Name string `json:"name"`

    Hosts []string `json:"hosts"`

    CreateTime int64 `json:"create_time""`     // 创建时间
    UpdateTime int64 `json:"update_time""`   // 最后一次更新时间
    Remark  string `json:"remark"`              // 备注
}
```

## 任务
```go
type Task {
    Id int `json:"id"`
    Name string `json:"name"`

    HostGroupId int `json:"host_group_id"` // 主机组Id
    HostGroup string `json:"host_group"`  // 主机组名称
    Timeout int `json:"timeout"`    // 0 为不启动超时控制
    Run int `json:"run"`  // 0 为不能运行 1 为可以运行

    ParentTaskIds []int `json:"parent_task_ids"`  // 父任务 运行任务前先运行父任务 以父或子任务运行时 任务不会执行自已的父子任务，防止循环依赖
    ParentRunParallel int `json:"parent_run_parallel"`  // 是否以并行运行父任务 0否 1是 

    ChildTaskIds []int `json:"child_task_ids"`    // 子任务 运行结束后运行子任务
    ChildRunParallel int   `json:"child_run_parallel"`  // 是否以并行运行子任务 0否 1是

    RunByTask int `json:"run_by_task"` // 通过其他的任务调用运行 如果是被其他任务依赖而运行就不会运行此任务的父子任务 1 为 true

    RunTime int `json:"run_time"`     // 运行次数 默认为0则不限制
    
    AlarmTotal int `json:"alarm_total"`  // 是否报警次数

    RetryRun int `json:"retry_run"`         // 重试次数 0 为不重试

    CreateByUserId int `json:"create_by_user_id"`  // 创建人ID
    CreateByUser string `json:"create_by_user"`   // 创建人
    
    Remark  string `json:"remark"`              // 备注
    CreateTime int64 `json:"create_time"`     // 创建时间
    UpdateTime int64 `json:"update_time""`   // 最后一次更新时间
}
```

## 技术点
- 用户权限管理
  casbin

- 用户认证
  jwt

- 任务调用
  grpc 用于运行任务时调用worker来运行指定的任务信息，
  master启动后会从配置文件中加载主从节点加载配置到etcd存储，然后worker会通过grpc来拉取自已的配置，然后加载配置运行，会会watch配置服务，实时更新，
  然后会将自已启动信息上报至服务端，
  worker启动后先会从拉取运行所需要的会将自已的IP:PORT注册至etcd中，
- 任务调度
  会将所有的cron表达式和对应ID加载至内存，计算下一次的运行时间，然后通过对应的任务ID去运行任务
- 日志管理
  zap
  
- 定时任务
只将任务ID、任务的cron表示式、下一次运行时间存储在内存中，
其余任务的信息存储在数据库中，任务到期后通过任务ID去数据库取相应的任务信息然后运行


## Required
尽量模块化设计，使得每一部分都可以单独使用

swagger API