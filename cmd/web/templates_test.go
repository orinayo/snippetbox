package main

import (
	"testing"
	"time"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name     string
		testTime time.Time
		want     string
	}{
		{
			name:     "UTC",
			testTime: time.Date(2020, 12, 17, 10, 0, 0, 0, time.UTC),
			want:     "17 Dec 2020 at 10:00",
		},
		{
			name:     "Empty",
			testTime: time.Time{},
			want:     "",
		},
		{
			name:     "CET",
			// UTC + 1
			testTime: time.Date(2020, 12, 17, 10, 0, 0, 0, time.FixedZone("CET", 1*60*60)),
			want:     "17 Dec 2020 at 09:00",
		},
	}

	for _, tt := range tests {
		// Use the t.Run() function to run a sub-test for each test case. The
		// first parameter to this is the name of the test (which is used to
		// identify the sub-test in any log output) and the second parameter is
		// and anonymous function containing the actual test for each case.
		t.Run(tt.name, func(t *testing.T) {

			humanReadableDate := humanDate(tt.testTime)

			if humanReadableDate != tt.want {
				t.Errorf("want %q; got %q", tt.want, humanReadableDate)
			}
		})

	}
}
