package lib

import (
	"fmt"
	"testing"
	"time"
)

// Test the adding of an hour interval to an article.Visits array
func TestAddHourInterval(t *testing.T) {
	article := &Article{}
	timeVal := time.Now()
	numToAdd := 20

	for i := 0; i < numToAdd; i++ {
		addHourInterval(article, RoundHourDown(timeVal), i)
		if len(article.Visits) != i+1 {
			t.Fatalf(fmt.Sprintf("Should have %d value in article.Visits, have %d", i+1, len(article.Visits)))
		}

		interval := article.Visits[i]
		if interval.Timestamp.Hour() != timeVal.Hour() {
			t.Fatalf(fmt.Sprintf("Hours should match. %d!=%d", interval.Timestamp.Hour(), timeVal.Hour()))
		}
		t.Logf("MInute: %d, second: %d, nanosecond: %d", interval.Timestamp.Minute(), interval.Timestamp.Second(), interval.Timestamp.Nanosecond())
		if interval.Timestamp.Minute() != 0 {
			t.Fatalf("Mintues should be rounded to zero")
		}
		if interval.Timestamp.Second() != 0 {
			t.Fatalf("Seconds shuold be rounded to zero")
		}
		if interval.Timestamp.Nanosecond() != 0 {
			t.Fatalf("Milliseconds should be rounded to zero")
		}
		if interval.Max != i {
			t.Fatalf(fmt.Sprintf("Max should match. %d != %d", interval.Max, i))
		}

		timeVal = timeVal.Add(1 * time.Hour)
	}
}

// Tests CheckHourlyMax when len(article.Visits) == 0
func TestHourlyMaxAdd(t *testing.T) {
	article := &Article{}
	timeVal := time.Now()
	visits := 100

	CheckHourlyMax(article, timeVal, visits)
	if len(article.Visits) != 1 {
		t.Fatalf("Should be exactly one value in article.Visits")
	}

	interval := article.Visits[0]
	if interval.Max != visits {
		t.Fatalf(fmt.Sprintf("interval.Max == %d, should be %d", interval.Max, visits))
	}
	if !RoundHourDown(timeVal).Equal(interval.Timestamp) {
		t.Fatalf(fmt.Sprintf("Timestamp is %v, should be %v", interval.Timestamp, timeVal))
	}
}

// Test a simple maximum comparison: given two values within a given hour range
// make sure the max is chosen
func TestSimpleHourlyMaxReplace(t *testing.T) {
	startTime := time.Now()
	article := &Article{}
	article.Visits = []TimeInterval{
		TimeInterval{
			100,
			startTime,
		},
	}
	newVisits := 101

	// Check Hourly max, replace
	CheckHourlyMax(article, startTime, newVisits)
	interval := article.Visits[0]
	if interval.Max != newVisits {
		t.Fatalf(fmt.Sprintf("Max should have been updated to %d, it's still %d", newVisits, interval.Max))
	}

	// Check hourly max, shouldn't replace
	newVisits = 1
	CheckHourlyMax(article, startTime, newVisits)
	interval = article.Visits[0]
	if interval.Max == newVisits {
		t.Fatalf(fmt.Sprintf("Max should not have been replaced, it is now %d", newVisits))
	}
}

func TestHourlyMaxNextHour(t *testing.T) {
	article := &Article{}
	timeVal := time.Now()
	visits := 100
	numToAdd := 100

	for i := 0; i < numToAdd; i++ {
		numIntervals := i + 1
		numVisits := numIntervals * visits
		CheckHourlyMax(article, timeVal, numVisits)

		if len(article.Visits) != numIntervals {
			t.Fatalf(fmt.Sprintf("Number of intervals: (Actual) %d, (Expected) %d", len(article.Visits), numIntervals))
		}

		interval := article.Visits[i]
		if interval.Max != numVisits {
			t.Fatalf(fmt.Sprintf("Number of visits: (Actual) %d, (Exptected) %d", interval.Max, numVisits))
		}

		// + 1 hour, should add another interval
		timeVal = timeVal.Add(1 * time.Hour)
	}
}

func TestRoundHourDown(t *testing.T) {
	timeVal := time.Now()
	numChanges := 100

	// Different intervals that will be randomly added to the value
	intervals := []time.Duration{
		time.Millisecond,
		time.Second,
		time.Minute,
		time.Hour,
	}

	for i := 0; i < numChanges; i++ {
		roundedTime := RoundHourDown(timeVal)

		if roundedTime.Minute() != 0 {
			t.Fatalf(fmt.Sprintf("Time %v should have been rounded down to the nearst hour", timeVal))
		}

		if roundedTime.Hour() != timeVal.Hour() {
			t.Fatalf(fmt.Sprintf("Times do not have matching hours. Non-rounded: %v, rounded: %v", timeVal, roundedTime))
		}

		// adjust the time val
		interval := intervals[RandomInt(len(intervals)-1)]
		randomVal := RandomInt(200)

		timeVal = timeVal.Add(time.Duration(randomVal) * interval)
	}
}

func TestIsSameDay(t *testing.T) {
	today := time.Now()

	if !IsSameDay(today, today) {
		t.Fatalf(fmt.Sprintf("Days %v and %v should be same day", today, today))
	}

	tomorrow := today.Add(36 * time.Hour)

	if IsSameDay(today, tomorrow) {
		t.Fatalf(fmt.Sprintf("Days %v and %v should not be the same day", today, tomorrow))
	}

	today = today.Add(36 * time.Hour)

	if !IsSameDay(today, tomorrow) {
		t.Fatalf(fmt.Sprintf("Days %v and %v should be on the same day", today, tomorrow))
	}
}
