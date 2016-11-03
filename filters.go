package blue

import "strings"

func joinkey(a, b string) string {
	if a != "" {
		return a + "_" + b
	}
	return b
}

func timeStamp(key string, value interface{}) bool {
	if key == "" {
		return false
	}
	low := strings.ToLower(key)
	if low == "timestamp" {
		return true
	}
	return false
}

func isTag(key string) bool {
	return false
}
func isField(key string) bool {
	return true
}

func isMeasurement(key string) bool {
	if key == "" {
		return false
	}
	low := strings.ToLower(key)
	return low == "measurement"
}
