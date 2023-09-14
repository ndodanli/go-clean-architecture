package utils

import "time"

const (
	ISO8601  = "2006-01-02T15:04:05.999999Z07:00"
	Min1     = time.Minute * 1
	Min5     = time.Minute * 5
	Min10    = time.Minute * 10
	Min15    = time.Minute * 15
	Min30    = time.Minute * 30
	Hour1    = time.Hour * 1
	Infinite = time.Duration(0)
)

func UTCNow() time.Time {
	return time.Now().UTC()
}

func UTCNowAddDuration(d time.Duration) time.Time {
	return UTCNow().Add(d)
}

func UTCNowISOAddDuration(d time.Duration) string {
	return UTCNow().Add(d).Format("2006-01-02T15:04:05.999999Z07:00")
}

func UTCNowUnix() int64 {
	return UTCNow().Unix()
}

func UTCNowUnixMilli() int64 {
	return UTCNow().UnixNano() / int64(time.Millisecond)
}

func UTCNowUnixNano() int64 {
	return UTCNow().UnixNano()
}
