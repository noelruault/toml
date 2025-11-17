/*
This TOML parser is a heavily customized, unstable, and workflow-specific implementation. It was written quickly to handle a narrow subset of TOML files and is not production-ready.

It works for simple `[section] key = value` config setups.
It does NOT support full TOML syntax, nested tables, arrays, escapes, or type safety.
*/
package toml

import (
	"strconv"
	"strings"
)

type TOMLData map[string]map[string]string

type Unmarshaler interface {
	UnmarshalTOML(data TOMLData)
}

func Parse(data []byte) (TOMLData, error) {
	return parse(string(data))
}

func UnmarshalTOML(data []byte, v Unmarshaler) error {
	parsed, err := Parse(data)
	if err != nil {
		return err
	}
	v.UnmarshalTOML(parsed)
	return nil
}

func (t TOMLData) GetString(section, key string) string {
	if sec, ok := t[strings.ToLower(section)]; ok {
		return unquote(sec[key])
	}
	return ""
}

func (t TOMLData) GetBool(section, key string) bool {
	if sec, ok := t[strings.ToLower(section)]; ok {
		if val, exists := sec[key]; exists {
			if b, err := strconv.ParseBool(strings.TrimSpace(val)); err == nil {
				return b
			}
		}
	}
	return false
}

func (t TOMLData) GetInt64(section, key string) int64 {
	if sec, ok := t[strings.ToLower(section)]; ok {
		if val, exists := sec[key]; exists {
			if i, err := strconv.ParseInt(strings.TrimSpace(val), 10, 64); err == nil {
				return i
			}
		}
	}
	return 0
}

func (t TOMLData) GetFloat64(section, key string) float64 {
	if sec, ok := t[strings.ToLower(section)]; ok {
		if val, exists := sec[key]; exists {
			if f, err := strconv.ParseFloat(strings.TrimSpace(val), 64); err == nil {
				return f
			}
		}
	}
	return 0.0
}

func parse(data string) (TOMLData, error) {
	sections := make(TOMLData)
	var currentSection string

	for _, line := range strings.Split(data, "\n") {
		line = strings.TrimSpace(stripComment(line))
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			currentSection = strings.ToLower(strings.TrimSpace(line[1 : len(line)-1]))
			sections[currentSection] = make(map[string]string)
			continue
		}

		if eq := strings.Index(line, "="); eq > 0 {
			key := strings.TrimSpace(line[:eq])
			val := strings.TrimSpace(line[eq+1:])
			if currentSection != "" {
				sections[currentSection][key] = val
			}
		}
	}
	return sections, nil
}

func unquote(s string) string {
	s = strings.TrimSpace(s)
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func stripComment(s string) string {
	if i := strings.Index(s, "#"); i >= 0 {
		return s[:i]
	}
	return s
}
