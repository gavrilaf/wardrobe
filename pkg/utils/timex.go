package utils

import "time"

func Date(year int, month time.Month, day int) time.Time {
	return DT(year, month, day, 0, 0, 0, 0)
}

func DT(year int, month time.Month, day, hour, min, sec, nsec int) time.Time {
	return time.Date(year, month, day, hour, min, sec, nsec, time.UTC)
}

func TimeToJsonString(t time.Time) string {
	return t.UTC().Format(time.RFC3339Nano)
}
