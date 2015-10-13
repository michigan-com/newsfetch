package link_parsing

import (
	"fmt"
	"net/url"
	"strings"

	gq "github.com/PuerkitoBio/goquery"
)

func BuildSearchURL(term string, page int) string {
	return fmt.Sprintf("http://freep.com/search/%s/%d/?ajax=true", url.QueryEscape(term), page)
}

func ExtractArticleURLsFromDocument(doc *gq.Document) []string {
	links := doc.Find("a.search-result-item-link[href]")

	urls := make([]string, 0, links.Length())
	links.Each(func(i int, link *gq.Selection) {
		url, exists := link.Attr("href")
		if exists {
			url = strings.TrimSpace(url)
			if len(url) > 0 {
				if !strings.Contains(url, "://") {
					if !strings.HasPrefix(url, "/") {
						url = "/" + url
					}
					url = "http://freep.com" + url
				}

				urls = append(urls, url)
			}
		}
	})

	return urls
}
