package utils

import "time"

func UnixToStr(i int64) string {
	return time.Unix(i, 0).Format("2006-01-02 03:04:05")
}

func StrToUnix(t string) int64 {
	tparse, err := time.Parse("2006-01-02 03:04:05", t)
	if err != nil {
		return 0
	}
	return tparse.Unix()
}
