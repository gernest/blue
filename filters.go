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

//IsTimeStamp checks whether the given key is timestamp.For this function,
//timestamp is expected to be a field with key=="timestamp" and the value is a
//float64, which is a normal json number.
//
// The value is also assument to be time in nanoseconds since the UNIX Epoch.
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

//IsTag checks if key is tag. This is always false.
func IsTag(key string) bool {
	return false
}

//IsField checks whether the key is a field. This is always true except when the
//key is either "measurement" or "timestamp"
func IsField(key string) bool {
	return key != "measurement" && key != "timestamp"
}

//IsMeasurement checks if the key is measuremet.
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
