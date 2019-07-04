package util

import (
	"github.com/google/uuid"
)

// 生成uuid

func GenerateID() (uid string) {

	uid = uuid.New().String()
	return
}
