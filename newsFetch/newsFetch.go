package newsFetch

import (
	//"gopkg.in/mgo.v2"
	"../lib/"
	"fmt"
	"github.com/bitly/go-simplejson"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

var db = lib.DB

type PhotoInfo struct {
	url    string
	width  int
	height int
}

type Photo struct {
	caption   string
	credit    string
	full      PhotoInfo
	thumbnail PhotoInfo
}

type Article struct {
	headline    string
	subheadline string
	section     string
	subsection  string
	source      string
	summary     string
	created_at  time.Time
	url         string
	photo       Photo
}

func getUrl(url string) {
	log.Print(fmt.Sprintf("Fetching %s", url))

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

		article := Article{
			headline:    articleJson.Get("headline").MustString(),
			subheadline: articleJson.Get("attrs").Get("brief").MustString(),
			section:     ssts.Get("section").MustString(),
			subsection:  ssts.Get("subsection").MustString(),
			source:      site,
			summary:     articleJson.Get("summary").MustString(),
			created_at:  time.Now(),
			url:         fmt.Sprintf("http://%s.com%s", site, articleJson.Get("url").MustString()),
			photo: Photo{
				caption: photoAttrs.Get("caption").MustString(),
				credit:  photoAttrs.Get("credit").MustString(),
				full: PhotoInfo{
					url:    strings.Join([]string{photoAttrs.Get("publishurl").MustString(), photoAttrs.Get("basename").MustString()}, ""),
					width:  photoAttrs.Get("oimagewidth").MustInt(),
					height: photoAttrs.Get("oimageheight").MustInt(),
				},
				thumbnail: PhotoInfo{
					url:    "TODO",
					width:  photoAttrs.Get("simagewidth").MustInt(),
					height: photoAttrs.Get("simageheight").MustInt(),
				},
			},
		}

		log.Print(article)
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

	for i := 0; i < 1; /*len(urls)*/ i++ {
		go getUrl(urls[i])
	}
}
