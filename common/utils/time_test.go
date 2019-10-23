package utils

import (
	"testing"
	"time"
)

func TestUnixToTime(t *testing.T) {
	now := time.Now().Unix()

	nowstr := UnixToStr(now)

	t.Logf("%d change to %s", now, nowstr)
}

func TestTimeToUnix(t *testing.T) {
	now := time.Now().Unix()

	nowstr := UnixToStr(now)

	t.Logf("%d change to %s", now, nowstr)

	nowunix := StrToUnix(nowstr)
	if nowunix == 0 {
		t.Errorf("Change Failed: %s", nowstr)
	}
	t.Logf("%s change to %d", nowstr, now)
}
