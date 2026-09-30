package engine

import (
	"regexp"
	"strings"
)

// Classification never supplies execution arguments. Complex shell programs
// conservatively use the generic reducer; their source remains opaque.
var simpleShell = regexp.MustCompile(`^[a-zA-Z0-9_./ :\\-]+$`)

func shellClassification(source string) []string {
	if !simpleShell.MatchString(source) {
		return nil
	}
	return strings.Fields(source)
}

func mergeEnv(inherited []string, overrides map[string]string) []string {
	result := make([]string, 0, len(inherited)+len(overrides))
	for _, item := range inherited {
		key, _, _ := strings.Cut(item, "=")
		if _, replaced := overrides[key]; !replaced {
			result = append(result, item)
		}
	}
	for key, value := range overrides {
		result = append(result, key+"="+value)
	}
	return result
}
