package utils

import "time"

var almatyLoc *time.Location

func init() {
	var err error
	almatyLoc, err = time.LoadLocation("Asia/Almaty")
	if err != nil {
		panic("failed to load almaty time")
	}
}

func Now() time.Time {
	return time.Now().In(almatyLoc)
}
