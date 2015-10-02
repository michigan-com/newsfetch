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

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var tokenizer = LoadTokenizer()

const maxarticles = 20 // Expected number of articles to be returned per URL

var articleIdIndex = mgo.Index{
	Key:      []string{"article_id"},
	Unique:   true,
	DropDups: true,
}

/*
 * DATA GOING OUT
 */
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

type Article struct {
	Id          bson.ObjectId  `bson:"_id,omitempty" json:"_id"`
	ArticleId   int            `bson:"article_id" json:"article_id"`
	Headline    string         `bson:"headline" json:"headline`
	Subheadline string         `bson:"subheadline" json:"subheadline"`
	Section     string         `bson:"section" json:"section"`
	Subsection  string         `bson:"subsection" json:"subsection"`
	Source      string         `bson:"source" json:"source"`
	Summary     interface{}    `bson:"summary" json:"summary"`
	Created_at  time.Time      `bson:"created_at" json:"created_at"`
	Updated_at  time.Time      `bson:"updated_at" json:"updated_at"`
	Timestamp   time.Time      `bson:"timestamp" json:"timestamp"`
	Url         string         `bson:"url" json:"url"`
	Photo       *Photo         `bson:"photo" json:"photo"`
	BodyText    string         `bson:"body" json:"body"`
	Visits      []TimeInterval `body:"visits" json:"visits"`
}

func (a *Article) String() string {
	return fmt.Sprintf("<Article Id: %d, Headline: %s, Url: %s>", a.Id, a.Headline, a.Url)
}

func (article *Article) Save(session *mgo.Session) error {
	// Save the snapshot
	articleCol := session.DB("").C("Article")
	art := Article{}
	err := articleCol.
		Find(bson.M{"article_id": article.ArticleId}).
		Select(bson.M{"_id": 1, "created_at": 1}).
		One(&art)
	if err == nil {
		article.Created_at = art.Created_at
		articleCol.Update(bson.M{"_id": art.Id}, article)
		Debugger.Println("Article updated: ", article)
	} else {
		//bulk.Insert(article)
		articleCol.Insert(article)
		Debugger.Println("Article added: ", article)
	}

	return nil
}

type TimeInterval struct {
	Max       int       `bson:"max"`
	Timestamp time.Time `bson:"timestamp"`
}

/*
 * DATA COMING IN
 */
type Feed struct {
	Site string
	Body *struct {
		Content []*Content `json:"content"`
	}
}

type Ssts struct {
	Section    string `json:"section"`
	Subsection string `json:"subsection"`
}

type Content struct {
	Url       string `json:"url"`
	Headline  string `json:"headline"`
	Summary   string `json:"summary"`
	Timestamp string `json:"timestamp"`
	*Ssts     `json:"ssts"`
	Photo     *struct {
		*Attrs `json:"attrs"`
	} `json:"photo"`
	Attrs *struct {
		Brief string `json:"brief"`
	} `json:"attrs"`
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

type MapiArticlesResponse struct {
	Articles []Article `json:"articles"`
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
	Debugger.Println("Fetching articles ...")

	// Fetch articles from urls
	var wg sync.WaitGroup

	mapMutex := &sync.RWMutex{}
	articleMap := map[int]*Article{}

	for i := 0; i < len(urls); i++ {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			feedContent, err := GetFeedContent(url)
			if err != nil {
				Debugger.Println("%v", err)
				return
			}

			for _, articleJson := range feedContent.Body.Content {
				url := articleJson.Url
				articleUrl := fmt.Sprintf("http://%s.com%s", feedContent.Site, url)
				articleId := GetArticleId(articleUrl)

				mapMutex.RLock()
				_, repeatArticle := articleMap[articleId]
				mapMutex.RUnlock()

				if !IsValidArticleId(articleId) || repeatArticle || isBlacklisted(articleUrl) {
					continue
				}

				article, err := ParseArticle(articleUrl, articleJson, extractBody)
				if err != nil {
					Debugger.Printf("%v", err)
					continue
				}

				article.Source = feedContent.Site

				mapMutex.Lock()
				articleMap[article.ArticleId] = article
				mapMutex.Unlock()
			}

		}(urls[i])
	}

	// Wait for all the fetching to return and save the data
	wg.Wait()
	Logger.Println("Done fetching and parsing URLs ...")

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
	Debugger.Println("Fetching ", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	Debugger.Println(fmt.Sprintf("Successfully fetched %s", url))

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
	if !IsValidArticleId(articleId) {
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

	var extracted *ExtractedBody
	var summary []string
	if extractBody {
		ch := make(chan *ExtractedBody)
		go ExtractBodyFromURL(ch, articleUrl, false)
		extracted = <-ch

		if extracted.Text != "" {
			Debugger.Println("Extracted extracted contains %d characters, %d paragraphs.", len(strings.Split(extracted.Text, "")), len(strings.Split(extracted.Text, "\n\n")))
			summarizer := NewPunktSummarizer(articleJson.Headline, extracted.Text, tokenizer)
			summary = summarizer.KeyPoints()
			Debugger.Println("Generated summary ...")
		}
	}

	timestamp, aerr := time.Parse("2006-1-2T15:04:05.0000000", articleJson.Timestamp)
	if aerr != nil {
		timestamp = time.Now()
		Debugger.Println("%v", aerr)
	}

	article := &Article{
		ArticleId:   articleId,
		Headline:    articleJson.Headline,
		Subheadline: articleJson.Attrs.Brief,
		Section:     ssts.Section,
		Subsection:  ssts.Subsection,
		Summary:     summary,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
		Timestamp:   timestamp,
		Url:         articleUrl,
		Photo:       &photo,
	}

	if extracted != nil {
		article.BodyText = extracted.Text
	}

	return article, nil
}

func RemoveArticles(mongoUri string) error {
	session := DBConnect(mongoUri)
	defer DBClose(session)

	Debugger.Println("Removing all articles from mongodb ...")

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
		art := Article{}
		err := articleCol.
			Find(bson.M{"article_id": article.ArticleId}).
			Select(bson.M{"_id": 1, "created_at": 1}).
			One(&art)
		if err == nil {
			article.Created_at = art.Created_at
			articleCol.Update(bson.M{"_id": art.Id}, article)
			Debugger.Println("Article updated: ", article.Url)
			totalUpdates++
		} else {
			//bulk.Insert(article)
			articleCol.Insert(article)
			Debugger.Println("Article added: ", article.Url)
			totalInserts++
		}
	}
	Logger.Printf("%d articles updated", totalUpdates)
	Logger.Printf("%d articles added", totalInserts)
	//_, err := bulk.Run()

	/*if err != nil {
		return err
	}*/

	Debugger.Println("Saved a batch of articles ...")
	return nil
}

func LoadArticles(mongoUri string) ([]*Article, error) {
	session := DBConnect(mongoUri)
	defer DBClose(session)

	articleCol := session.DB("").C("Article")

	var result []*Article
	err := articleCol.Find(nil).All(&result)
	return result, err
}

func LoadArticleById(mongoUri string, articleId int) (*Article, error) {
	session := DBConnect(mongoUri)
	defer DBClose(session)

	articleCol := session.DB("").C("Article")

	var result *Article
	err := articleCol.Find(bson.M{"article_id": articleId}).One(&result)
	return result, err
}

func LoadRemoteArticles(url string) ([]*Article, error) {
	Debugger.Println("Fetching ", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	Debugger.Println(fmt.Sprintf("Successfully fetched %s", url))

	var response MapiArticlesResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return nil, err
	}

	return sliceOfArticlesToSliceOfPointers(response.Articles), nil
}

func FilterArticlesBySubsection(articles []*Article, section string, subsection string) []*Article {
	result := make([]*Article, 0, len(articles))
	for _, el := range articles {
		if (el.Section == section) && (el.Subsection == subsection) {
			result = append(result, el)
		}
	}
	return result
}

func FilterArticlesForRecipeExtraction(articles []*Article) []*Article {
	return FilterArticlesBySubsection(articles, "life", "food")
}

func sliceOfArticlesToSliceOfPointers(articles []Article) []*Article {
	result := make([]*Article, 0, len(articles))
	for _, el := range articles {
		copy := el
		result = append(result, &copy)
	}
	return result
}
