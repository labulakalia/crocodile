package model

import (
	"context"
	"testing"

	"github.com/labulaka521/crocodile/core/utils/define"
)

func TestCreateTaskv2(t *testing.T) {

	type args struct {
		ctx  context.Context
		task *Task
	}
	t.Run("create user", TestAddUserv2)
	users, count, err := GetUsersv2(context.Background(), nil, 0, 10)
	if err != nil {
		t.Fatalf("get users v2 failed: %v", err)
	}
	if count == 0 {
		t.Fatal("users count should not get 0")
	}
	t.Run("create hostgroup", TestCreateHostgroupv2)
	hgs, count, err := GetHostGroupsv2(context.Background(), 0, 10)
	if count == 0 {
		t.Fatal("hgs count should not get 0")
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		// 等 User
		{
			name: "create task",
			args: args{
				ctx: context.Background(),
				task: &Task{
					Name:          "test task name",
					TaskType:      define.API,
					Run:           true,
					ParentTaskIDs: IDs{"111111111111111111"},
					CreateUID:     users[0].ID,
					HostgroupID:   hgs[0].ID,
					Cronexpr:      "* * * * * * *",
					Timeout:       -1,
					AlarmUIDs:     IDs{hgs[0].ID},
					RoutePolicy:   define.LeastTask,
					ExpectCode:    0,
					AlarmStatus:   define.Fail,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateTaskv2(tt.args.ctx, tt.args.task)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTaskv2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestChangeTaskv2(t *testing.T) {
	t.Run("create tasks", TestCreateTaskv2)
	tasks, count, err := GetTasksv2(context.Background(), 0, 10, "", "")
	if err != nil {
		t.Fatal("get task failed ", err)
	}
	if count == 0 {
		t.Fatal("count not should 0")
	}
	task := tasks[0]
	task.Run = false
	task.RoutePolicy = define.Weight

	notexistTask := task
	notexistTask.ID = "not exist idwo;enowenf"
	type args struct {
		ctx  context.Context
		task *Task
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "chnage task",
			args: args{
				ctx:  context.Background(),
				task: task,
			},
			wantErr: false,
		},
		{
			name: "chnage not existtask",
			args: args{
				ctx:  context.Background(),
				task: notexistTask,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ChangeTaskv2(tt.args.ctx, tt.args.task); (err != nil) != tt.wantErr {
				t.Errorf("ChangeTaskv2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTasksv2(t *testing.T) {
	t.Run("create tasks", TestCreateTaskv2)
	type args struct {
		ctx               context.Context
		offset            int
		limit             int
		presearchtaskname string
		createUID         string
	}
	tests := []struct {
		name    string
		args    args
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get tasks",
			args: args{
				ctx:    context.Background(),
				offset: 0,
				limit:  10,
			},
			want1:   1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1, err := GetTasksv2(tt.args.ctx, tt.args.offset, tt.args.limit, tt.args.presearchtaskname, tt.args.createUID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTasksv2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got1 != tt.want1 {
				t.Errorf("GetTasksv2() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestDeleteTaskv2(t *testing.T) {
	t.Run("create tasks", TestCreateTaskv2)
	tasks, count, err := GetTasksv2(context.Background(), 0, 10, "", "")
	if err != nil {
		t.Fatal("get task failed ", err)
	}
	if count == 0 {
		t.Fatal("count not should 0")
	}
	task := tasks[0]
	type args struct {
		ctx    context.Context
		taskid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "delete task",
			args: args{
				ctx:    context.Background(),
				taskid: task.ID,
			},
			wantErr: false,
		},
		{
			name: "delete not exist task",
			args: args{
				ctx:    context.Background(),
				taskid: "not exist id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteTaskv2(tt.args.ctx, tt.args.taskid); (err != nil) != tt.wantErr {
				t.Errorf("DeleteTaskv2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTaskIsUsev2(t *testing.T) {
	t.Run("create tasks", TestCreateTaskv2)
	tasks, count, err := GetTasksv2(context.Background(), 0, 10, "", "")
	if err != nil {
		t.Fatal("get task failed ", err)
	}
	if count == 0 {
		t.Fatal("count not should 0")
	}
	task := tasks[0]
	id := task.ID

	newtask := task
	newtask.Name = "newtask"
	newtask.ParentTaskIDs = IDs{task.ID}
	type args struct {
		ctx    context.Context
		taskid string
	}
	_, err = CreateTaskv2(context.Background(), newtask)
	if err != nil {
		t.Fatal("create new task failed ", err)
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "task is use",
			args: args{
				ctx:    context.Background(),
				taskid: id,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "task is not use",
			args: args{
				ctx:    context.Background(),
				taskid: "not_exist_id",
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TaskIsUsev2(gormdb, tt.args.taskid)
			if (err != nil) != tt.wantErr {
				t.Errorf("TaskIsUsev2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("TaskIsUsev2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTaskByIDv2(t *testing.T) {
	t.Run("create tasks", TestCreateTaskv2)
	tasks, count, err := GetTasksv2(context.Background(), 0, 10, "", "")
	if err != nil {
		t.Fatal("get task failed ", err)
	}
	if count == 0 {
		t.Fatal("count not should 0")
	}
	task := tasks[0]
	type args struct {
		ctx    context.Context
		taskid string
	}
	tests := []struct {
		name string
		args args
		// want    *Task
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get task",
			args: args{
				ctx:    context.Background(),
				taskid: task.ID,
			},
			wantErr: false,
		},
		{
			name: "get not exist task",
			args: args{
				ctx:    context.Background(),
				taskid: "not exist",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetTaskByIDv2(tt.args.ctx, tt.args.taskid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTaskByIDv2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetTaskByNamev2(t *testing.T) {
	t.Run("create tasks", TestCreateTaskv2)
	tasks, count, err := GetTasksv2(context.Background(), 0, 10, "", "")
	if err != nil {
		t.Fatal("get task failed ", err)
	}
	if count == 0 {
		t.Fatal("count not should 0")
	}
	task := tasks[0]

	type args struct {
		ctx      context.Context
		taskname string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get task",
			args: args{
				ctx:      context.Background(),
				taskname: task.Name,
			},
			wantErr: false,
		},
		{
			name: "get not exist task",
			args: args{
				ctx:      context.Background(),
				taskname: "not exist",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetTaskByNamev2(tt.args.ctx, tt.args.taskname)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTaskByNamev2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
