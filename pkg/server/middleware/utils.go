package middleware

import "time"

func GetDurationInMilliseconds(start time.Time) float64 {
	end := time.Now().UTC()
	duration := end.Sub(start)
	milliseconds := float64(duration) / float64(time.Millisecond)
	rounded := float64(int(milliseconds*100+.5)) / 100
	return rounded
}
