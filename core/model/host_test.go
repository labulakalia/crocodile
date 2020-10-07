package model

import (
	"context"
	"testing"

	pb "github.com/labulaka521/crocodile/core/proto"
)

func setup() {
	InitGormSqlite()
}

func TestCreateOrUpdateHost(t *testing.T) {
	type args struct {
		ctx context.Context
		req *pb.RegistryReq
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "creat1",
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
				wantErr: false,
			},
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
