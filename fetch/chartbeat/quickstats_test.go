package fetch

import (
	"os"
	"testing"

	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model"
	"gopkg.in/mgo.v2/bson"
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

func TestSaveQuickStats(t *testing.T) {
	t.Skip("No mongo tests allowed MIKE")
	mongoUri := os.Getenv("MONGO_URI")
	if mongoUri == "" {
		t.Fatalf("%v", "No mongo URI specified, failing test")
	}

	numStats := 20
	quickStats := make([]*m.QuickStats, 0, numStats)
	for i := 0; i < numStats; i++ {
		stat := &m.QuickStats{}
		stat.Visits = lib.RandomInt(100)

		quickStats = append(quickStats, stat)
	}

	quickStats = SortQuickStats(quickStats)

	// Add it a bunch of times
	SaveQuickStats(quickStats, mongoUri)
	SaveQuickStats(quickStats, mongoUri)
	SaveQuickStats(quickStats, mongoUri)
	SaveQuickStats(quickStats, mongoUri)

	// Now verify
	session := lib.DBConnect(mongoUri)
	defer lib.DBClose(session)

	col := session.DB("").C("Quickstats")
	count, err := col.Count()

	if err != nil {
		t.Fatalf("%v", err)
	}
	if count != 1 {
		t.Fatalf("Should be 1 Quickstats snapshot, there are %d", count)
	}

	snapshot := &m.QuickStatsSnapshot{}
	col.Find(bson.M{}).One(&snapshot)
	expectedLen, actualLen := len(quickStats), len(snapshot.Stats)

	if expectedLen != actualLen {
		t.Fatalf("Epected %d, actual %d", expectedLen, actualLen)
	}

	for i := 0; i < len(quickStats); i++ {
		val1, val2 := quickStats[i].Visits, snapshot.Stats[i].Visits
		if val1 != val2 {
			t.Fatalf("quickstats[%d].Visits == %d, snapshot[i].Visits == %d. They should be equal", i, val1, i, val2)
		}
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
