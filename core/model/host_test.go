package model

import (
	"reflect"
	"testing"
)

func Test_deletefromslice(t *testing.T) {
	type args struct {
		deleteid string
		ids      []string
	}
	tests := []struct {
		name  string
		args  args
		want  []string
		want1 bool
	}{
		// TODO: Add test cases.
		struct{name string; args args; want []string; want1 bool}{
			name: "case1",
			args: args{
				deleteid: "3",
				ids: []string{"1","2","3","4","5"},
			},
			want: []string{"1","2","4","5"},
			want1: true,
		},
		{
			name: "case2",
			args: args{
				deleteid: "6",
				ids: []string{"1","2","3","4","5"},
			},
			want: []string{"1","2","3","4","5"},
			want1: false,
		},
		{
			name: "case3",
			args: args{
				deleteid: "6",
				ids: []string{"6"},
			},
			want: []string{},
			want1: true,
		},
		{
			name: "case4",
			args: args{
				deleteid: "6",
				ids: []string{},
			},
			want: []string{},
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := deletefromslice(tt.args.deleteid, tt.args.ids)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deletefromslice() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("deletefromslice() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
