package schedule

import (
	"testing"

	"github.com/labulaka521/crocodile/core/utils/define"
)

func Test_generatename(t *testing.T) {
	id := "233903600084979712"
	name := generatename(id, define.MasterTask)
	t.Log(name)

	getid, tasktype, err := splitname(name)
	if err != nil {
		t.Error(err)
	}
	if getid != id {
		t.Errorf("want get id %s, but get %s", id, getid)
	}
	if tasktype != define.MasterTask {
		t.Errorf("want get tasktype %d, but get %d", define.MasterTask, tasktype)
	}
}
