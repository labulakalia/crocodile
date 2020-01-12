package tasktype

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"github.com/labulaka521/crocodile/core/utils/resp"
)

var _ TaskRuner = &DataShell{}

// DataShell task run shell
type DataShell struct {
	Command string `json:"command"`
}


// Run implment TaskRuner
// run shell command
// return io.ReadCloser
func (ds *DataShell)Run(ctx context.Context) io.ReadCloser {
	pr, pw := io.Pipe()
	go func() {
		var exitCode = DefaultExitCode
		defer pw.Close()
		defer func() {
			pw.Write([]byte(fmt.Sprintf("%3d", exitCode))) // write exitCode,total 3 byte
		}()
		// tell the command to write to our pipe
		shell := os.Getenv("SHELL")

		cmd := exec.CommandContext(ctx, shell, "-c", ds.Command)
		cmd.Stdout = pw
		cmd.Stderr = pw
		err := cmd.Start()
		if err != nil {
			pw.Write([]byte(err.Error()))
			return
		}

		err = cmd.Wait()
		if err != nil {
			var customerr bytes.Buffer
			// deal err
			// if context err,will change err to custom msg
			switch ctx.Err() {
			case context.DeadlineExceeded:
				customerr.WriteString(resp.GetMsg(resp.ErrCtxDeadlineExceeded))
			case context.Canceled:
				customerr.WriteString(resp.GetMsg(resp.ErrCtxCanceled))
			default:
				customerr.WriteString(err.Error())
			}
			pw.Write(customerr.Bytes())
			// try to get the exit code
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			}
		} else {
			exitCode = 0
		}

	}()
	return pr
}