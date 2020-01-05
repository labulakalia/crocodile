package tasktype

import (
	"context"
	pb "github.com/labulaka521/crocodile/core/proto"
	"os"
	"os/exec"
)

// DataShell task run shell 
type DataShell struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

// Run implment TaskRuner
// run shell command
// do not return error
func (ds *DataShell) Run(ctx context.Context) (resp *pb.TaskResp) {
	shell := os.Getenv("SHELL")
	resp = &pb.TaskResp{}
	cmd := exec.CommandContext(ctx, shell, "-c", ds.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		resp.Code = -1
		resp.ErrMsg = []byte(err.Error())
		return resp
	}
	resp.RespData = output
	resp.Code = 0
	return resp
}
