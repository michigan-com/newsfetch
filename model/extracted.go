package model

import (
	"fmt"

	"github.com/michigan-com/newsfetch/model/recipetypes"
	"github.com/michigan-com/newsfetch/util/messages"
)

type RecipeExtractionResult struct {
	Recipes             []*recipetypes.Recipe
	UnusedParagraphs    []recipetypes.RecipeFragment
	EmbeddedArticleUrls []string
}

type ExtractedBody struct {
	Text       string
	RecipeData RecipeExtractionResult
	Messages   *messages.Messages
}

// TODO only use the Sections array for this, don't have section/subsection distinction
type ExtractedSection struct {
	Section    string
	Subsection string
	Sections   []string
}

func (e *ExtractedBody) String() string {
	return fmt.Sprintf("<ExtractedBody %s\n %s>\n", e.Text, e.Messages)
}
