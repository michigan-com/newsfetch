package fetch

import (
	"testing"
)

func TestFetchReferrers(t *testing.T) {
	url := "http://api.chartbeat.com/live/referrers/v3/?apikey=317a25eccba186e0f6b558f45214c0e7&host=gizmodo.com"

	referrers := FetchReferrers(url)
	if len(referrers) != 1 {
		t.Fatalf("urls should have length %d, instead has length %d", 1, len(urls))
	}
}
