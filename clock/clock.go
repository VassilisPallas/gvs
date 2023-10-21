// Package clock provides an interface to for using
// time and duration functions.
package clock

import "time"

// Clock is the interface that wraps the basic methods for time and duration.
type Clock interface {
	// Now returns a time.Time instance of the local time.
	Now() time.Time

	// GetDiffInHoursFromNow returns the difference in float64 between now and the given datetime
	// in hours.
	GetDiffInHoursFromNow(u time.Time) float64

	// IsBefore returns whether the first date is before the second date.
	IsBefore(u1 time.Time, u2 time.Time) bool

	// IsAfter returns whether the first date is after the second date.
	IsAfter(u1 time.Time, u2 time.Time) bool
}

// RealClock is the struct that implements the Clock interface
type RealClock struct {
}

// Now returns a time.Time instance of the local time.
func (RealClock) Now() time.Time {
	return time.Now()
}

// GetDiffInHoursFromNow returns the difference in float64 between now and the given datetime
// in hours.
func (c RealClock) GetDiffInHoursFromNow(u time.Time) float64 {
	return c.Now().Sub(u).Hours()
}

// IsBefore returns whether the first date is before the second date.
func (RealClock) IsBefore(u1 time.Time, u2 time.Time) bool {
	return u1.Before(u2)
}

// IsAfter returns whether the first date is after the second date.
func (RealClock) IsAfter(u1 time.Time, u2 time.Time) bool {
	return u1.After(u2)
}
