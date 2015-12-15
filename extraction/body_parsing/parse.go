package body_parsing

import (
	"strconv"
	"strings"
	"time"

	gq "github.com/PuerkitoBio/goquery"
	"github.com/michigan-com/newsfetch/extraction/classify"
	"github.com/michigan-com/newsfetch/extraction/dateline"
	"github.com/michigan-com/newsfetch/extraction/recipe_parsing"
	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model"
)

var Debugger = lib.NewCondLogger("newsfetch:extraction:body_parsing")

func withoutEmptyStrings(strings []string) []string {
	result := make([]string, 0, len(strings))
	for _, el := range strings {
		if el != "" {
			result = append(result, el)
		}
	}
	return result
}

func ExtractBodyFromDocument(doc *gq.Document, fromJSON bool, includeTitle bool) *m.ExtractedBody {
	msg := new(m.Messages)

	var paragraphs *gq.Selection
	if fromJSON {
		paragraphs = doc.Find("p")
	} else {
		if len(doc.Find(".longform-body").Nodes) == 0 {
			paragraphs = doc.Find("div[itemprop=articleBody] > p")
		} else {
			paragraphs = doc.Find("div[itemprop=articleBody] > .longform-body > p")
		}
	}

	// remove contact info at the end of the article (might not be needed any more when parsing
	// HTML from JSON?)
	paragraphs.Find("span.-newsgate-paragraph-cci-endnote-contact-").Remove()
	paragraphs.Find("span.-newsgate-paragraph-cci-endnote-contrib-").Remove()

	ignoreRemaining := false
	paragraphStrings := paragraphs.Map(func(i int, paragraph *gq.Selection) string {
		if ignoreRemaining {
			return ""
		}
		for _, selector := range [...]string{"span.-newsgate-character-cci-tagline-name-", "span.-newsgate-paragraph-cci-infobox-head-"} {
			if el := paragraph.Find(selector); el.Length() > 0 {
				ignoreRemaining = true
				return ""
			}
		}

		text := strings.TrimSpace(paragraph.Text())

		if worthy, _ := classify.IsWorthyParagraph(text); !worthy {
			return ""
		}

		//marker := ""

		for _, selector := range [...]string{"span.-newsgate-paragraph-cci-subhead-lead-", "span.-newsgate-paragraph-cci-subhead-"} {
			if el := paragraph.Find(selector); el.Length() > 0 {
				//marker = "### "
				return ""
				break
			}
		}

		return text
	})

	if len(paragraphStrings) > 0 {
		paragraphStrings[0] = dateline.RmDateline(paragraphStrings[0])
	}

	content := make([]string, 0, len(paragraphStrings)+1)
	if includeTitle {
		title := ExtractTitleFromDocument(doc)
		content = append(content, title)
	}

	content = append(content, withoutEmptyStrings(paragraphStrings)...)

	body := strings.Join(content, "\n")
	recipeData, recipeMsg := recipe_parsing.ExtractRecipes(doc)
	msg.AddMessages("recipes", recipeMsg)
	extracted := m.ExtractedBody{body, recipeData, msg}
	return &extracted
}

func ExtractTitleFromDocument(doc *gq.Document) string {
	title := doc.Find("h1[itemprop=headline]").Text()

	if title == "" {
		title = doc.Find("h1[itemprop=name]").Text()
	}
	return strings.TrimSpace(title)
}

func ExtractSubheadlineFromDocument(doc *gq.Document) (subheadline string) {
	content, exists := doc.Find("meta[itemprop=description]").Attr("content")
	if !exists {
		return
	}
	subheadline = content
	return
}

func ExtractSectionInfo(doc *gq.Document) (sectionInfo *m.ExtractedSection) {
	sectionInfo = &m.ExtractedSection{}
	sectionString, exists := doc.Find("meta[itemprop=articleSection]").Attr("content")
	if !exists {
		return
	}

	sections := strings.Split(sectionString, ",")
	if len(sections) == 0 {
		return
	}

	section := sections[0]
	subsection := ""
	if len(sections) > 1 {
		subsection = sections[1]
	}
	sectionInfo = &m.ExtractedSection{
		section,
		subsection,
		sections,
	}
	return
}

func ExtractTimestamp(doc *gq.Document) (timestamp time.Time) {
	timestamp = time.Now()

	timeString := doc.Find("asset-metabar-time").Text()

	parsedTime, err := time.Parse("11:15 a.m. EST November 20, 2015", timeString)
	if err != nil {
		return
	}
	timestamp = parsedTime
	return
}

func ExtractPhotoInfo(doc *gq.Document) (photo *m.Photo) {
	// TODO figure out a way to make the following DOM query work. og:image throws an error
	// because of the colon
	//
	//ogImage, ogImageExists := doc.Find("meta[property=og:image]").Attr("content")
	//thumbImage, thumbImageExitsts := doc.Find("meta[property=thumbnailUrl]").Attr("content")
	var ogImage, thumbImage string
	var ogWidth, ogHeight, thumbWidth, thumbHeight int
	doc.Find("head meta").Each(func(index int, meta *gq.Selection) {

		// Yay for consistency. Have to check property and itemprop
		property, propExists := meta.Attr("property")
		if !propExists {
			property, propExists = meta.Attr("itemprop")
			if !propExists {
				return
			}
		}

		switch property {
		case "og:image":
			ogImage, _ = meta.Attr("content")
		case "og:image:width":
			widthText, _ := meta.Attr("content")
			ogWidth, _ = strconv.Atoi(widthText)
		case "og:image:height":
			heightText, _ := meta.Attr("content")
			ogHeight, _ = strconv.Atoi(heightText)
		case "thumbnailUrl":
			thumbImage, _ = meta.Attr("content")
		case "thumbnailWidth":
			widthText, _ := meta.Attr("content")
			thumbWidth, _ = strconv.Atoi(widthText)
		case "thumbnailHeight":
			heightText, _ := meta.Attr("content")
			thumbHeight, _ = strconv.Atoi(heightText)
		}
	})
	caption := doc.Find(".cutline").Text()
	credit := doc.Find(".credit").Text()

	if ogImage == "" || thumbImage == "" {
		return
	}

	photo = &m.Photo{}
	photo.Caption = caption
	photo.Credit = credit
	if ogImage != "" {
		photo.Full.Url = ogImage
		photo.Full.Width = ogWidth
		photo.Full.Height = ogHeight
	}
	if thumbImage != "" {
		photo.Thumbnail.Url = thumbImage
		photo.Thumbnail.Width = thumbWidth
		photo.Thumbnail.Height = thumbHeight
	}

	return
}
