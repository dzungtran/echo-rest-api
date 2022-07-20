package utils

import (
	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
)

func GenerateUUID() string {
	return shortuuid.New()
}

func GenerateLongUUID() string {
	return uuid.New().String()
}
