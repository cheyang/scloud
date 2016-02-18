package util

import (
	"fmt"
	"strconv"
)

// Support hostname extention from 001-fff(4095)

const (
	MaxValueForExtention = 4096
)

func GetCurrentMaxExt(nums []string) string {

	maxValue := 0

	vMap := make(map[int]string, 0)

	for _, num := range nums {
		key, err := strconv.ParseInt(num, 16, 0)

		if err != nil {
			fmt.Println(error, "in GetCurrentMaxExt")
			continue
		}

		vMap[key] = num

		if key > maxValue {
			maxValue = key
		}

	}

	return nil
}
