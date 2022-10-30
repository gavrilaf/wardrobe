package utils

import "time"

func TimeToJsonString(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
