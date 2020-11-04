package model

import (
	"context"
	"testing"
)

func TestGetIDNameDIct(t *testing.T) {
	t.Run("create hostGroup", TestCreateHostgroupv2)
	hgs, count, err := GetHostGroupsv2(context.Background(), 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("can get hg")
	}
	hg := hgs[0]

	type args struct {
		ctx   context.Context
		model interface{}
		ids   []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "id name",
			args: args{
				ctx:   context.Background(),
				model: &HostGroup{},
				ids: []string{
					hg.ID,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIDNameDict(tt.args.ctx, tt.args.ids, tt.args.model)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetIDNameDIct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%+v", got)
		})
	}
}
