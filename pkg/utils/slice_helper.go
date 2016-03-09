package utils

// check if the element in the slice
func Contains(list []interface{}, elem interface{}) bool {
	for _, t := range list {
		if t == elem {
			return true
		}
	}
	return false
}
