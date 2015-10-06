package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

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
	return fmt.Sprintf("<Article Id: %d, Headline: %s, Url: %s>", a.ArticleId, a.Headline, a.Url)
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

func SaveArticles(mongoUri string, articles []*Article) error {
	session := DBConnect(mongoUri)
	defer DBClose(session)

	for _, article := range articles {
		article.Save(session)
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
	Id        string `json:"id"`
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
