package fetch

import (
	"testing"

	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model/chartbeat"
)

func TestGetQuickStats(t *testing.T) {
	// Test API URL from chartbeat docs
	url := "http://api.chartbeat.com/live/quickstats/v4/?apikey=317a25eccba186e0f6b558f45214c0e7&host=gizmodo.com"

	quickStats, err := GetQuickStats(url)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if quickStats.Source != "gizmodo" {
		t.Fatalf("Source should be %s, it's %s", "gizmodo", quickStats.Source)
	}

	url = "this is an invalid url"
	_, err = GetQuickStats(url)
	if err == nil {
		t.Fatalf("Error should have been thrown for url '%s'", url)
	}

	url = "http://google.com"
	_, err = GetQuickStats(url)
	if err == nil {
		t.Fatalf("Error should have been thrown for url '%s'", url)
	}
}

func TestGetHostFromParams(t *testing.T) {
	// Test API URL from chartbeat docs
	url := "http://api.chartbeat.com/live/quickstats/v4/?apikey=317a25eccba186e0f6b558f45214c0e7&host=gizmodo.com"
	host, err := GetHostFromParams(url)
	if err != nil {
		t.Fatalf("%v", err)
	} else if host == "" {
		t.Fatalf("Url %s should have host gizmodo.com", url)
	}

	// This is an invalud url
	url = "aasdf asdf asdf"
	host, err = GetHostFromParams(url)
	if host != "" {
		t.Fatalf("Host should be an empty string")
	}

	// This should get the host
	url = "this.com?host=freep.com&asdfasdf=this&what=huh"
	host, err = GetHostFromParams(url)
	if err != nil {
		t.Fatalf("%v", err)
	} else if host == "" {
		t.Fatalf("Url %s should have host freep.com", url)
	}

}

func TestSortQuickStats(t *testing.T) {
	numQuickStats := 20
	quickStats := make([]*m.QuickStats, 0, numQuickStats)
	for i := 0; i < numQuickStats; i++ {
		quickStat := &m.QuickStats{}
		quickStat.Visits = lib.RandomInt(100)
		quickStats = append(quickStats, quickStat)
	}

	sorted := SortQuickStats(quickStats)

	if !confirmQuickStatsSort(sorted) {
		t.Fatalf("Sort failed: %v", sorted)
	}
}

func confirmQuickStatsSort(sorted []*m.QuickStats) bool {
	lastVal := -1
	for i, stats := range sorted {
		if lastVal == -1 {
			lastVal = stats.Visits
			continue
		}

		if lastVal < stats.Visits {
			lib.Debugger.Printf("sorted[%d] == %d, sorted[%d] == %d, should be sorted in descending order", i-1, lastVal, i, stats.Visits)
			return false
		}
	}

	return true
}
