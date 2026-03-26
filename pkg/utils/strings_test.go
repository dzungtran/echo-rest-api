package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToSnake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello_world"},
		{"helloWorld", "hello_world"},
		{"Hello_World", "hello_world"},
		{"hello-world", "hello_world"},
		{"HelloWorld123", "hello_world_123"}, // Note: underscore before numbers
		{"XMLParser", "xml_parser"},
		{"simple", "simple"},
		{"", ""},
		{"AlreadySnake", "already_snake"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := ToSnake(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToKebab(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HelloWorld", "hello-world"},
		{"helloWorld", "hello-world"},
		{"Hello_World", "hello-world"},
		{"Hello-World", "hello-world"},
		{"HelloWorld123", "hello-world-123"}, // Note: hyphen before numbers
		{"XMLParser", "xml-parser"},
		{"simple", "simple"},
		{"", ""},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := ToKebab(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToCamel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"hello-world", "HelloWorld"},
		{"hello", "Hello"},
		{"alreadyCamel", "AlreadyCamel"},
		{"xml_parser", "XmlParser"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := ToCamel(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToLowerCamel(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "helloWorld"},
		{"hello-world", "helloWorld"},
		{"HelloWorld", "helloWorld"},
		{"hello", "hello"},
		{"Hello", "hello"},
		{"", ""},
		{"AlreadyCamel", "alreadyCamel"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := ToLowerCamel(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToScreamingSnake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HELLO_WORLD"},
		{"helloWorld", "HELLO_WORLD"},
		{"HelloWorld", "HELLO_WORLD"},
		{"hello-world", "HELLO_WORLD"},
		{"simple", "SIMPLE"},
		{"", ""},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := ToScreamingSnake(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestToScreamingKebab(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HELLO-WORLD"},
		{"helloWorld", "HELLO-WORLD"},
		{"HelloWorld", "HELLO-WORLD"},
		{"hello-kebab", "HELLO-KEBAB"},
		{"simple", "SIMPLE"},
		{"", ""},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := ToScreamingKebab(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestUcFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"Hello", "Hello"},
		{"h", "H"},
		{"123abc", "123abc"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := UcFirst(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestLcFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "hello"},
		{"hello", "hello"},
		{"H", "h"},
		{"123ABC", "123ABC"},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := LcFirst(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestStringSlice(t *testing.T) {
	tests := []struct {
		input    string
		sep      string
		expected []string
	}{
		{"a,b,c", ",", []string{"a", "b", "c"}},
		{"a, b, c", ",", []string{"a", "b", "c"}},
		{"a | b | c", " | ", []string{"a", "b", "c"}},
		{"single", ",", []string{"single"}},
		{"", ",", nil},
		{"a,,c", ",", []string{"a", "c"}},
		{"  a , b  , c  ", ",", []string{"a", "b", "c"}},
	}

	for _, tc := range tests {
		t.Run(tc.input, func(t *testing.T) {
			result := StringSlice(tc.input, tc.sep)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsStringSliceCaseInsensitiveContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		val      string
		expected bool
	}{
		{"found exact", []string{"Hello", "World"}, "Hello", true},
		{"found case insensitive", []string{"Hello", "World"}, "hello", true},
		{"found mixed case", []string{"Hello", "World"}, "WORLD", true},
		{"not found", []string{"Hello", "World"}, "foo", false},
		{"empty slice", []string{}, "hello", false},
		{"empty val", []string{"Hello"}, "", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsStringSliceCaseInsensitiveContains(tc.slice, tc.val)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestIsStringContainsAnyKeywords(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		keywords []string
		expected bool
	}{
		{"contains one", "hello world", []string{"foo", "bar", "world"}, true},
		{"contains none", "hello world", []string{"foo", "bar", "baz"}, false},
		{"empty string", "", []string{"foo"}, false},
		{"empty keywords", "hello world", []string{}, false},
		{"exact match", "error", []string{"error"}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsStringContainsAnyKeywords(tc.s, tc.keywords)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestRemoveStringSliceContains(t *testing.T) {
	tests := []struct {
		name     string
		slice    []string
		val      string
		expected []string
	}{
		{"remove one", []string{"a", "b", "c"}, "b", []string{"a", "c"}},
		{"remove first", []string{"a", "b", "c"}, "a", []string{"b", "c"}},
		{"remove last", []string{"a", "b", "c"}, "c", []string{"a", "b"}},
		{"remove not found", []string{"a", "b", "c"}, "d", []string{"a", "b", "c"}},
		{"remove all same", []string{"a", "a", "a"}, "a", []string{}},
		{"empty slice", []string{}, "a", []string{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := RemoveStringSliceContains(tc.slice, tc.val)
			assert.Equal(t, tc.expected, result)
		})
	}
}
