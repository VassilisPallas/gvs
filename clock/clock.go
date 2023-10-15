package clock

import "time"

type Clock struct {
}

func (Clock) Now() time.Time {
	return time.Now()
}

func (c Clock) GetDiffInHoursFromNow(u time.Time) float64 {
	return c.Now().Sub(u).Hours()
}

func (Clock) IsBefore(t1 time.Time, t2 time.Time) bool {
	return t1.Before(t2)
}

func (Clock) IsAfter(t1 time.Time, t2 time.Time) bool {
	return t1.After(t2)
}
