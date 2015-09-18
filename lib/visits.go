package lib

import (
	"time"
)

// Given an article, a currentTime, and a currentVisits variable, check the
// last value in the article.Visits array (lastInterval).
// Does an in-memory replace of the values if necessary
//
// If currentTime.Hour() == lastInterval.Hour()
//		if currentVisits > lastInterval.Max
//				lastInterval.Max = currentVisits
// If currentTime.Hour() > lastInterval.Hour()
//		append onto article.Visits array using this hour
// If currentTime.Hour() < lastInterval.Hour()
//		ignore
func CheckHourlyMax(article *Article, currentTime time.Time, currentVisits int) {
	length := len(article.Visits)
	if length == 0 {
		addHourInterval(article, currentTime, currentVisits)
		return
	}

	lastInterval := &article.Visits[length-1]

	roundedTime := RoundHourDown(currentTime)

	// We only care about two cases: when currentTime is on the same day, and
	// when currentTime is on the next day
	//
	// If these two time.Time objects are the same day...
	if IsSameDay(roundedTime, lastInterval.Timestamp) {

		// ...compare the hours. If they're equal...
		if roundedTime.Hour() == lastInterval.Timestamp.Hour() {

			// ...compare the maxes and replace as necessary
			if currentVisits > lastInterval.Max {
				lastInterval.Max = currentVisits
				lastInterval.Timestamp = currentTime
			}
			// ...else, if it's an hour after, then add an hour interval
		} else if roundedTime.Hour() > lastInterval.Timestamp.Hour() {
			addHourInterval(article, currentTime, currentVisits)
		}
		// ...else if currentTime time.Time is after lastInterval...
	} else if currentTime.After(lastInterval.Timestamp) {
		addHourInterval(article, currentTime, currentVisits)
	}
}

// Append a new interval onto the end of article.Visits
func addHourInterval(article *Article, currentTime time.Time, currentVisits int) {
	newInterval := &TimeInterval{
		currentVisits,
		currentTime,
	}

	currentLength := len(article.Visits)
	newVisits := make([]TimeInterval, currentLength+1, currentLength+1)

	for i := 0; i < currentLength; i++ {
		newVisits[i] = article.Visits[i]
	}

	newVisits[currentLength] = *newInterval

	article.Visits = newVisits
}

func RoundHourDown(t time.Time) time.Time {
	return t.Truncate(time.Hour)
}

// Are t1 and t2 on the same day?
func IsSameDay(t1 time.Time, t2 time.Time) bool {
	return t1.Year() == t2.Year() && t1.Month() == t2.Month() && t1.Day() == t2.Day()
}
