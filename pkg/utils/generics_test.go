package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsTSliceContains(t *testing.T) {
	intsl := []int{1, 2, 3, 4, 5, 6}
	assert.True(t, IsSliceContains(intsl, 1))

	strsl := []string{"a", "b", "c"}
	assert.True(t, IsSliceContains(strsl, "c"))

}
