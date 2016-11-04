package blue

import "strings"

func joinkey(a, b string) string {
	if a != "" {
		return a + "_" + b
	}
	return b
}

func IsTimeStamp(key string, value interface{}) bool {
	if key == "" {
		return false
	}
	low := strings.ToLower(key)
	if low == "timestamp" {
		return true
	}
	return false
}

func IsTag(key string) bool {
	return false
}

func IsField(key string) bool {
	return true
}

func IsMeasurement(key string, value interface{}) (string, bool) {
	if key == "" {
		return "", false
	}
	if key == "measurement" {
		if s, ok := value.(string); ok {
			return s, true
		}
	}
	return "", false
}
