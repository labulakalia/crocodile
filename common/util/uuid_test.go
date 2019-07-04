package util

import "testing"

func TestGenerateID(t *testing.T) {
	for i := 0; i < 100; i++ {
		uuid := GenerateID()
		t.Log(uuid)
	}
}
