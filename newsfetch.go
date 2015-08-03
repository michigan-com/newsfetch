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

const maxarticles = 20 // Expected number of articles to be returned per URL

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
		return -1
	}
}

// Fetch a feed url and parse the articles that get retuned
// Each successfully parsed article
func getFeedUrl(url string) ([]Article, error) {
	fmt.Println("Fetching %s", url)

	articles := make([]Article, 0, maxarticles)

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
		ssts := articleJson.Get("ssts")
		articleUrl := fmt.Sprintf("http://%s.com%s", site, articleJson.Get("url").MustString())
		articleId := getArticleId(articleUrl)

		// Check to make sure we could parse the ID
		if articleId < 0 {
			log.Print(fmt.Sprintf("Failed to parse an article ID, likely not a news article: %s", articleUrl))
			continue
		}

		photoAttrs, ok := articleJson.Get("photo").CheckGet("attrs")
		photo := Photo{}
		if !ok {
			log.Print(fmt.Sprintf("Failed to get photos for %s", articleUrl))
		} else {
			// Height/width stuff
			owidth, _ := strconv.Atoi(photoAttrs.Get("oimagewidth").MustString())
			oheight, _ := strconv.Atoi(photoAttrs.Get("oimageheight").MustString())
			swidth, _ := strconv.Atoi(photoAttrs.Get("simageWidth").MustString())
			sheight, _ := strconv.Atoi(photoAttrs.Get("simageheight").MustString())

			// URLs
			publishUrl := photoAttrs.Get("publishurl").MustString()
			photoUrl := strings.Join([]string{publishUrl, photoAttrs.Get("basename").MustString()}, "")
			thumbUrl := ""
			if smallBaseName, ok := photoAttrs.CheckGet("smallbasename"); ok {
				thumbUrl = strings.Join([]string{publishUrl, smallBaseName.MustString()}, "")
			} else if thumbPath, ok := photoAttrs.CheckGet("thumbnailPath"); ok {
				thumbUrl = strings.Join([]string{publishUrl, thumbPath.MustString()}, "")
			}

			photo = Photo{
				photoAttrs.Get("caption").MustString(),
				photoAttrs.Get("credit").MustString(),
				PhotoInfo{
					photoUrl,
					owidth,
					oheight,
				},
				PhotoInfo{
					thumbUrl,
					swidth,
					sheight,
				},
			}
		}

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
			Photo:       photo,
		}

		articles = append(articles, article)
	}

	return articles, nil
}

func formatUrls() []string {

	sites := lib.Sites
	sections := lib.Sections
	urls := make([]string, 0, len(sites)*len(sections))

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
	articles := make([]Article, 0, len(urls)*maxarticles)

	for i := 0; i < len(urls); i++ {
		wg.Add(1)
		go func(url string) {
			// Fetch the url
			returnedArticles, err := getFeedUrl(url)
			if err != nil {
				log.Print(err)
			} else {
				// If we returned successfully, append all the articles we found
				articles = append(articles, returnedArticles...)
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
