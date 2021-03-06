package fetch

import (
	"testing"
)

func TestGetRecents(t *testing.T) {
	url := "http://api.chartbeat.com/live/recent/v3/?apikey=317a25eccba186e0f6b558f45214c0e7&host=gizmodo.com&limit=100"
	resp, err := GetRecents(url)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if resp.Source != "gizmodo" {
		t.Fatalf("Should be gizmodo, is %s", resp.Source)
	}
	if len(resp.Recents) != 100 {
		t.Fatalf("Should be 100 recents, there are %d recents", len(resp.Recents))
	}

	// Now try some failure cases
	url = "http://google.com"
	resp, err = GetRecents(url)
	if err == nil {
		t.Fatalf("Url %s should have failed, it didn't", url)
	}

	url = "asdfasdf asdf asdf"
	resp, err = GetRecents(url)
	if err == nil {
		t.Fatalf("Url should have failed. It didn't")
	}
}
