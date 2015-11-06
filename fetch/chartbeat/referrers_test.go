package fetch

import (
	"testing"
)

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
