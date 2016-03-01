package utils

import (
	"fmt"
	"strconv"
)

// Support hostname extention from 001-fff(4095)

const (
	MaxValueForExtention int64 = 4096
)

func GetCurrentMaxExt(nums []string) string {

	var maxValue int64 = 0

	vMap := make(map[int64]string, 0)

	for _, num := range nums {
		key, err := strconv.ParseInt(num, 16, 0)

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s in GetCurrentMaxExt", err)
			continue
		}

		vMap[key] = num

		if key > maxValue {
			maxValue = key
		}

	}

	return vMap[maxValue]
}
