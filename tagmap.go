package techlemm

import (
	"strings"
)

// TagMap is the main structure for looking up canonical tags
type TagMap struct {
	values map[string]string
}

// NewTagMap creates a new, empty TagMap for the purpose of looking up canonical tags
func NewTagMap(tags []string) *TagMap {
	result := &TagMap{
		values: make(map[string]string),
	}
	for _, tag := range tags {
		key := normalize(tag)
		result.values[key] = tag
	}
	return result
}

// normalize returns a string suitable as a key for tag lookup, removing dots and dashes and converting to lowercase
func normalize(s string) string {
	result := make([]rune, 0)

	for index, value := range s {
		if index == 0 {
			// Leading dots are meaningful and should not be removed, for example ".net"
			result = append(result, value)
			continue
		}
		if value == '.' || value == '-' {
			continue
		}
		result = append(result, value)
	}
	return strings.ToLower(string(result))
}
