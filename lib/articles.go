package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var logger = GetLogger()

const maxarticles = 20 // Expected number of articles to be returned per URL

type PhotoInfo struct {
	Url    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Photo struct {
	Caption   string    `json:"caption"`
	Credit    string    `json:"credit"`
	Full      PhotoInfo `json:"full"`
	Thumbnail PhotoInfo `json:"thumbnail"`
}

type ArticleId struct {
	Id int `json:"_id"`
}

type Article struct {
	ArticleId   int       `json:"article_id"`
	Headline    string    `json:"headline"`
	Subheadline string    `json:"subheadline"`
	Section     string    `json:"section"`
	Subsection  string    `json:"subsection"`
	Source      string    `json:"source"`
	Summary     string    `json:"summary"`
	Created_at  time.Time `json:"created_at"`
	Timestamp   time.Time `json:"timestamp"`
	Url         string    `json:"url"`
	Photo       *Photo    `json:"photo"`
	BodyText    string    `json:"body"`
}

type Feed struct {
	Site string
	Body *struct {
		Content []*Content `json:"content"`
	}
}

type Content struct {
	Url       string `json:"url"`
	Headline  string `json:"headline"`
	Summary   string `json:"summary"`
	Timestamp string `json:"timestamp"`
	Ssts      *struct {
		Section    string `json:"section"`
		Subsection string `json:"subsection"`
	} `json:"ssts"`
	Photo *struct {
		*Attrs `json:"attrs"`
	} `json:"photo"`
}

type Attrs struct {
	Owidth        string `json:"oimagewidth"`
	OWidth        string `json:"oimageWidth"`
	Oheight       string `json:"oimageheight"`
	Swidth        string `json:"simagewidth"`
	Sheight       string `json:"simageheight"`
	Basename      string `json:"basename"`
	PublishUrl    string `json:"publishurl"`
	SmallBasename string `json:"smallbasename"`
	ThumbnailPath string `json:"thumbnailPath"`
	Caption       string `json:"caption"`
	Credit        string `json:"credit"`
	Brief         string `json:"brief"`
}

func isBlacklisted(url string) bool {
	blacklist := []string{
		"/videos/",
		"/police-blotter/",
	}

	for _, item := range blacklist {
		if strings.Contains(url, item) {
			return true
		}
	}

	return false
}

func FetchAndParseArticles(urls []string, extractBody bool) []*Article {
	logger.Info("Fetching articles ...")

	// Fetch articles from urls
	var wg sync.WaitGroup
	logger.Debug("%v", urls)

	articleMap := map[string]*Article{}

	for i := 0; i < len(urls); i++ {
		wg.Add(1)
		go func(url string) {
			feedContent, err := GetFeedContent(url)
			if err != nil {
				logger.Warning("%v", err)
				wg.Done()
				return
			}

			for _, articleJson := range feedContent.Body.Content {
				url := articleJson.Url
				articleUrl := fmt.Sprintf("http://%s.com%s", feedContent.Site, url)

				if articleMap[articleUrl] != nil || isBlacklisted(articleUrl) {
					continue
				}

				article, err := ParseArticle(articleUrl, articleJson, extractBody)
				if err != nil {
					logger.Warning("%v", err)
					continue
				}

				article.Source = feedContent.Site
				articleMap[article.Url] = article
			}

			wg.Done()
		}(urls[i])
	}

	// Wait for all the fetching to return and save the data
	wg.Wait()
	logger.Info("Done fetching and parsing URLs ...")

	articles := make([]*Article, 0, len(articleMap))
	for _, art := range articleMap {
		articles = append(articles, art)
	}

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
	logger.Debug("Fetching %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	logger.Debug(fmt.Sprintf("Successfully fetched %s", url))

	feed := &Feed{}

	var body []byte
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&feed.Body)
	if err != nil {
		return feed, err
	}

	json.Unmarshal(body, feed.Body)

	replace := regexp.MustCompile("^w{3}[.](.+)[.].+$")
	match := replace.FindStringSubmatch(resp.Request.URL.Host)
	if len(match) < 2 {
		return nil, fmt.Errorf("Could not parse %s for host", resp.Request.URL.Host)
	}

	feed.Site = match[1]

	return feed, nil
}

func ParseArticle(articleUrl string, articleJson *Content, extractBody bool) (*Article, error) {
	ssts := articleJson.Ssts
	articleId := GetArticleId(articleUrl)

	// Check to make sure we could parse the ID
	if articleId < 0 {
		return nil, fmt.Errorf("Failed to parse an article ID, likely not a news article: %s", articleUrl)
	}

	if articleJson.Photo == nil {
		return nil, fmt.Errorf("Failed to find photo object for %s", articleUrl)
	}

	if articleJson.Photo.Attrs == nil {
		return nil, fmt.Errorf("Failed to find photo.attrs object for %s", articleUrl)
	}

	attrs := articleJson.Photo.Attrs

	photoUrl := strings.Join([]string{attrs.PublishUrl, attrs.Basename}, "")
	thumbUrl := ""
	if attrs.SmallBasename != "" {
		thumbUrl = strings.Join([]string{attrs.PublishUrl, attrs.SmallBasename}, "")
	} else if attrs.ThumbnailPath != "" {
		thumbUrl = strings.Join([]string{attrs.PublishUrl, attrs.ThumbnailPath}, "")
	}

	owidth, _ := strconv.Atoi(attrs.Owidth)
	if owidth == 0 {
		owidth, _ = strconv.Atoi(attrs.OWidth)
	}
	oheight, _ := strconv.Atoi(attrs.Oheight)
	swidth, _ := strconv.Atoi(attrs.Swidth)
	sheight, _ := strconv.Atoi(attrs.Sheight)

	photo := Photo{
		attrs.Caption,
		attrs.Credit,
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
		ch := make(chan string)
		go ExtractBodyFromURL(ch, articleUrl, false)
		body = <-ch

		logger.Debug("Extracted body contains %d characters, %d paragraphs.", len(strings.Split(body, "")), len(strings.Split(body, "\n\n")))
	}

	timestamp, aerr := time.Parse("2006-1-2T15:04:05.0000000", articleJson.Timestamp)
	if aerr != nil {
		timestamp = time.Now()
		logger.Warning("%v", aerr)
	}

	article := &Article{
		ArticleId:   articleId,
		Headline:    articleJson.Headline,
		Subheadline: attrs.Brief,
		Section:     ssts.Section,
		Subsection:  ssts.Subsection,
		Summary:     articleJson.Summary,
		Created_at:  time.Now(),
		Timestamp:   timestamp,
		Url:         articleUrl,
		Photo:       &photo,
		BodyText:    body,
	}

	return article, nil
}

func RemoveArticles(mongoUri string) error {
	session := DBConnect(mongoUri)
	defer DBClose(session)

	logger.Info("Removing all articles from mongodb ...")

	articles := session.DB("").C("Article")
	_, err := articles.RemoveAll(nil)

	return err
}

func SaveArticles(mongoUri string, articles []*Article) error {
	// DB stuff
	session := DBConnect(mongoUri)
	defer DBClose(session)

	// Save the snapshot
	articleCol := session.DB("").C("Article")
	//bulk := articleCol.Bulk()
	totalUpdates := 0
	totalInserts := 0
	for _, article := range articles {
		art := ArticleId{}
		err := articleCol.Find(bson.M{"url": article.Url}).One(&art)
		if err == nil {
			logger.Debug("Article updated: %s", article.Url)
			totalUpdates++
			articleCol.Update(bson.M{"_id": art.Id}, article)
		} else {
			//bulk.Insert(article)
			logger.Debug("Article added: %s", article.Url)
			totalInserts++
			articleCol.Insert(article)
		}
	}
	logger.Info("%d articles updated", totalUpdates)
	logger.Info("%d articles added", totalInserts)
	//_, err := bulk.Run()

	/*if err != nil {
		return err
	}*/

	logger.Info("Saved a batch of articles ...")
	return nil
}
