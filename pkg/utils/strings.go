package utils

import (
	"crypto/rand"
	"encoding/base64"
	"regexp"
	"strings"
	"unicode"
)

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

// IsStringSliceCaseInsensitiveContains -- check slice contain string
func IsStringSliceCaseInsensitiveContains(stringSlice []string, searchString string) bool {
	for _, value := range stringSlice {
		if strings.EqualFold(value, searchString) {
			return true
		}
	}
	return false
}

func IsStringContainsAnyKeywords(s string, keywords []string) bool {
	contain := false
	for i := range keywords {
		if strings.Contains(s, keywords[i]) {
			contain = true
			break
		}
	}
	return contain
}

// RemoveStringSliceContains -- check slice contain string and remove
func RemoveStringSliceContains(stringSlice []string, searchString string) []string {
	newData := []string{}
	for _, value := range stringSlice {
		if value != searchString {
			newData = append(newData, value)
		}
	}
	return newData
}

// StringSlice -- slice string by separate
func StringSlice(s, sep string) []string {
	var sl []string

	for _, p := range strings.Split(s, sep) {
		if str := strings.TrimSpace(p); len(str) > 0 {
			sl = append(sl, strings.TrimSpace(p))
		}
	}

	return sl
}

// ToSnake Converts a string to snake_case
func ToSnake(s string) string {
	return ToDelimited(s, '_')
}

// ToScreamingSnake Converts a string to SCREAMING_SNAKE_CASE
func ToScreamingSnake(s string) string {
	return ToScreamingDelimited(s, '_', true)
}

// ToKebab Converts a string to kebab-case
func ToKebab(s string) string {
	return ToDelimited(s, '-')
}

// ToScreamingKebab Converts a string to SCREAMING-KEBAB-CASE
func ToScreamingKebab(s string) string {
	return ToScreamingDelimited(s, '-', true)
}

// ToDelimited Converts a string to delimited.snake.case (in this case `del = '.'`)
func ToDelimited(s string, del uint8) string {
	return ToScreamingDelimited(s, del, false)
}

// ToScreamingDelimited Converts a string to SCREAMING.DELIMITED.SNAKE.CASE (in this case `del = '.'; screaming = true`) or delimited.snake.case (in this case `del = '.'; screaming = false`)
func ToScreamingDelimited(s string, del uint8, screaming bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	for i, v := range s {
		// treat acronyms as words, eg for JSONData -> JSON is a whole word
		nextCaseIsChanged := false
		if i+1 < len(s) {
			next := rune(s[i+1])
			nextCaseIsChanged = (isRuneInRange(v, 'A', 'Z') && isRuneInRange(next, 'a', 'z')) ||
				(isRuneInRange(v, 'a', 'z') && isRuneInRange(next, 'A', 'Z'))
		}

		if i > 0 && n[len(n)-1] != del && nextCaseIsChanged {
			// add underscore if next letter case type is changed
			if isRuneInRange(v, 'A', 'Z') {
				n += string(del) + string(v)
			}

			if isRuneInRange(v, 'a', 'z') {
				n += string(v) + string(del)
			}
			continue
		}

		if v == ' ' || v == '_' || v == '-' {
			// replace spaces/underscores with delimiters
			n += string(del)
			continue
		}

		n = n + string(v)
	}

	n = strings.ToLower(n)
	if screaming {
		n = strings.ToUpper(n)
	}

	return n
}

func isRuneInRange(chr rune, fromRune rune, toRune rune) bool {
	return chr >= fromRune && chr <= toRune
}

// UcFirst Upper case first character
func UcFirst(str string) string {
	r := []rune(str)
	return string(unicode.ToUpper(r[0])) + string(r[1:])
}

// LcFirst Lower case first char
func LcFirst(str string) string {
	r := []rune(str)
	return string(unicode.ToLower(r[0])) + string(r[1:])
}

// Converts a string to CamelCase
func toCamelInitCase(s string, initCase bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	capNext := initCase
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			n += string(v)
		}
		if v >= '0' && v <= '9' {
			n += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' {
			capNext = true
		} else {
			capNext = false
		}
	}
	return n
}

// ToCamel Converts a string to CamelCase
func ToCamel(s string) string {
	return toCamelInitCase(s, true)
}

// ToLowerCamel Converts a string to lowerCamelCase
func ToLowerCamel(s string) string {
	if s == "" {
		return s
	}
	if r := rune(s[0]); r >= 'A' && r <= 'Z' {
		s = strings.ToLower(string(r)) + s[1:]
	}
	return toCamelInitCase(s, false)
}

// ToCamelInitCaseKeepAll Converts a string to CamelCase
func ToCamelInitCaseKeepAll(s string, initCase bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	capNext := initCase
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			n += string(v)
		}
		if v >= '0' && v <= '9' {
			n += string(v)
		}

		if v == '_' || v == ' ' || v == '-' {
			capNext = true
		} else {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
			capNext = false
		}
	}
	return n
}

var numberSequence = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
var numberReplacement = []byte(`$1 $2 $3`)

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)
	return string(b)
}
