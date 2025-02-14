package util

import (
	"runtime"
)

func GetOsLineSep() string {
	switch runtime.GOOS {
	case "windows":
		return "\r\n"
	case "linux":
		return "\n"
	case "darwin":
		return "\n"
	default:
		return "\n"
	}
}
