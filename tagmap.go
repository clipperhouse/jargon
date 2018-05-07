package techlemm

import (
	"fmt"
	"strings"
)

// TagMap is the main structure for looking up canonical tags
type TagMap struct {
	values map[string]string
}

// NewTagMap creates a new, empty TagMap for the purpose of looking up canonical tags
func NewTagMap() *TagMap {
	return &TagMap{
		values: make(map[string]string),
	}
}

func (t *TagMap) AddTag(s string) error {
	key := normalize(s)
	t.values[key] = s
	return nil
}

// normalize returns a string suitable as a key for tag lookup, removing dots and dashes and converting to lowercase
func normalize(s string) string {
	result := make([]rune, 0)

	for index, value := range s {
		fmt.Printf("index: %q, value: %q\n", index, value)
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
	fmt.Printf("%v\n", result)
	fmt.Printf("%q\n", string(result))
	return strings.ToLower(string(result))
}
