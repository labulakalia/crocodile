package model

import (
	"context"
	"testing"
	"time"

	"github.com/labulaka521/crocodile/core/utils/define"
)

func TestCreateLogv2(t *testing.T) {
	gormdb = InitGormSqlite()
	type args struct {
		ctx context.Context
		log *Log
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "create log",
			args: args{
				ctx: context.Background(),
				log: &Log{
					TaskName:    "test",
					TaskID:      "111111111111111221",
					StartTime:   time.Now(),
					EndTime:     time.Now(),
					Status:      1,
					TriggerType: define.Auto,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateLogv2(tt.args.ctx, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("CreateLogv2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetLogv2(t *testing.T) {
	t.Run("create lof", TestCreateLogv2)
	type args struct {
		ctx      context.Context
		taskname string
		status   int
		offset   int
		limit    int
	}
	tests := []struct {
		name    string
		args    args
		want    []*Log
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get all log ",
			args: args{
				ctx:      context.Background(),
				taskname: "",
				status:   1,
				offset:   0,
				limit:    10,
			},
			want1:   1,
			wantErr: false,
		},
		{
			name: "get log not task nam ",
			args: args{
				ctx:      context.Background(),
				taskname: "test",
				status:   1,
				offset:   0,
				limit:    10,
			},
			want1:   1,
			wantErr: false,
		},
		{
			name: "get log failed",
			args: args{
				ctx:      context.Background(),
				taskname: "test",
				status:   -1,
				offset:   0,
				limit:    10,
			},
			want1:   0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1, err := GetLogv2(tt.args.ctx, tt.args.taskname, tt.args.status, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLogv2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got1 != tt.want1 {
				t.Errorf("GetLogv2() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCleanLogv2(t *testing.T) {
	t.Run("create lof", TestCreateLogv2)
	logs, _, err := GetLogv2(context.Background(), "", 0, 0, 10)
	if err != nil {
		t.Fatalf("get log failed %v", err)
	}
	newlog := logs[0]
	type args struct {
		ctx             context.Context
		taskid          string
		beforeStartTime time.Time
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "clean log",
			args: args{
				ctx:             context.Background(),
				taskid:          newlog.TaskID,
				beforeStartTime: time.Now(),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CleanLogv2(tt.args.ctx, tt.args.taskid, tt.args.beforeStartTime); (err != nil) != tt.wantErr {
				t.Errorf("CleanLogv2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateNotify(t *testing.T) {
	gormdb = InitGormSqlite()
	type args struct {
		ctx    context.Context
		notify *Notify
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "create notify",
			args: args{
				ctx: context.Background(),
				notify: &Notify{
					Type:    define.TaskNotify,
					UID:     "11212121212121212121",
					Title:   "test",
					Content: "test content",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateNotify(tt.args.ctx, tt.args.notify); (err != nil) != tt.wantErr {
				t.Errorf("CreateNotify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetNotify(t *testing.T) {
	t.Run("create notify", TestCreateNotify)
	type args struct {
		ctx context.Context
		uid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get notify",
			args: args{
				ctx: context.Background(),
				uid: "11212121212121212121",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNotify(tt.args.ctx, tt.args.uid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNotify() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != 1 {
				t.Errorf("GetNotify() want count = %v, get %v", 1, len(got))
			}
		})
	}
}

func TestReadNotify(t *testing.T) {
	t.Run("create notify", TestCreateNotify)
	type args struct {
		ctx context.Context
		id  uint
		uid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "read notify",
			args: args{
				ctx: context.Background(),
				uid: "11212121212121212121",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadNotify(tt.args.ctx, tt.args.id, tt.args.uid); (err != nil) != tt.wantErr {
				t.Errorf("ReadNotify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
