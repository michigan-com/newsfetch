package fetch

import (
	"testing"
)

func TestFetchReferrers(t *testing.T) {
	url := "http://api.chartbeat.com/live/referrers/v3/?apikey=317a25eccba186e0f6b558f45214c0e7&host=gizmodo.com"

	referrers := FetchReferrers([]string{url})
	if len(referrers) != 1 {
		t.Fatalf("urls should have length %d, instead has length %d", 1, len(referrers))
	}

	if referrers[0].Source != "gizmodo" {
		t.Fatalf("Should have gizmodo as the source, instead has %s", referrers[0].Source)
	}

	// Look for direct traffic
	directFound := false
	for source, _ := range referrers[0].Referrers {
		if source == "" {
			directFound = true
		}
	}

	if !directFound {
		t.Fatalf("direct traffic not found")
	}
}

func TestGetReferrers(t *testing.T) {
	url := "http://api.chartbeat.com/live/referrers/v3/?apikey=317a25eccba186e0f6b558f45214c0e7&host=gizmodo.com"
	_, err := getReferrers(url)
	if err != nil {
		t.Fatalf("%v", err)
	}

	url = "http://google.com"
	_, err = getReferrers(url)
	if err == nil {
		t.Fatalf("Url %s should have thrown an error", url)
	}

	url = "asdf"
	_, err = getReferrers(url)
	if err == nil {
		t.Fatalf("Url %s should have thrown an error", url)
	}
}
