package util

import (
	"math/rand"
	"time"
)

func SleepMS(s int, max int) {
	var t int
	if max > s {
		t = rand.Intn(max-s) + s
	} else {
		t = s
	}
	time.Sleep(time.Duration(t) * time.Millisecond)
}

func SleepS(s int) {
	time.Sleep(time.Duration(s) * time.Second)
}
