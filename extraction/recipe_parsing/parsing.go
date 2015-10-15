package recipe_parsing

import (
	"strings"

	gq "github.com/PuerkitoBio/goquery"
	m "github.com/michigan-com/newsfetch/model"
)

type recipeState int

const (
	none recipeState = iota
	unconfirmed
	confirmed
)

func ExtractRecipes(doc *gq.Document) (m.RecipeExtractionResult, *m.Messages) {
	msg := new(m.Messages)

	var embeddedArticleUrls []string
	doc.Find(".story-asset.oembed-asset a").Each(func(_ int, s *gq.Selection) {
		url, exists := s.Attr("href")
		if exists {
			if IsPotentialRecipeUrl(url) {
				embeddedArticleUrls = append(embeddedArticleUrls, url)
			}
		}
	})

	var photo *m.RecipePhoto
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
			photo = new(m.RecipePhoto)
			if exists {
				photo.FullSizeImage = &m.RecipeImage{Url: url}
			}
			if smallExists {
				photo.SmallImage = &m.RecipeImage{Url: smallUrl}
			} else if srcExists {
				photo.SmallImage = &m.RecipeImage{Url: srcUrl}
			}
		}
	})

	fragments := extractRecipeFragments(doc, msg)

	markDefiniteRecipeTitles(fragments)
	if hasFragmentWithTag(fragments, m.TitleTag) {
		fragments = omitFragmentsBeforeFirstTitleTag(fragments)
	}

	fixupParagraphTexts(fragments)
	fragments = skipEmptyParagraphFragments(fragments)

	recipes := make([]*m.Recipe, 0)

	state := none
	var recipeFragments []m.RecipeFragment
	var nextRecipeFragments []m.RecipeFragment
	var unused []m.RecipeFragment

	setState := func(newState recipeState) {
		if newState == state {
			return
		}

		if state == none {
			recipeFragments = make([]m.RecipeFragment, 0)
		}
		if newState == none {
			if state == confirmed {
				recipe := m.NewRecipe()
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
		case m.TitleTag:
			setState(none)
			setState(unconfirmed)

		case m.PossibleTitleTag:
			if state == none {
				setState(unconfirmed)
			}

		case m.TimingTag, m.NutritionDataTag, m.SignatureTag, m.IngredientsSubtitleTag, m.DirectionsSubtitleTag:
			if state == unconfirmed {
				setState(confirmed)
			}

		case m.ServingSizeAltTimingTag:
			if state == none {
				nextRecipeFragments = append(nextRecipeFragments, fragment)
			}

		case m.PossibleIngredientTag:
			if state == confirmed {
				fragment.Mark(m.IngredientTag)
			}

		case m.PossibleIngredientSubdivisionTag:
			if state != none {
				fragment.Mark(m.IngredientTag)
			}

		case m.EndMarkerTag:
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

func walkFragments(fragments []m.RecipeFragment) {
	for {
		changed := false

		// PossibleIngredientSubdivisionTag before an ingredient becomes IngredientSubdivisionTag
		walkFragmentsBackward(fragments, func(cur m.RecipeFragment, curTag m.RecipeFragmentTag, nextTag m.RecipeFragmentTag) {
			switch curTag {
			case m.PossibleIngredientSubdivisionTag:
				switch nextTag {
				case m.IngredientTag, m.PossibleIngredientTag:
					cur.Mark(m.IngredientSubdivisionTag)
					changed = true
				}
			}
		})

		// ShortParagraphTag before and after an ingredient becomes IngredientTag
		walkFragmentsBackward(fragments, func(cur m.RecipeFragment, curTag m.RecipeFragmentTag, nextTag m.RecipeFragmentTag) {
			switch curTag {
			case m.ShortParagraphTag:
				switch nextTag {
				case m.IngredientTag, m.PossibleIngredientTag, m.IngredientSubdivisionTag:
					cur.Mark(m.IngredientTag)
					changed = true
				}
			}
		})
		walkFragmentsForward(fragments, func(cur m.RecipeFragment, curTag m.RecipeFragmentTag, prevTag m.RecipeFragmentTag) {
			switch curTag {
			case m.ShortParagraphTag:
				switch prevTag {
				case m.IngredientTag, m.PossibleIngredientTag, m.IngredientSubdivisionTag:
					cur.Mark(m.IngredientTag)
					changed = true
				}
			}
		})

		if !changed {
			break
		}
	}
}

func markDefiniteRecipeTitles(fragments []m.RecipeFragment) {
	walkFragmentsBackward(fragments, func(cur m.RecipeFragment, curTag m.RecipeFragmentTag, nextTag m.RecipeFragmentTag) {
		switch curTag {
		case m.PossibleTitleTag:
			switch nextTag {
			case m.TimingTag:
				cur.Mark(m.TitleTag)
			}
		}
	})
}

func hasFragmentWithTag(fragments []m.RecipeFragment, tag m.RecipeFragmentTag) bool {
	for _, frag := range fragments {
		if frag.Tag() == tag {
			return true
		}
	}
	return false
}

func omitFragmentsBeforeFirstTitleTag(fragments []m.RecipeFragment) []m.RecipeFragment {
	result := make([]m.RecipeFragment, 0, len(fragments))
	found := false
	for _, frag := range fragments {
		if frag.Tag() == m.TitleTag {
			found = true
		}
		if found {
			result = append(result, frag)
		}
	}
	return result
}

func skipEmptyParagraphFragments(fragments []m.RecipeFragment) []m.RecipeFragment {
	result := make([]m.RecipeFragment, 0, len(fragments))
	for _, frag := range fragments {
		if frag, ok := frag.(*m.ParagraphFragment); ok {
			if len(frag.Text) == 0 {
				continue
			}
		}
		result = append(result, frag)
	}
	return result
}

func fixupParagraphTexts(fragments []m.RecipeFragment) {
	for _, frag := range fragments {
		if frag, ok := frag.(*m.ParagraphFragment); ok {
			frag.Text = strings.TrimSpace(frag.Text)
			oldText := frag.Text
			frag.Text = strings.Trim(frag.Text, "■ ")

			switch frag.Tag() {
			case m.IngredientTag, m.PossibleIngredientTag:
			case m.ParagraphTag, m.ShortParagraphTag:
				if oldText != frag.Text {
					frag.Mark(m.IngredientTag)
				}
			}
		}
	}
}

func walkFragmentsForward(fragments []m.RecipeFragment, iter func(cur m.RecipeFragment, curTag m.RecipeFragmentTag, prevTag m.RecipeFragmentTag)) {
	prevTag := m.NoTag
	for _, cur := range fragments {
		curTag := cur.Tag()
		iter(cur, curTag, prevTag)
		prevTag = cur.Tag() // might have changed
	}
}

func walkFragmentsBackward(fragments []m.RecipeFragment, iter func(cur m.RecipeFragment, curTag m.RecipeFragmentTag, nextTag m.RecipeFragmentTag)) {
	nextTag := m.NoTag
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
