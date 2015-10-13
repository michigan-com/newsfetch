package fetch

import (
	"testing"
)

func TestArticleModelWithPhoto(t *testing.T) {
	url := "http://www.detroitnews.com/story/news/nation/2015/10/13/planned-parenthood-fetal-tissue/73861022/"

	processor := ParseArticleAtURL(url, true)
	t.Log(processor)
	if processor.Err != nil {
		t.Fatalf("Failed to process article: %s", processor.Err)
	}

}

func TestArticleModelWithoutPhoto(t *testing.T) {
	//url := "http://www.detroitnews.com/story/news/nation/2015/10/13/chicago-schools-indictment/73858210/"
}
