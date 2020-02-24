package utils

import "testing"

func TestCheckEmail(t *testing.T) {
	type args struct {
		email string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "email1",
			args:    args{email: "labulaka522@163.com"},
			wantErr: false,
		},
		{
			name:    "email2",
			args:    args{email: "labulaka86218^>[1~@163.com"},
			wantErr: true,
		},
		{
			name:    "email3",
			args:    args{email: "idowe*skwqun ejijiji"},
			wantErr: true,
		},
		{
			name:    "email4",
			args:    args{email: "labulaka@123.com"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CheckEmail(tt.args.email); (err != nil) != tt.wantErr {
				t.Errorf("CheckEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
