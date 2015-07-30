package main

import (
	"github.com/michigan-com/newsFetch/lib"
	"fmt"
	"github.com/bitly/go-simplejson"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type PhotoInfo struct {
	Url    string
	Width  int
	Height int
}

type Photo struct {
	Caption   string
	Credit    string
	Full      PhotoInfo
	Thumbnail PhotoInfo
}

type Article struct {
	Id          int
	Headline    string
	Subheadline string
	Section     string
	Subsection  string
	Source      string
	Summary     string
	Created_at  time.Time
	Url         string
	Photo       Photo
}

type Snapshot struct {
	Articles   []Article
	Created_at time.Time
}

func getArticleId(url string) int {
	// Given an article url, get the ID from it
	r := regexp.MustCompile("/([0-9]+)/{0,1}$")
	match := r.FindStringSubmatch(url)

	if len(match) > 1 {
		i, err := strconv.Atoi(match[1])
		if err != nil {
			return -1
		}
		return i
	} else {

		fmt.Println("Failed to get ID from %s", url)
		return -1
	}
}

// Fetch a feed url and parse the articles that get retuned
// Each successfully parsed article
func getFeedUrl(url string) []Article {
	fmt.Println("Fetching %s", url)

	articles := make([]Article, 0)

	resp, err := http.Get(url)
	if err != nil {
		log.Print(fmt.Sprintf("Error fetching %s: %e", url, err))
		return articles
	}
	log.Print(fmt.Sprintf("Successfully fetched %s", url))

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		log.Print(fmt.Sprintf("Error parseing response body for %s: %e", url, err))
		return articles
	}

	content := json.Get("content")
	arrContent := content.MustArray()

	replace := regexp.MustCompile("^w{3}[.](.+)[.].+$")
	match := replace.FindStringSubmatch(resp.Request.URL.Host)
	if len(match) < 2 {
		log.Print(fmt.Println("Could not parse %s for host", resp.Request.URL.Host))
		return articles
	}
	site := match[1]

	for i := 0; i < len(arrContent); i++ {
		articleJson := content.GetIndex(i)
		photoAttrs := articleJson.Get("photo_attrs")
		ssts := articleJson.Get("ssts")
		articleUrl := fmt.Sprintf("http://%s.com%s", site, articleJson.Get("url").MustString())
		articleId := getArticleId(articleUrl)

		// Check to make sure we could parse the ID
		if articleId < 0 {
			continue
		}

		//fmt.Println("Saving article %s", article.headline)
		article := Article{
			Id:          articleId,
			Headline:    articleJson.Get("headline").MustString(),
			Subheadline: articleJson.Get("attrs").Get("brief").MustString(),
			Section:     ssts.Get("section").MustString(),
			Subsection:  ssts.Get("subsection").MustString(),
			Source:      site,
			Summary:     articleJson.Get("summary").MustString(),
			Created_at:  time.Now(),
			Url:         articleUrl,
			Photo: Photo{
				photoAttrs.Get("caption").MustString(),
				photoAttrs.Get("credit").MustString(),
				PhotoInfo{
					strings.Join([]string{photoAttrs.Get("publishurl").MustString(), photoAttrs.Get("basename").MustString()}, ""),
					photoAttrs.Get("oimagewidth").MustInt(),
					photoAttrs.Get("oimageheight").MustInt(),
				},
				PhotoInfo{
					"TODO",
					photoAttrs.Get("simagewidth").MustInt(),
					photoAttrs.Get("simageheight").MustInt(),
				},
			},
		}

		articles = append(articles, article)
	}

	return articles
}

func formatUrls() []string {

	sites := lib.Sites
	sections := lib.Sections
	urls := make([]string, 0)

	for i := 0; i < len(sites); i++ {
		site := sites[i]
		for j := 0; j < len(sections); j++ {
			section := sections[j]

			if strings.Contains(site, "detroitnews") && section == "life" {
				section += "-home"
			}
			url := fmt.Sprintf("http://%s/feeds/live/%s/json", site, section)
			urls = append(urls, url)
		}
	}

	return urls
}

func FetchArticles() {
	log.Print("Fetching articles")

	// Fetch articles from urls
	urls := formatUrls()
	c := make(chan []Article)
	articles := make([]Article, 0)
	for i := 0; i < len(urls); i++ {

		go func(url string, c chan<- []Article) {
			c <- getFeedUrl(url)

		}(urls[i], c)
	}

	// Wait until all return and parse the results
	for i := 0; i < len(urls); i++ {
		returnedArticles := <-c
		for j := 0; j < len(returnedArticles); j++ {
			articles = append(articles, returnedArticles[j])
		}
	}

	// DB stuff
	session := lib.DBConnect()
	defer session.Close()

	snapshotCollection := session.DB("mapi").C("Snapshot")
	err := snapshotCollection.Insert(&Snapshot{
		Articles:   articles,
		Created_at: time.Now(),
	})

	if err != nil {
		panic(err)
	}

	log.Print("Saved a batch of articles")
}
