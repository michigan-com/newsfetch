package fetch

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var articleIdIndex = mgo.Index{
	Key:      []string{"article_id"},
	Unique:   true,
	DropDups: true,
}

func SaveArticles(mongoUri string, articles []*m.Article) error {
	session := lib.DBConnect(mongoUri)
	defer lib.DBClose(session)

	for _, article := range articles {
		article.Save(session)
	}

	return nil
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
	Articles []m.Article `json:"articles"`
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

func LoadArticles(mongoUri string) ([]*m.Article, error) {
	session := lib.DBConnect(mongoUri)
	defer lib.DBClose(session)

	articleCol := session.DB("").C("Article")

	var result []*m.Article
	err := articleCol.Find(nil).All(&result)
	return result, err
}

func LoadArticleById(mongoUri string, articleId int) (*m.Article, error) {
	session := lib.DBConnect(mongoUri)
	defer lib.DBClose(session)

	articleCol := session.DB("").C("Article")

	var result *m.Article
	err := articleCol.Find(bson.M{"article_id": articleId}).One(&result)
	return result, err
}

func LoadRemoteArticles(url string) ([]*m.Article, error) {
	artDebugger.Println("Fetching ", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	artDebugger.Println(fmt.Sprintf("Successfully fetched %s", url))

	var response MapiArticlesResponse
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return nil, err
	}

	return sliceOfArticlesToSliceOfPointers(response.Articles), nil
}

func FilterArticlesBySubsection(articles []*m.Article, section string, subsection string) []*m.Article {
	result := make([]*m.Article, 0, len(articles))
	for _, el := range articles {
		if (el.Section == section) && (el.Subsection == subsection) {
			result = append(result, el)
		}
	}
	return result
}

func FilterArticlesForRecipeExtraction(articles []*m.Article) []*m.Article {
	return FilterArticlesBySubsection(articles, "life", "food")
}

func sliceOfArticlesToSliceOfPointers(articles []m.Article) []*m.Article {
	result := make([]*m.Article, 0, len(articles))
	for _, el := range articles {
		copy := el
		result = append(result, &copy)
	}
	return result
}
