package properties

import (
	"fmt"
	"strings"
)

// Map is a map of server properties.
type Map map[string]interface{}

// Read reads the bytes of a server.properties file into a `map`.
//
// An error is returned if any lines are malformed (i.e., cannot
// be split into a key/value pair).
func Read(raw []byte) (Map, error) {
	ret := make(Map, 0)

	lines := strings.Split(string(raw), "\n")
	for i, line := range lines {
		// Skip lines that are comments or blank
		if strings.HasPrefix(line, "#") || len(line) == 0 {
			continue
		}

		// Split the line into a key/value pair
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			return nil, fmt.Errorf("malformed line: cannot split key and value:\n\tLine %d: '%s'", i, line)
		}

		// Shove em in the map
		ret[parts[0]] = parts[1]
	}

	return ret, nil
}
