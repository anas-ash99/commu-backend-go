package utils

import "time"

func GetCurrentTime() string {
	localTime := time.Now()
	formattedTime := localTime.Format("2006-01-02 15:04:05 -0700 MST")
	return formattedTime
}
