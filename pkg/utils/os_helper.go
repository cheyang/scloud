package utils

import (
	"os"
	"runtime"
)

func GetHomedir() string {

	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}
