package util

import "time"

func SleepMS(s int) {
	time.Sleep(time.Duration(s) * time.Millisecond)
}

func SleepS(s int) {
	time.Sleep(time.Duration(s) * time.Second)
}
