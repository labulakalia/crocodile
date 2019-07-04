package util

import (
	"testing"
)

var (
	password = "testpassword"
)

func TestCheckgeneratepass(t *testing.T) {
	var (
		err          error
		hashpassword string
	)
	if hashpassword, err = GenerateHashPass(password); err != nil {
		t.Errorf("GenerateHashPass Err :%v\n", err)
	}

	if err = CheckHashPass(hashpassword, password); err != nil {
		t.Errorf("Check Pass Err HashPass: %s Password %s", hashpassword, password)
	}
}
