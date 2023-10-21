package testutils

import (
	"time"

	"github.com/VassilisPallas/gvs/clock"
)

type FakeClock struct {
	GetDiffInHoursFromNowRes float64
	UseRealIsBefore          bool
	MockIsBefore             bool
	UseRealIsAfter           bool
	MockIsfter               bool
}

func (FakeClock) Now() time.Time {
	return time.Time{}
}

func (c FakeClock) GetDiffInHoursFromNow(u time.Time) float64 {
	return c.GetDiffInHoursFromNowRes
}

func (c FakeClock) IsBefore(u1 time.Time, u2 time.Time) bool {
	if c.UseRealIsBefore {
		return clock.RealClock{}.IsBefore(u1, u2)
	}

	return c.MockIsBefore
}

func (c FakeClock) IsAfter(u1 time.Time, u2 time.Time) bool {
	if c.UseRealIsAfter {
		return clock.RealClock{}.IsAfter(u1, u2)
	}

	return c.MockIsfter
}
