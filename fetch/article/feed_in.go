package fetch

import (
	"encoding/json"
	"net/http"
)

type ArticleFeedIn struct {
	Content []*struct {
		Url string `json:"url"`
	} `json:"content"`
}

func (a *ArticleFeedIn) Urls() []string {
	articleUrls := make([]string, 0, len(a.Content))
	for _, content := range a.Content {
		articleUrls = append(articleUrls, content.Url)
	}
	return articleUrls
}

type ArticleUrlsChan struct {
	Urls []string
	Err  error
}

func GetArticleUrlsFromFeed(url string, ch chan *ArticleUrlsChan) {
	artDebugger.Println("Fetching: ", url)

	articleChan := &ArticleUrlsChan{}

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		articleChan.Err = err
		ch <- articleChan
		return
	}

	a := ArticleFeedIn{}
	var jso []byte
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&a)
	if err != nil {
		articleChan.Err = err
		ch <- articleChan
		return
	}

	json.Unmarshal(jso, a)
	articleChan.Urls = a.Urls()
	ch <- articleChan
}
