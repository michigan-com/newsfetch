package lib

import (
	"strings"

	gq "github.com/PuerkitoBio/goquery"
)

type RecipeExtractionResult struct {
	Recipes             []*Recipe
	UnusedParagraphs    []RecipeFragment
	EmbeddedArticleUrls []string
}

type recipeState int

const (
	none recipeState = iota
	unconfirmed
	confirmed
)

func ExtractRecipes(doc *gq.Document) RecipeExtractionResult {
	var embeddedArticleUrls []string
	doc.Find(".story-asset.oembed-asset a").Each(func(_ int, s *gq.Selection) {
		url, exists := s.Attr("href")
		if exists {
			if IsPotentialRecipeUrl(url) {
				embeddedArticleUrls = append(embeddedArticleUrls, url)
			}
		}
	})

	var photo *RecipePhoto
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
			photo = new(RecipePhoto)
			if exists {
				photo.FullSizeImage = &RecipeImage{Url: url}
			}
			if smallExists {
				photo.SmallImage = &RecipeImage{Url: smallUrl}
			} else if srcExists {
				photo.SmallImage = &RecipeImage{Url: srcUrl}
			}
		}
	})

	fragments := extractRecipeFragments(doc)

	recipes := make([]*Recipe, 0)

	state := none
	var recipeFragments []RecipeFragment
	var unused []RecipeFragment

	setState := func(newState recipeState) {
		// println("setState", newState)
		if newState == state {
			return
		}

		if state == none {
			recipeFragments = make([]RecipeFragment, 0)
		}
		if newState == none {
			if state == confirmed {
				recipe := NewRecipe()
				recipe.Photo = photo
				walkFragments(recipeFragments)
				for _, fragment := range recipeFragments {
					fragment.AddToRecipe(recipe)
				}
				recipes = append(recipes, recipe)
			}
			recipeFragments = nil
		}

		state = newState
	}

	for _, fragment := range fragments {
		tag := fragment.Tag()
		// println("fragment", tag)

		switch tag {
		case TitleTag:
			setState(none)
			setState(unconfirmed)

		case PossibleTitleTag:
			if state == none {
				setState(unconfirmed)
			}

		case TimingTag, NutritionDataTag, SignatureTag, IngredientsSubtitleTag, DirectionsSubtitleTag:
			if state == unconfirmed {
				setState(confirmed)
			}

		case PossibleIngredientTag:
			if state == confirmed {
				fragment.Mark(IngredientTag)
			}

		case PossibleIngredientSubdivisionTag:
			if state != none {
				fragment.Mark(IngredientTag)
			}

		case EndMarkerTag:
			setState(none)
		}

		if state != none {
			recipeFragments = append(recipeFragments, fragment)
		} else {
			unused = append(unused, fragment)
		}
	}
	setState(none)

	return RecipeExtractionResult{Recipes: recipes, UnusedParagraphs: unused, EmbeddedArticleUrls: embeddedArticleUrls}
}

func walkFragments(fragments []RecipeFragment) {
	for {
		changed := false

		// PossibleIngredientSubdivisionTag before an ingredient becomes IngredientSubdivisionTag
		walkFragmentsBackward(fragments, func(cur RecipeFragment, curTag RecipeFragmentTag, nextTag RecipeFragmentTag) {
			switch curTag {
			case PossibleIngredientSubdivisionTag:
				switch nextTag {
				case IngredientTag, PossibleIngredientTag:
					cur.Mark(IngredientSubdivisionTag)
					changed = true
				}
			}
		})

		// ShortParagraphTag before and after an ingredient becomes IngredientTag
		walkFragmentsBackward(fragments, func(cur RecipeFragment, curTag RecipeFragmentTag, nextTag RecipeFragmentTag) {
			switch curTag {
			case ShortParagraphTag:
				switch nextTag {
				case IngredientTag, PossibleIngredientTag, IngredientSubdivisionTag:
					cur.Mark(IngredientTag)
					changed = true
				}
			}
		})
		walkFragmentsForward(fragments, func(cur RecipeFragment, curTag RecipeFragmentTag, prevTag RecipeFragmentTag) {
			switch curTag {
			case ShortParagraphTag:
				switch prevTag {
				case IngredientTag, PossibleIngredientTag, IngredientSubdivisionTag:
					cur.Mark(IngredientTag)
					changed = true
				}
			}
		})

		if !changed {
			break
		}
	}
}

func walkFragmentsForward(fragments []RecipeFragment, iter func(cur RecipeFragment, curTag RecipeFragmentTag, prevTag RecipeFragmentTag)) {
	prevTag := NoTag
	for _, cur := range fragments {
		curTag := cur.Tag()
		iter(cur, curTag, prevTag)
		prevTag = cur.Tag() // might have changed
	}
}

func walkFragmentsBackward(fragments []RecipeFragment, iter func(cur RecipeFragment, curTag RecipeFragmentTag, nextTag RecipeFragmentTag)) {
	nextTag := NoTag
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
