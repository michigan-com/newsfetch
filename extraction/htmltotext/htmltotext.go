package htmltotext

import (
	gq "github.com/PuerkitoBio/goquery"

	t "github.com/michigan-com/newsfetch/model/bodytypes"
)

func ConvertDocument(doc *gq.Document) []t.Paragraph {
	paragraphs := doc.Find("div[itemprop=articleBody] > p, div[itemprop=articleBody] li")
	if paragraphs.Length() == 0 {
		paragraphs = doc.Find("body > p")
	}

	result := make([]t.Paragraph, 0, paragraphs.Length())

	for _, node := range paragraphs.Nodes {
		para := ConvertParagraphNode(node)
		result = append(result, para)
	}

	return result
}
