package utils

import "github.com/lithammer/shortuuid/v4"

func GenerateUUID() string {
	return shortuuid.New()
}
