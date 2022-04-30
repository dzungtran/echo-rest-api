package utils

// IsStringSliceContains -- check slice contain string
// func Is[T any]SliceContains(itemSlice []T, searchItem T) bool {
// 	for _, value := range itemSlice {
// 		if value == searchItem {
// 			return true
// 		}
// 	}
// 	return false
// }

// IsStringSliceContains -- check slice contain string
func IsStringSliceContains(stringSlice []string, searchString string) bool {
	for _, value := range stringSlice {
		if value == searchString {
			return true
		}
	}
	return false
}
