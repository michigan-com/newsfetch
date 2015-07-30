package main

import (
	"errors"
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/michigan-com/newsFetch/lib"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
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
func getFeedUrl(url string) ([]Article, error) {
	fmt.Println("Fetching %s", url)

	articles := make([]Article, 0)

	resp, err := http.Get(url)
	if err != nil {
		return articles, err
	}
	log.Print(fmt.Sprintf("Successfully fetched %s", url))

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return articles, err
	}

	content := json.Get("content")
	arrContent := content.MustArray()

	replace := regexp.MustCompile("^w{3}[.](.+)[.].+$")
	match := replace.FindStringSubmatch(resp.Request.URL.Host)
	if len(match) < 2 {
		return articles, errors.New(fmt.Sprintf("Could not parse %s for host", resp.Request.URL.Host))
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

	return articles, nil
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
	var wg sync.WaitGroup
	urls := formatUrls()
	articles := make([]Article, 0)

	for i := 0; i < len(urls); i++ {
		wg.Add(1)
		go func(url string) {
			// Fetch the url
			returnedArticles, err := getFeedUrl(url)
			if err != nil {
				log.Print(err)
			} else {
				// If we returned successfully, append all the articles we found
				for j := 0; j < len(returnedArticles); j++ {
					articles = append(articles, returnedArticles[j])
				}
			}
			// Donezo
			wg.Done()
		}(urls[i])
	}

	// Wait for all the fetching to return and save the data
	wg.Wait()
	log.Print("Done fetching and parsing URLs")

	// DB stuff
	session := lib.DBConnect()
	defer session.Close()

	// Save the snapshot
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
