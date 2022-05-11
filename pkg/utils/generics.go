package utils

// var (
// 	IsStringSliceContains = IsSliceContains[string]
// 	IsIntSliceContains    = IsSliceContains[int]
// 	IsInt64SliceContains  = IsSliceContains[int64]
// )

// func IsSliceContains[T constraints.Ordered](itemSlice []T, searchItem T) bool {
// 	for _, value := range itemSlice {
// 		if value == searchItem {
// 			return true
// 		}
// 	}
// 	return false
// }

func IsStringSliceContains(itemSlice []string, searchItem string) bool {
	for _, value := range itemSlice {
		if value == searchItem {
			return true
		}
	}
	return false
}

func IsIntSliceContains(itemSlice []int, searchItem int) bool {
	for _, value := range itemSlice {
		if value == searchItem {
			return true
		}
	}
	return false
}

func IsInt64SliceContains(itemSlice []int64, searchItem int64) bool {
	for _, value := range itemSlice {
		if value == searchItem {
			return true
		}
	}
	return false
}
