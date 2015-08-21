package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var logger = GetLogger()

const maxarticles = 20 // Expected number of articles to be returned per URL

type PhotoInfo struct {
	Url    string `bson:"url"`
	Width  int    `bson:"width"`
	Height int    `bson:"height"`
}

type Photo struct {
	Caption   string    `bson:"caption"`
	Credit    string    `bson:"credit"`
	Full      PhotoInfo `bson:"full"`
	Thumbnail PhotoInfo `bson:"thumbnail"`
}

type ArticleId struct {
	Id int `bson:"_id"`
}

type Article struct {
	ArticleId   int       `bson:"article_id"`
	Headline    string    `bson:"headline"`
	Subheadline string    `bson:"subheadline"`
	Section     string    `bson:"section"`
	Subsection  string    `bson:"subsection"`
	Source      string    `bson:"source"`
	Summary     string    `bson:"summary"`
	Created_at  time.Time `bson:"created_at"`
	Timestamp   time.Time `bson:"timestamp"`
	Url         string    `bson:"url"`
	Photo       *Photo    `bson:"photo"`
	BodyText    string    `bson:"body"`
}

type Feed struct {
	Site string
	Body map[string]interface{} //*simplejson.Json
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

	//articles := make([]*Article, 0, len(urls)*maxarticles)
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

			content := feedContent.Body //.Get("content")
			cslice, _ := content["content"].([]interface{})

			for _, articleJson := range cslice {
				jso := articleJson.(map[string]interface{})
				url := jso["url"].(string)
				articleUrl := fmt.Sprintf("http://%s.com%s", feedContent.Site, url)

				if articleMap[articleUrl] != nil || isBlacklisted(articleUrl) {
					continue
				}

				article, err := ParseArticle(articleUrl, jso, extractBody)
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

func ParseArticle(articleUrl string, articleJson map[string]interface{}, extractBody bool) (*Article, error) {
	//logger.Debug(site)
	//logger.Debug("%v", articleJson)

	ssts := articleJson["ssts"].(map[string]interface{})
	articleId := GetArticleId(articleUrl)

	// Check to make sure we could parse the ID
	if articleId < 0 {
		return &Article{}, fmt.Errorf("Failed to parse an article ID, likely not a news article: %s", articleUrl)
	}

	photoJson, _ := articleJson["photo"].(map[string]interface{})
	attrs, _ := photoJson["attrs"].(map[string]interface{})
	if attrs == nil {
		return nil, fmt.Errorf("Failed to get photos for %s", articleUrl)
	}
	photo := Photo{}

	// Height/width stuff
	owidth, _ := attrs["oimagewidth"].(int)
	oheight, _ := attrs["oimageheight"].(int)
	swidth, _ := attrs["simageWidth"].(int)
	sheight, _ := attrs["simageheight"].(int)

	// URLs
	publishUrl, ok := attrs["publishurl"].(string)
	if !ok {
		return nil, fmt.Errorf("Failed to get photo url for %s", articleUrl)
	}
	basename, ok := attrs["basename"].(string)
	if !ok {
		return nil, fmt.Errorf("Failed to get photo filename for %s", articleUrl)
	}

	photoUrl := strings.Join([]string{publishUrl, basename}, "")
	thumbUrl := ""
	if smallBaseName, ok := attrs["smallbasename"].(string); ok {
		thumbUrl = strings.Join([]string{publishUrl, smallBaseName}, "")
	} else if thumbPath, ok := attrs["thumbnailPath"].(string); ok {
		thumbUrl = strings.Join([]string{publishUrl, thumbPath}, "")
	}

	caption, _ := attrs["caption"].(string)
	credit, _ := attrs["credit"].(string)

	photo = Photo{
		caption,
		credit,
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
		/*if aerr != nil {
			return &Article{}, fmt.Errorf("Failed to extract body from article at %s", articleUrl)
		}*/

		logger.Debug("Extracted body contains %d characters, %d paragraphs.", len(strings.Split(body, "")), len(strings.Split(body, "\n\n")))
	}

	timestamp, aerr := time.Parse("2006-1-2T15:04:05.0000000", articleJson["timestamp"].(string))
	if aerr != nil {
		timestamp = time.Now()
		logger.Warning("%v", aerr)
	}

	headline, _ := articleJson["headline"].(string)
	subheadline, _ := attrs["brief"].(string)
	section, _ := ssts["section"].(string)
	subsection, _ := ssts["subsection"].(string)
	summary, _ := articleJson["summary"].(string)

	article := &Article{
		ArticleId:   articleId,
		Headline:    headline,
		Subheadline: subheadline,
		Section:     section,
		Subsection:  subsection,
		Summary:     summary,
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
