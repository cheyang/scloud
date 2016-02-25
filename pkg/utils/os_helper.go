package utils

import (
	"os"
	"runtime"
)

func GetHomedir() {

	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}
