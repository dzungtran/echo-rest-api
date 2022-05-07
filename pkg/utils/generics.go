package utils

import (
	"golang.org/x/exp/constraints"
)

var (
	IsStringSliceContains = IsSliceContains[string]
	IsIntSliceContains    = IsSliceContains[int]
	IsInt64SliceContains  = IsSliceContains[int64]
)

func IsSliceContains[T constraints.Ordered](itemSlice []T, searchItem T) bool {
	for _, value := range itemSlice {
		if value == searchItem {
			return true
		}
	}
	return false
}
