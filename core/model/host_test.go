package model

import (
	"context"
	"reflect"
	"testing"

	pb "github.com/labulaka521/crocodile/core/proto"
)

func TestCreateOrUpdateHost(t *testing.T) {
	gormdb = InitGormSqlite()
	type args struct {
		ctx context.Context
		req *pb.RegistryReq
	}
	// TODO create hostgroup
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "create1",
			args: args{
				ctx: context.Background(),
				req: &pb.RegistryReq{
					Ip:        "192.168.1.1",
					Port:      8080,
					Weight:    100,
					Hostgroup: "hostname1",
					Version:   "1.9.1",
					Remark:    "remark1",
				},
			},
			wantErr: false,
		},
		{
			name: "update1",
			args: args{
				ctx: context.Background(),
				req: &pb.RegistryReq{
					Ip:        "192.168.1.1",
					Port:      8080,
					Weight:    90,
					Hostgroup: "hostname2",
					Version:   "1.9.2",
					Remark:    "remark2",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateOrUpdateHost(tt.args.ctx, tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("CreateOrUpdateHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateHostHearbeatv2(t *testing.T) {
	if !t.Run("create hosts", TestCreateOrUpdateHost) {
		t.Error("Run create host failed")
		return
	}
	type args struct {
		ctx           context.Context
		addr          string
		countRunTasks int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "update exist addr host hb",
			args: args{
				ctx:           context.TODO(),
				addr:          "192.168.1.1:8080",
				countRunTasks: 10,
			},
			wantErr: false,
		},
		{
			name: "update not exist addr host hb",
			args: args{
				ctx:           context.TODO(),
				addr:          "not_exist_addr",
				countRunTasks: 10,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateHostHearbeatv2(tt.args.ctx, tt.args.addr, tt.args.countRunTasks); (err != nil) != tt.wantErr {
				t.Errorf("UpdateHostHearbeatv2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetHostsv2(t *testing.T) {
	if !t.Run("create hosts", TestCreateOrUpdateHost) {
		t.Error("Run create host failed")
		return
	}
	type args struct {
		ctx    context.Context
		offset int
		limit  int
	}
	tests := []struct {
		name    string
		args    args
		want    []*Host
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get host",
			args: args{
				ctx:    context.TODO(),
				offset: 0,
				limit:  10,
			},
			want: []*Host{
				{
					Addr:    "192.168.1.1:8080",
					Weight:  90,
					Version: "1.9.2",
					Remark:  "remark2",
					Stop:    false,
				},
			},
			want1:   1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, count, err := GetHostsv2(tt.args.ctx, tt.args.offset, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostsv2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for i := range tt.want {
				gothost := &Host{
					Addr:    got[i].Addr,
					Weight:  got[i].Weight,
					Version: got[i].Version,
					Remark:  got[i].Remark,
					Stop:    got[i].Stop,
				}
				if !reflect.DeepEqual(gothost, tt.want[i]) {
					t.Errorf("GetHostsv2() got = %v, want %v", gothost, tt.want)
				}
			}

			if count != tt.want1 {
				t.Errorf("GetHostsv2() count = %v, want %v", count, tt.want1)
			}
		})
	}
}

func TestGetHostByIDv2(t *testing.T) {
	if !t.Run("create hosts", TestCreateOrUpdateHost) {
		t.Error("Run create host failed")
		return
	}
	hosts, count, err := GetHostsv2(context.Background(), 0, 10)
	if err != nil {
		t.Errorf("get all hosts failed %v", err)
	}
	if count == 0 {
		t.Log("can not get valid hosts")
		return
	}
	id := hosts[0].ID
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		want    *Host
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get host by exists id",
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			want: &Host{
				Addr:    "192.168.1.1:8080",
				Weight:  90,
				Version: "1.9.2",
				Remark:  "remark2",
				Stop:    false,
			},
			wantErr: false,
		},
		{
			name: "get host by not exists id",
			args: args{
				ctx: context.Background(),
				id:  "not exists host id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHostByIDv2(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostByIDv2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				gothost := &Host{
					Addr:    got.Addr,
					Weight:  got.Weight,
					Version: got.Version,
					Remark:  got.Remark,
					Stop:    got.Stop,
				}
				if !reflect.DeepEqual(gothost, tt.want) {
					t.Errorf("GetHostsv2() got = %v, want %v", gothost, tt.want)
				}
			}
		})
	}
}

func TestGetHostsByIDSv2(t *testing.T) {
	if !t.Run("create hosts", TestCreateOrUpdateHost) {
		t.Error("Run create host failed")
		return
	}
	hosts, count, err := GetHostsv2(context.Background(), 0, 10)
	if err != nil {
		t.Errorf("get all hosts failed %v", err)
	}
	if count == 0 {
		t.Log("can not get valid hosts")
		return
	}
	ids := []string{}
	for _, h := range hosts {
		ids = append(ids, h.ID)
	}
	type args struct {
		ctx context.Context
		ids []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Host
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get host by ids",
			args: args{
				ctx: context.TODO(),
				ids: ids,
			},
			want: []*Host{
				{
					Addr:    "192.168.1.1:8080",
					Weight:  90,
					Version: "1.9.2",
					Remark:  "remark2",
					Stop:    false,
				},
			},
			wantErr: false,
		},
		{
			name: "get host by not exist ids",
			args: args{
				ctx: context.TODO(),
				ids: []string{"not_exists_id"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHostsByIDSv2(tt.args.ctx, tt.args.ids)
			if len(got) != len(tt.args.ids) {
				t.Log("can find hosts ids ", tt.args.ids)
				if !tt.wantErr {
					t.Errorf("count got %d, want %d", len(got), len(tt.args.ids))
				}
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostsByIDSv2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for i := range got {
				gothost := &Host{
					Addr:    got[i].Addr,
					Weight:  got[i].Weight,
					Version: got[i].Version,
					Remark:  got[i].Remark,
					Stop:    got[i].Stop,
				}
				if !reflect.DeepEqual(gothost, tt.want[i]) {
					t.Errorf("GetHostsv2() got = %+v, want %+v", gothost, tt.want[i])
				}
			}
			if len(got) != int(count) {
				t.Errorf("GetHostsv2() count = %d, want %d", len(got), count)
			}
		})
	}
}

func TestChangeHostStopStatus(t *testing.T) {
	if !t.Run("create hosts", TestCreateOrUpdateHost) {
		t.Error("Run create host failed")
		return
	}
	hosts, count, err := GetHostsv2(context.Background(), 0, 10)
	if err != nil {
		t.Errorf("get hosts failed %v", err)
		return
	}
	if count == 0 {
		t.Log("can not get valid hosts")
		return
	}
	id := hosts[0].ID
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "change exist host status",
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			wantErr: false,
		},
		{
			name: "change not exist host status",
			args: args{
				ctx: context.Background(),
				id:  "not_exist_id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ChangeHostStopStatus(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("ChangeHostStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteHostv2(t *testing.T) {
	if !t.Run("create hosts", TestCreateOrUpdateHost) {
		t.Error("Run create host failed")
		return
	}

	hosts, count, err := GetHostsv2(context.Background(), 0, 10)
	if err != nil {
		t.Errorf("get hosts failed %v", err)
		return
	}
	if count == 0 {
		t.Log("can not get valid hosts")
		return
	}
	id := hosts[0].ID
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test delete host",
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			wantErr: false,
		},
		{
			name: "test delete not exist host",
			args: args{
				ctx: context.Background(),
				id:  "not exist error",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteHostv2(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteHostv2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
