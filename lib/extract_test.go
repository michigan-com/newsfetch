package lib

import (
	"testing"
)

func TestBodyExtractor(t *testing.T) {
	t.Log("Ensure body extractor produces non-empty string.")

	url := "http://www.freep.com/story/news/local/michigan/oakland/2015/08/20/police-chase-troy-bloomfield-hills-warren-absconder-shooting/32056645/"

	var actual string
	ch := make(chan string)
	go ExtractBodyFromURL(ch, url, false)

	actual = <-ch
	if actual == "" {
		t.Errorf("Body extractor returned no text.")
	}
}
