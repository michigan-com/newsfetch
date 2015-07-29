package newsFetch

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/michigan-com/newsFetch/lib"
	"log"
	"net/http"
	"regexp"
	"strings"
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

func getUrl(url string) {
	fmt.Println("Fetching %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Print(fmt.Sprintf("Error fetching %s: %e", url, err))
		return
	}

	json, err := simplejson.NewFromReader(resp.Body)
	if err != nil {
		log.Print(fmt.Sprintf("Error parseing response body for %s: %e", url, err))
		return
	}

	session := lib.DBConnect()
	db := session.DB("mapi")
	defer session.Close()

	content := json.Get("content")
	arrContent := content.MustArray()

	replace := regexp.MustCompile("^w{3}[.](.+)[.].+$")
	match := replace.FindStringSubmatch(resp.Request.URL.Host)
	if len(match) < 2 {
		log.Print(fmt.Println("Could not parse %s for host", resp.Request.URL.Host))
		return
	}
	site := match[1]

	for i := 0; i < len(arrContent); i++ {
		articleJson := content.GetIndex(i)
		photoAttrs := articleJson.Get("photo_attrs")
		ssts := articleJson.Get("ssts")

		//fmt.Println("Saving article %s", article.headline)

		err = db.C("articles").Insert(&Article{
			articleJson.Get("headline").MustString(),
			articleJson.Get("attrs").Get("brief").MustString(),
			ssts.Get("section").MustString(),
			ssts.Get("subsection").MustString(),
			site,
			articleJson.Get("summary").MustString(),
			time.Now(),
			fmt.Sprintf("http://%s.com%s", site, articleJson.Get("url").MustString()),
			Photo{
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
		})

		if err != nil {
			panic(err)
		}
	}

	log.Print(fmt.Sprintf("Successfully fetched %s", url))
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

	urls := formatUrls()
	for i := 0; i < len(urls); i++ {

		go getUrl(urls[i])
	}
}
