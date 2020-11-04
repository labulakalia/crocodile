package model

import (
	"context"
	"testing"

	pb "github.com/labulaka521/crocodile/core/proto"
)

func TestCreateHostgroupv2(t *testing.T) {
	gormdb = InitGormSqlite()
	type args struct {
		ctx     context.Context
		name    string
		remark  string
		userid  string
		hostids []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "create new hostgroup",
			args: args{
				ctx:     context.Background(),
				name:    "hg1",
				remark:  "remark1",
				userid:  "uid",
				hostids: []string{"hid127272222222222", "hid127272222222221"},
			},
			wantErr: false,
		},
		{
			name: "create exist new hostgroup",
			args: args{
				ctx:     context.Background(),
				name:    "hg1",
				remark:  "remark1",
				userid:  "uid",
				hostids: []string{"hid127272222222222", "hid127272222222221"},
			},
			wantErr: true,
		},
		{
			name: "create not valid hostid",
			args: args{
				ctx:     context.Background(),
				name:    "hg1",
				remark:  "remark1",
				userid:  "uid",
				hostids: []string{"valid hostid", "valid hostid2"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateHostgroupv2(tt.args.ctx, tt.args.name, tt.args.remark, tt.args.userid, tt.args.hostids); (err != nil) != tt.wantErr {
				t.Errorf("CreateHostgroupv2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestChangeHostGroupv2(t *testing.T) {
	t.Run("create hostgroup", TestCreateHostgroupv2)
	hgs, count, err := GetHostGroupsv2(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("get hostgroup failed %v", err)
	}
	if count == 0 {
		t.Fatalf("can find hostgroup")
	}
	hg := hgs[0]
	type args struct {
		ctx        context.Context
		hostids    []string
		id         string
		currentUID string
		remark     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "change hostgroup",
			args: args{
				ctx:     context.Background(),
				hostids: []string{"hid127272222222222", "hid127272222222221"},
				id:      hg.ID,
				remark:  "remark2",
			},
			wantErr: false,
		},
		{
			name: "change not exist hostgroup",
			args: args{
				ctx:     context.Background(),
				hostids: []string{"hid127272222222233", "hid127272222222234"},
				id:      "not exist hg id",
				remark:  "remark2",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ChangeHostGroupv2(tt.args.ctx, tt.args.hostids, tt.args.id, tt.args.remark, tt.args.currentUID); (err != nil) != tt.wantErr {
				t.Errorf("ChangeHostGroupv2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteHostGroupv2(t *testing.T) {
	t.Run("create hostgroup", TestCreateHostgroupv2)
	hgs, count, err := GetHostGroupsv2(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("get hostgroup failed %v", err)
	}
	if count == 0 {
		t.Fatalf("can find hostgroup")
	}
	hg := hgs[0]
	type args struct {
		ctx        context.Context
		id         string
		currentUID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "delete exist hostgroup",
			args: args{
				ctx: context.Background(),
				id:  hg.ID,
			},
			wantErr: false,
		},
		{
			name: "delete not exist hostgroup",
			args: args{
				ctx: context.Background(),
				id:  "not exist id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteHostGroupv2(tt.args.ctx, tt.args.id, tt.args.currentUID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteHostGroupv2() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetHostGroupsv2(t *testing.T) {
	t.Run("create hostgroup", TestCreateHostgroupv2)
	type args struct {
		ctx    context.Context
		limit  int
		offset int
	}
	tests := []struct {
		name string
		args args
		// want    []*HostGroup
		want1   int64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get all hostgroup",
			args: args{
				ctx:    context.Background(),
				limit:  10,
				offset: 0,
			},
			want1:   1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1, err := GetHostGroupsv2(tt.args.ctx, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostGroupsv2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got1 != tt.want1 {
				t.Errorf("GetHostGroupsv2() got1 = %v, want %v", got1, tt.want1)
			}

		})
	}
}

func TestGetHostGroupByIDv2(t *testing.T) {
	t.Run("create hostgroup", TestCreateHostgroupv2)
	hgs, count, err := GetHostGroupsv2(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("get hostgroup failed %v", err)
	}
	if count == 0 {
		t.Fatalf("can find hostgroup")
	}
	hg := hgs[0]
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
			name: "get exist hg id",
			args: args{
				ctx: context.Background(),
				id:  hg.ID,
			},
			wantErr: false,
		},
		{
			name: "get not exist hg id",
			args: args{
				ctx: context.Background(),
				id:  "not_exist_id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHostGroupByIDv2(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostGroupByIDv2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if got.Name != hg.Name {
					t.Errorf("get hostgroup name error = %v, want %v", got.Name, hg.Name)
				}
			}
		})
	}
}

func TestGetHostGroupByNamev2(t *testing.T) {
	t.Run("create hostgroup", TestCreateHostgroupv2)
	hgs, count, err := GetHostGroupsv2(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("get hostgroup failed %v", err)
	}
	if count == 0 {
		t.Fatalf("can find hostgroup")
	}
	hg := hgs[0]

	type args struct {
		ctx    context.Context
		hgname string
	}
	tests := []struct {
		name    string
		args    args
		want    *HostGroup
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get exist hg id",
			args: args{
				ctx:    context.Background(),
				hgname: hg.Name,
			},
			wantErr: false,
		},
		{
			name: "get not exist hg id",
			args: args{
				ctx:    context.Background(),
				hgname: "not_exist_id",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHostGroupByNamev2(tt.args.ctx, tt.args.hgname)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostGroupByNamev2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if got.ID != hg.ID {
					t.Errorf("get hostgroup id error = %v, want %v", got.ID, hg.ID)
				}
			}
		})
	}
}

func TestGetHostsByHGIDv2(t *testing.T) {
	t.Run("create hostgroup", TestCreateHostgroupv2)

	hgs, count, err := GetHostGroupsv2(context.Background(), 0, 10)
	if err != nil {
		t.Fatalf("get hostgroup failed %v", err)
	}
	if count == 0 {
		t.Fatalf("can find hostgroup")
	}
	hg := hgs[0]

	reghost := &pb.RegistryReq{
		Ip:        "1.1.1.1",
		Port:      80,
		Weight:    100,
		Hostname:  "tetst",
		Hostgroup: "hg1",
	}
	err = CreateOrUpdateHost(context.Background(), reghost)
	if err != nil {
		t.Fatalf("create host failed %v", err)
	}
	type args struct {
		ctx  context.Context
		hgid string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "get host by hg id",
			args: args{
				ctx:  context.Background(),
				hgid: hg.ID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetHostsByHGIDv2(tt.args.ctx, tt.args.hgid)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetHostsByHGIDv2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Fatal("can find hosts by hg id")
			}
			if got[0].Addr != got[0].Addr {
				t.Fatalf("get hosts failed want addr %v, get addr %v", got[0].Addr, got[0].Addr)
			}
		})
	}
}
