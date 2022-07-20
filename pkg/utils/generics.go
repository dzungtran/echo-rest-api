package utils

import (
	"golang.org/x/exp/constraints"
)

func IsSliceContains[T constraints.Ordered](itemSlice []T, searchItem T) bool {
	for _, val := range itemSlice {
		if val == searchItem {
			return true
		}
	}
	return false
}
