package recipeparser

import (
	"strings"

	gq "github.com/PuerkitoBio/goquery"

	"github.com/michigan-com/newsfetch/extraction/htmltotext"
	"github.com/michigan-com/newsfetch/extraction/recipematcher"
	m "github.com/michigan-com/newsfetch/model"
	bt "github.com/michigan-com/newsfetch/model/bodytypes"
	t "github.com/michigan-com/newsfetch/model/recipetypes"
	"github.com/michigan-com/newsfetch/util/messages"
)

type recipeState int

const (
	none recipeState = iota
	unconfirmed
	confirmed
)

func ExtractRecipes(doc *gq.Document) (m.RecipeExtractionResult, *messages.Messages) {
	msg := new(messages.Messages)

	var embeddedArticleUrls []string
	doc.Find(".story-asset.oembed-asset a").Each(func(_ int, s *gq.Selection) {
		url, exists := s.Attr("href")
		if exists {
			if IsPotentialRecipeUrl(url) {
				embeddedArticleUrls = append(embeddedArticleUrls, url)
			}
		}
	})

	var photo *t.RecipePhoto
	doc.Find("img[data-mycapture-src]").Each(func(_ int, s *gq.Selection) {
		url, exists := s.Attr("data-mycapture-src")
		if exists && len(url) == 0 {
			exists = false
		}

		smallUrl, smallExists := s.Attr("data-mycapture-sm-src")
		if smallExists && len(smallUrl) == 0 {
			smallExists = false
		}

		srcUrl, srcExists := s.Attr("src")
		if srcExists && len(srcUrl) == 0 {
			srcExists = false
		}

		if exists || smallExists || srcExists {
			photo = new(t.RecipePhoto)
			if exists {
				photo.FullSizeImage = &t.RecipeImage{Url: url}
			}
			if smallExists {
				photo.SmallImage = &t.RecipeImage{Url: smallUrl}
			} else if srcExists {
				photo.SmallImage = &t.RecipeImage{Url: srcUrl}
			}
		}
	})

	paragraphs := htmltotext.ConvertDocument(doc)

	recipeParagraphs := parseParagraphs(paragraphs, msg)

	fragments := convertRecipeParagraphsToFragments(recipeParagraphs)

	markDefiniteRecipeTitles(fragments)
	if hasFragmentWithTag(fragments, t.TitleTag) {
		fragments = omitFragmentsBeforeFirstTitleTag(fragments)
	}

	fixupParagraphTexts(fragments)
	fragments = skipEmptyParagraphFragments(fragments)

	recipes := make([]*t.Recipe, 0)

	state := none
	var recipeFragments []t.RecipeFragment
	var nextRecipeFragments []t.RecipeFragment
	var unused []t.RecipeFragment

	setState := func(newState recipeState) {
		if newState == state {
			return
		}

		if state == none {
			recipeFragments = make([]t.RecipeFragment, 0)
		}
		if newState == none {
			if state == confirmed {
				recipe := t.NewRecipe()
				recipe.Photo = photo
				walkFragments(recipeFragments)
				for _, fragment := range nextRecipeFragments {
					fragment.AddToRecipe(recipe)
				}
				for _, fragment := range recipeFragments {
					fragment.AddToRecipe(recipe)
				}
				recipes = append(recipes, recipe)
				nextRecipeFragments = nil
			}
			recipeFragments = nil
		}

		state = newState
	}

	for _, fragment := range fragments {
		tag := fragment.Tag()

		switch tag {
		case t.TitleTag:
			setState(none)
			setState(unconfirmed)

		case t.PossibleTitleTag:
			if state == none {
				setState(unconfirmed)
			}

		case t.TimingTag, t.NutritionDataTag, t.SignatureTag, t.IngredientsSubtitleTag, t.DirectionsSubtitleTag:
			if state == unconfirmed {
				setState(confirmed)
			}

		case t.ServingSizeAltTimingTag:
			if state == none {
				nextRecipeFragments = append(nextRecipeFragments, fragment)
			}

		case t.PossibleIngredientTag:
			if state == confirmed {
				fragment.Mark(t.IngredientTag)
			}

		case t.PossibleIngredientSubdivisionTag:
			if state != none {
				fragment.Mark(t.IngredientTag)
			}

		case t.EndMarkerTag:
			setState(none)
		}

		if state != none {
			recipeFragments = append(recipeFragments, fragment)
		} else {
			unused = append(unused, fragment)
		}
	}
	setState(none)

	return m.RecipeExtractionResult{Recipes: recipes, UnusedParagraphs: unused, EmbeddedArticleUrls: embeddedArticleUrls}, msg
}

func parseParagraphs(paragraphs []bt.Paragraph, msg *messages.Messages) []t.Paragraph {
	result := make([]t.Paragraph, 0, len(paragraphs))
	for _, el := range paragraphs {
		role, match, fragments := recipematcher.Match(el, msg)

		paragraph := t.Paragraph{el, role, match, fragments}
		result = append(result, paragraph)
	}
	return result
}

func walkFragments(fragments []t.RecipeFragment) {
	for {
		changed := false

		// PossibleIngredientSubdivisionTag before an ingredient becomes IngredientSubdivisionTag
		walkFragmentsBackward(fragments, func(cur t.RecipeFragment, curTag t.RecipeFragmentTag, nextTag t.RecipeFragmentTag) {
			switch curTag {
			case t.PossibleIngredientSubdivisionTag:
				switch nextTag {
				case t.IngredientTag, t.PossibleIngredientTag:
					cur.Mark(t.IngredientSubdivisionTag)
					changed = true
				}
			}
		})

		// ShortParagraphTag before and after an ingredient becomes IngredientTag
		walkFragmentsBackward(fragments, func(cur t.RecipeFragment, curTag t.RecipeFragmentTag, nextTag t.RecipeFragmentTag) {
			switch curTag {
			case t.ShortParagraphTag:
				switch nextTag {
				case t.IngredientTag, t.PossibleIngredientTag, t.IngredientSubdivisionTag:
					cur.Mark(t.IngredientTag)
					changed = true
				}
			}
		})
		walkFragmentsForward(fragments, func(cur t.RecipeFragment, curTag t.RecipeFragmentTag, prevTag t.RecipeFragmentTag) {
			switch curTag {
			case t.ShortParagraphTag:
				switch prevTag {
				case t.IngredientTag, t.PossibleIngredientTag, t.IngredientSubdivisionTag:
					cur.Mark(t.IngredientTag)
					changed = true
				}
			}
		})

		if !changed {
			break
		}
	}
}

func convertRecipeParagraphsToFragments(paragraphs []t.Paragraph) []t.RecipeFragment {
	fragments := make([]t.RecipeFragment, 0, len(paragraphs))
	for _, para := range paragraphs {
		var fragment t.RecipeFragment
		if len(para.Fragments) > 0 {
			fragment = para.Fragments[0].(t.RecipeFragment)
		} else {
			switch para.Role {
			case t.Conflict, t.ShitTail:

			case t.EndMarker:
				fragment = &t.RecipeMarkerFragment{TagF: t.EndMarkerTag}

			case t.Title:
				if para.Confidence >= bt.Likely {
					fragment = &t.ParagraphFragment{TagF: t.TitleTag, Text: para.Text}
				} else {
					fragment = &t.ParagraphFragment{TagF: t.PossibleTitleTag, Text: para.Text}
				}

			case t.IngredientsHeading:
				fragment = &t.ParagraphFragment{TagF: t.IngredientsSubtitleTag, Text: para.Text}
			case t.DirectionsHeading:
				fragment = &t.ParagraphFragment{TagF: t.DirectionsSubtitleTag, Text: para.Text}

			case t.IngredientSubsectionHeading:
				if para.Confidence >= bt.Likely {
					fragment = &t.ParagraphFragment{TagF: t.IngredientSubdivisionTag, Text: para.Text}
				} else {
					fragment = &t.ParagraphFragment{TagF: t.PossibleIngredientSubdivisionTag, Text: para.Text}

				}

			case t.Ingredient:
				if para.Confidence >= bt.Likely {
					fragment = &t.ParagraphFragment{TagF: t.IngredientTag, Text: para.Text}
				} else {
					fragment = &t.ParagraphFragment{TagF: t.PossibleIngredientTag, Text: para.Text}
				}

			case t.Direction:
				fragment = &t.ParagraphFragment{TagF: t.InstructionTag, Text: para.Text}

			default:
				panic("Unexpected paragraph type")
			}
		}

		if fragment != nil {
			fragments = append(fragments, fragment)
		}
	}
	return fragments
}

func markDefiniteRecipeTitles(fragments []t.RecipeFragment) {
	walkFragmentsBackward(fragments, func(cur t.RecipeFragment, curTag t.RecipeFragmentTag, nextTag t.RecipeFragmentTag) {
		switch curTag {
		case t.PossibleTitleTag:
			switch nextTag {
			case t.TimingTag:
				cur.Mark(t.TitleTag)
			}
		}
	})
}

func hasFragmentWithTag(fragments []t.RecipeFragment, tag t.RecipeFragmentTag) bool {
	for _, frag := range fragments {
		if frag.Tag() == tag {
			return true
		}
	}
	return false
}

func omitFragmentsBeforeFirstTitleTag(fragments []t.RecipeFragment) []t.RecipeFragment {
	result := make([]t.RecipeFragment, 0, len(fragments))
	found := false
	for _, frag := range fragments {
		if frag.Tag() == t.TitleTag {
			found = true
		}
		if found {
			result = append(result, frag)
		}
	}
	return result
}

func skipEmptyParagraphFragments(fragments []t.RecipeFragment) []t.RecipeFragment {
	result := make([]t.RecipeFragment, 0, len(fragments))
	for _, frag := range fragments {
		if frag, ok := frag.(*t.ParagraphFragment); ok {
			if len(frag.Text) == 0 {
				continue
			}
		}
		result = append(result, frag)
	}
	return result
}

func fixupParagraphTexts(fragments []t.RecipeFragment) {
	for _, frag := range fragments {
		if frag, ok := frag.(*t.ParagraphFragment); ok {
			frag.Text = strings.TrimSpace(frag.Text)
			oldText := frag.Text
			frag.Text = strings.Trim(frag.Text, "■ ")

			switch frag.Tag() {
			case t.IngredientTag, t.PossibleIngredientTag:
			case t.ParagraphTag, t.ShortParagraphTag:
				if oldText != frag.Text {
					frag.Mark(t.IngredientTag)
				}
			}
		}
	}
}

func walkFragmentsForward(fragments []t.RecipeFragment, iter func(cur t.RecipeFragment, curTag t.RecipeFragmentTag, prevTag t.RecipeFragmentTag)) {
	prevTag := t.NoTag
	for _, cur := range fragments {
		curTag := cur.Tag()
		iter(cur, curTag, prevTag)
		prevTag = cur.Tag() // might have changed
	}
}

func walkFragmentsBackward(fragments []t.RecipeFragment, iter func(cur t.RecipeFragment, curTag t.RecipeFragmentTag, nextTag t.RecipeFragmentTag)) {
	nextTag := t.NoTag
	for i := len(fragments) - 1; i >= 0; i-- {
		cur := fragments[i]
		curTag := cur.Tag()
		iter(cur, curTag, nextTag)
		nextTag = cur.Tag() // might have changed
	}
}

func IsPotentialRecipeUrl(url string) bool {
	return strings.Contains(url, "freep.com") && strings.Contains(url, "life/food")
}
