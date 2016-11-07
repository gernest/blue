package blue

import (
	"strings"
	"time"
)

func joinkey(a, b string) string {
	if a != "" {
		return a + "_" + b
	}
	return b
}

func IsTimeStamp(key string, value interface{}) (time.Time, bool) {
	if key == "" {
		return time.Time{}, false
	}
	low := strings.ToLower(key)
	if low == "timestamp" {
		if ns, ok := value.(float64); ok {
			return time.Unix(0, int64(ns)), true
		}
	}
	return time.Time{}, false
}

func IsTag(key string) bool {
	return false
}

func IsField(key string) bool {
	return key != "measurement" && key != "timestamp"
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
