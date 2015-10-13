package model

import (
	"testing"

	a "github.com/michigan-com/newsfetch/fetch/article"
)

func TestArticleModelWithPhoto(t *testing.T) {
	url := "http://www.detroitnews.com/story/news/nation/2015/10/13/planned-parenthood-fetal-tissue/73861022/"

	article, _, _, err := a.ParseArticleAtURL(articleUrl, body /* global flag */)
	if err != nil {
		artDebugger.Println("Failed to process article: ", err)
		return
	}
}

func TestArticleModelWithoutPhoto(t *testing.T) {
	//url := "http://www.detroitnews.com/story/news/nation/2015/10/13/chicago-schools-indictment/73858210/"
}
