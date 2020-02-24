package utils

import (
	"time"
)

// UnixToStr change int to strc
func UnixToStr(i int64) string {
	// return time.Unix(i,0).Location().String()
	return time.Unix(i, 0).Local().Format("2006-01-02 15:04:05")

}

// StrToUnix change str to unix
func StrToUnix(t string) int64 {
	tparse, err := time.Parse("2006-01-02 15:04:05", t)
	if err != nil {
		tparse, err = time.Parse("2006-01-02T15:04:05Z", t)
		if err != nil {
			return 0
		}
	}
	return tparse.Unix()
}
