package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
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
	Debugger.Println("Fetching: ", url)

	articleChan := &ArticleUrlsChan{}

	resp, err := http.Get(url)
	if err != nil {
		articleChan.Err = err
		ch <- articleChan
		return
	}
	defer resp.Body.Close()

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

func ProcessSummaries(ch chan error) {
	url := "http://brevity.detroitnow.io/newsfetch-summarize/"
	Debugger.Println("Fetching: ", url)

	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		ch <- err
	}

	ch <- nil
}

type Middleware interface {
	Process(*Article) error
}

type ArticleIn struct {
	Site    string
	Url     string
	Article *struct {
		Id       int `json:"id"`
		*Ssts    `json:"ssts"`
		Metadata *struct {
			Dates *struct {
				Timestamp string `json:"lastupdated"`
			} `json:"dates"`
		} `json:"metadata"`
		Photo *struct {
			AssetMetadata *struct {
				Attrs *Attrs `json:"items"`
			} `json:"asset_metadata"`
		} `json:"lead_photo"`
	} `json:"article"`
	Metadata *struct {
		Headline    string `json:"headline"`
		Description string `json:"description"`
		Brief       string `json:"description"`
	} `json:"metadata"`
}

func NewArticleIn(url string) *ArticleIn {
	article := &ArticleIn{Url: url}
	return article
}

func (a *ArticleIn) String() string {
	return fmt.Sprintf("<ArticleIn Site: %s, Id: %d, Url: %s>", a.Site, a.Article.Id, a.Url)
}

func (a *ArticleIn) GetData() error {
	json_url := ""
	if strings.HasSuffix(a.Url, "/") {
		json_url = fmt.Sprintf("%s%s", a.Url, "json")
	} else {
		json_url = fmt.Sprintf("%s/%s", a.Url, "json")
	}

	Debugger.Println("Fetching: ", json_url)

	resp, err := http.Get(json_url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	site, err := GetSiteFromHost(resp.Request.URL.Host)
	if err != nil {
		return err
	}

	a.Site = site

	var jso []byte
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&a)
	if err != nil {
		return err
	}

	json.Unmarshal(jso, a)

	return nil
}

func (a *ArticleIn) IsValid() bool {
	if a.Article == nil {
		Debugger.Println("Article struct missing ...")
		return false
	}

	if a.Article.Id == 0 {
		Debugger.Println("Article ID missing ...")
		return false
	}

	if isBlacklisted(a.Url) {
		Debugger.Println("Article URL has been blacklisted: ", a)
		return false
	}

	if a.Article.Photo == nil {
		Debugger.Println("Failed to find photo object: ", a)
		return false
	}

	if a.Article.Photo.AssetMetadata == nil {
		Debugger.Println("Failed to find asset_metadata object: ", a)
		return false
	}

	if a.Article.Photo.AssetMetadata.Attrs == nil {
		Debugger.Println("Failed to find photo.attrs object: ", a)
		return false
	}

	return true
}

func GetSiteFromHost(host string) (string, error) {
	replace := regexp.MustCompile("^w{3}[.](.+)[.].+$")
	match := replace.FindStringSubmatch(host)
	if len(match) < 2 {
		return "", fmt.Errorf("Could not parse %s for host", host)
	}

	return match[1], nil
}

func (a *ArticleIn) Process(article *Article) error {
	art := a.Article
	ssts := art.Ssts
	attrs := art.Photo.AssetMetadata.Attrs

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

	timestamp, aerr := time.Parse("2006-1-2T15:04:05.0000000", art.Metadata.Dates.Timestamp)
	if aerr != nil {
		timestamp = time.Now()
		Debugger.Println("Error parsing timestamp: ", aerr)
	}

	article.ArticleId = art.Id
	article.Headline = a.Metadata.Headline
	article.Subheadline = a.Metadata.Description
	article.Section = ssts.Section
	article.Subsection = ssts.Subsection
	article.Created_at = time.Now()
	article.Updated_at = time.Now()
	article.Timestamp = timestamp
	article.Url = a.Url
	article.Photo = &photo

	return nil
}
