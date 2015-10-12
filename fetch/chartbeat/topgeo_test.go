package fetch

import (
	"testing"
)

func TestGetTopGeo(t *testing.T) {
	url := "http://api.chartbeat.com/live/top_geo/v1/?apikey=317a25eccba186e0f6b558f45214c0e7&host=gizmodo.com"

	topGeos := FetchTopGeo([]string{url})

	if len(topGeos) != 1 {
		t.Fatalf("Should return array of length 1, actual length==%d", len(topGeos))
	}

	topGeo := topGeos[0]
	if topGeo.Source != "gizmodo" {
		t.Fatalf("Shold return gizmodo, returned %s", topGeo.Source)
	}
}
