package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/michigan-com/newsfetch/lib"
)

var logger = lib.GetLogger()

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
	Timestamp   time.Time
	Url         string
	Photo       *Photo
	BodyText    string
}

type Snapshot struct {
	Articles   []*Article
	Created_at time.Time
}

type Feed struct {
	Site string
	Body *simplejson.Json
}

func FetchAndParseArticles(sites []string, sections []string, extractBody bool) []*Article {
	logger.Info("Fetching articles ...")

	// Fetch articles from urls
	var wg sync.WaitGroup
	urls := FormatFeedUrls(sites, sections)
	logger.Debug("%v", urls)

	articles := make([]*Article, 0, len(urls)*maxarticles)

	for i := 0; i < len(urls); i++ {
		wg.Add(1)
		go func(url string) {
			feedContent, err := GetFeedContent(url)
			if err != nil {
				logger.Warning("%v", err)
				wg.Done()
				return
			}

			logger.Info("Parsing feed content ...")
			content := feedContent.Body.Get("content")
			contentArr, err := content.Array()
			for i := 0; i < len(contentArr); i++ {
				article, err := ParseArticle(feedContent.Site, content.GetIndex(i), extractBody)
				if err != nil {
					logger.Warning("%v", err)
					continue
				}
				articles = append(articles, article)
			}
			wg.Done()
		}(urls[i])
	}

	// Wait for all the fetching to return and save the data
	wg.Wait()
	logger.Info("Done fetching and parsing URLs ...")
	return articles
}

func FormatFeedUrls(sites []string, sections []string) []string {
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

func GetFeedContent(url string) (*Feed, error) {
	logger.Info("Fetching %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("Successfully fetched %s", url))

	body, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	replace := regexp.MustCompile("^w{3}[.](.+)[.].+$")
	match := replace.FindStringSubmatch(resp.Request.URL.Host)
	if len(match) < 2 {
		return nil, fmt.Errorf("Could not parse %s for host", resp.Request.URL.Host)
	}
	site := match[1]

	feedContent := &Feed{
		site,
		body,
	}

	return feedContent, nil
}

func ParseArticle(site string, articleJson *simplejson.Json, extractBody bool) (*Article, error) {
	//logger.Debug(site)
	//logger.Debug("%v", articleJson)

	ssts := articleJson.Get("ssts")
	articleUrl := fmt.Sprintf("http://%s.com%s", site, articleJson.Get("url").MustString())
	articleId := lib.GetArticleId(articleUrl)

	// Check to make sure we could parse the ID
	if articleId < 0 {
		return &Article{}, fmt.Errorf("Failed to parse an article ID, likely not a news article: %s", articleUrl)
	}

	photoAttrs, err := articleJson.Get("photo").CheckGet("attrs")
	photo := Photo{}
	if err == false {
		return &Article{}, fmt.Errorf("Failed to get photos for %s", articleUrl)
	}

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

	body := ""
	var aerr error
	if extractBody {
		body, aerr = lib.ExtractBodyFromURL(articleUrl)
		if aerr != nil {
			return &Article{}, fmt.Errorf("Failed to extract body from article at %s", articleUrl)
		}

		logger.Info("Extracted body contains %d characters, %d paragraphs.", len(strings.Split(body, "")), len(strings.Split(body, "\n\n")))
	}

	timestamp, aerr := time.Parse("2006-1-2T15:04:05.0000000", articleJson.Get("timestamp").MustString())
	if aerr != nil {
		timestamp = time.Now()
		logger.Warning("%v", aerr)
	}

	article := &Article{
		Id:          articleId,
		Headline:    articleJson.Get("headline").MustString(),
		Subheadline: articleJson.Get("attrs").Get("brief").MustString(),
		Section:     ssts.Get("section").MustString(),
		Subsection:  ssts.Get("subsection").MustString(),
		Source:      site,
		Summary:     articleJson.Get("summary").MustString(),
		Created_at:  time.Now(),
		Timestamp:   timestamp,
		Url:         articleUrl,
		Photo:       &photo,
		BodyText:    body,
	}

	return article, nil
}

func SaveArticles(mongoUri string, articles []*Article) error {
	// DB stuff
	session := lib.DBConnect(mongoUri)
	defer lib.DBClose(session)

	// Save the snapshot
	snapshotCollection := session.DB("mapi").C("Snapshot")
	err := snapshotCollection.Insert(&Snapshot{
		Articles:   articles,
		Created_at: time.Now(),
	})

	if err != nil {
		return err
	}

	logger.Info("Saved a batch of articles ...")
	return nil
}
