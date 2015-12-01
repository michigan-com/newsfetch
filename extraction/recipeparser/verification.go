package recipeparser

// import (
// 	"github.com/michigan-com/newsfetch/recipes/recipestats"
// )

// func VerifyRecipe(recipe *Recipe, stats *recipestats.Stats, threshold int) {
// 	results := ClassifyRecipeParagraphs(recipe)

// 	totalCount := 0
// 	matchedCount := 0

// 	for _, result := range results {
// 		if result.DetectedRole == IngredientHeader {
// 			continue
// 		}

// 		totalCount++

// 		if result.DetectedRole == Ingredient && result.Confidence >= rp.Likely {
// 			matchedCount++
// 		}
// 	}

// 	ranking := totalCount - matchedCount

// 	for _, result := range results {
// 		if result.DetectedRole == Ingredient {
// 			stats.AddIngredient(result.CanonicalText, ranking)
// 		}
// 	}

// 	for _, result := range results {
// 		if result.DetectedRole == Conflict {

// 		}
// 	}

// 	if (threshold < 0) || (ranking < threshold) {
// 		for _, result := range results {
// 		}
// 	}
// }

// func ClassifyRecipeParagraphs(recipe *Recipe, matcher *rp.Matcher) []ParagraphResult {
// 	var results []ParagraphResult

// 	for _, ingred := range recipe.Ingredients {
// 		result := matcher.ClassifyParagraph(ingred.Text)
// 		result.AssignedRole = Ingredient
// 		results = append(results, result)
// 	}

// 	for _, dir := range recipe.Instructions {
// 		result := matcher.ClassifyParagraph(dir.Text)
// 		result.AssignedRole = Direction
// 		results = append(results, result)
// 	}

// 	return results
// }
