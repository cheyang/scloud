package util

import (
	"fmt"
	"strconv"
)

// Support hostname extention from 001-fff(4095)

const (
	MaxValueForExtention int64 = 4096
)

func GetCurrentMaxExt(nums []string) string {

	maxValue int64 = 0

	vMap := make(map[int64]string, 0)

	for _, num := range nums {
		key, err := strconv.ParseInt(num, 16, 0)

		if err != nil {
			fmt.Println(err, "in GetCurrentMaxExt")
			continue
		}

		vMap[key] = num

		if key > maxValue {
			maxValue = key
		}

	}

	return nil
}
