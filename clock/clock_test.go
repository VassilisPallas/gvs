package clock_test

import (
	"testing"
	"time"

	"github.com/VassilisPallas/gvs/clock"
)

func TestNow(t *testing.T) {
	c := clock.RealClock{}
	res := c.Now().Format("20060102150405")
	now := time.Now().Format("20060102150405")

	if res != now {
		t.Errorf("time should be %q, instead got %q", now, res)
	}
}

func TestGetDiffInHoursFromNow(t *testing.T) {
	c := clock.RealClock{}

	now := time.Now().Add(-time.Hour * 24) // subtract 24 hours from now
	res := c.GetDiffInHoursFromNow(now)

	// if the difference is not 24 hours raise error
	if int(res) != 24 {
		t.Errorf("time should be %f", res)
	}
}

func TestIsBefore(t *testing.T) {
	testCases := []struct {
		testTitle      string
		u1             time.Time
		u2             time.Time
		expectedResult bool
	}{
		{
			testTitle:      "should return true",
			u1:             time.Now().Add(-time.Hour * 1), // subtract one hour from now
			u2:             time.Now(),
			expectedResult: true,
		},
		{
			testTitle:      "should return false",
			u1:             time.Now(),
			u2:             time.Now().Add(-time.Hour * 1), // subtract one hour from now
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			c := clock.RealClock{}
			is := c.IsBefore(tc.u1, tc.u2)
			if tc.expectedResult != is {
				t.Errorf("value should be %t, instead got %t", tc.expectedResult, is)
			}
		})
	}

}

func TestIsAfter(t *testing.T) {
	testCases := []struct {
		testTitle      string
		u1             time.Time
		u2             time.Time
		expectedResult bool
	}{
		{
			testTitle:      "should return true",
			u1:             time.Now().Add(time.Hour * 1),
			u2:             time.Now(),
			expectedResult: true,
		},
		{
			testTitle:      "should return false",
			u1:             time.Now(),
			u2:             time.Now().Add(time.Hour * 1),
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testTitle, func(t *testing.T) {
			c := clock.RealClock{}
			is := c.IsAfter(tc.u1, tc.u2)
			if tc.expectedResult != is {
				t.Errorf("value should be %t, instead got %t", tc.expectedResult, is)
			}
		})
	}

}
