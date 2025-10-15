package common

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var TemplateRegex = regexp.MustCompile(`{{\s*([a-zA-Z0-9_\.]+)\s*}}`)

func GetNestedField(data map[string]any, path string) (any, error) {
	parts := strings.Split(path, ".")
	var current any = data
	for _, part := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("cannot access field %q in non-map type", part)
		}
		val, exists := m[part]
		if !exists {
			return nil, fmt.Errorf("field %q not found", part)
		}
		current = val
	}
	return current, nil
}

func ProcessTemplate(msg string, attrs map[string]any) string {
	return TemplateRegex.ReplaceAllStringFunc(msg, func(s string) string {
		matches := TemplateRegex.FindStringSubmatch(s)
		if len(matches) < 2 {
			return s
		}
		fieldPath := matches[1]

		val, err := GetNestedField(attrs, fieldPath)
		if err != nil {
			return fmt.Sprintf("{{%s:ERR}}", fieldPath)
		}

		switch v := val.(type) {
		case string:
			return v
		case int:
			return strconv.Itoa(v)
		case float64:
			return strconv.FormatFloat(v, 'f', -1, 64)
		case bool:
			return strconv.FormatBool(v)
		default:
			return fmt.Sprintf("%v", v)
		}
	})
}
