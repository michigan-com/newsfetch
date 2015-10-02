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

func PArticle(proc ArticleProcessor, body Parser) {
	articleIn, err := proc.GetData()
	if !proc.IsValid() {
		Debugger.Println(err)
		return
	}

	article, err := proc.Parse(articleIn)
	rawBody, err := body.GetData()
	if err != nil {
	}
	article.BodyText, err = body.Parse(rawBody)
	article.Save(nil)
}

type Parser interface {
	GetData() (interface{}, error)
	Parse() (interface{}, error)
}

type ArticleProcessor interface {
	Parser
	IsValid() bool
}

type DefaultBodyParser struct {
	Text string
}

func (b *DefaultBodyParser) Parse() string {
	return ""
}

type DefaultArticleProcessor struct {
	Processor ArticleProcessor
	Parser
	*Article
}

func NewArticleProcessor(url string) *DefaultArticleProcessor {
	return &DefaultArticleProcessor{
		Processor:  &ArticleIn{Url: url},
		BodyParser: &DefaultBodyParser{},
	}
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
	return &ArticleIn{Url: url}
}

func (a *ArticleIn) String() string {
	return fmt.Sprintf("<ArticleIn Site: %s, Id: %s, Url: %s>", a.Site, a.Article.Id, a.Url)
}

func (a *ArticleIn) IsValid() bool {
	if !IsValidArticleId(a.Article.Id) {
		Debugger.Println("Could not parse article id: %s", a)
		return false
	}

	if isBlacklisted(a.Url) {
		Debugger.Println("Article URL has been blacklisted: %s", a)
		return false
	}

	if a.Article.Photo == nil {
		Debugger.Println("Failed to find photo object for %s", a)
		return false
	}

	if a.Article.Photo.AssetMetadata == nil {
		Debugger.Println("Failed to find asset_metadata object for %s", a)
		return false
	}

	if a.Article.Photo.AssetMetadata.Attrs == nil {
		Debugger.Println("Failed to find photo.attrs object for %s", a)
		return false
	}

	return true
}

func (a *ArticleIn) GetData() (interface{}, error) {
	json_url := ""
	if strings.HasSuffix(a.Url, "/") {
		json_url = fmt.Sprintf("%s%s", a.Url, "json")
	} else {
		json_url = fmt.Sprintf("%s/%s", a.Url, "json")
	}

	Debugger.Println("Fetching: ", json_url)

	resp, err := http.Get(json_url)
	if err != nil {
		return a, err
	}

	site, err := GetSiteFromHost(resp.Request.URL.Host)
	if err != nil {
		return a, err
	}

	a.Site = site

	var jso []byte
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&a.Article)
	if err != nil {
		return a, err
	}

	json.Unmarshal(jso, a.Article)

	return a, nil
}

func (a *ArticleIn) Parse() (*Article, error) {
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
		Debugger.Println("Error parsing timestamp: %v", aerr)
	}

	article := &Article{
		ArticleId:   art.Id,
		Headline:    a.Metadata.Headline,
		Subheadline: a.Metadata.Description,
		Section:     ssts.Section,
		Subsection:  ssts.Subsection,
		Created_at:  time.Now(),
		Updated_at:  time.Now(),
		Timestamp:   timestamp,
		Url:         a.Url,
		Photo:       &photo,
	}

	return article, nil
}

func GetSiteFromHost(host string) (string, error) {
	replace := regexp.MustCompile("^w{3}[.](.+)[.].+$")
	match := replace.FindStringSubmatch(host)
	if len(match) < 2 {
		return "", fmt.Errorf("Could not parse %s for host", host)
	}

	return match[1], nil
}
