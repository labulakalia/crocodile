package alarm

import "github.com/labulaka521/crocodile/core/utils/define"

// CheckAlarm will if err != nil will rediret return err
// othersie will check code and task output
func CheckAlarm(id string, runbyid string, tasktype define.TaskRespType, code int, output []byte, err error) error {
	return nil
}