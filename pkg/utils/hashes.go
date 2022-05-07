package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

// GetMD5Hash -- get md5 hash from a string
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetSHA256Hash -- get sha_256 hash from a string
func GetSHA256Hash(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

// GetSHA512Hash -- get sha_512 hash from a string
func GetSHA512Hash(text string) string {
	shaHash := sha512.New()
	shaHash.Write([]byte(text))
	return hex.EncodeToString(shaHash.Sum(nil))
}
