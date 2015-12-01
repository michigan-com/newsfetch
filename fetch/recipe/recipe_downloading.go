package fetch

import (
	"fmt"

	"github.com/michigan-com/newsfetch/extraction"
	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model"
	"github.com/michigan-com/newsfetch/model/recipetypes"
)

var recipeDebugger = lib.NewCondLogger("newsfetch:fetch:recipe")

func DownloadAndSaveRecipesForArticles(mongoUrl string, articles []*m.Article) error {
	for _, article := range articles {
		err := DownloadAndSaveRecipesForArticle(mongoUrl, article)
		if err != nil {
			return err
		}
	}
	return nil
}

func DownloadAndSaveRecipesForArticle(mongoUrl string, article *m.Article) error {
	recipes := DownloadRecipesForArticle(article)

	if mongoUrl != "" {
		err := SaveRecipes(mongoUrl, recipes)
		return err
	} else {
		return nil
	}
}

type DownloadRecipesResult struct {
	Recipes            []*recipetypes.Recipe
	URLs               []string
	URLsWithoutRecipes []string
}

func DownloadRecipesForArticle(article *m.Article) []*recipetypes.Recipe {
	return DownloadRecipesFromUrls([]string{article.Url}).Recipes
}

func DownloadRecipesFromUrls(urls []string) DownloadRecipesResult {
	result := DownloadRecipesResult{}
	result.URLs = make([]string, 0)
	result.URLsWithoutRecipes = make([]string, 0)

	visited := make(map[string]bool)
	for len(urls) > 0 {
		url := urls[0]
		urls = urls[1:]

		if visited[url] {
			continue
		}
		visited[url] = true

		articleId := lib.GetArticleId(url)
		if articleId < 1 {
			recipeDebugger.Println("Skipped, cannot determine article ID")
			continue
		}

		extracted := extraction.ExtractDataFromHTMLAtURL(url, false)

		recipesInArticle := extracted.RecipeData.Recipes

		for _, recipe := range recipesInArticle {
			recipe.ArticleId = articleId
		}

		if recipeDebugger.IsEnabled() {
			fmt.Printf("Found %d recipes + %d links in %s\n", len(recipesInArticle), len(extracted.RecipeData.EmbeddedArticleUrls), url)
		}

		if len(recipesInArticle) == 0 && len(extracted.RecipeData.EmbeddedArticleUrls) == 0 {
			result.URLsWithoutRecipes = append(result.URLsWithoutRecipes, url)
		}
		result.URLs = append(result.URLs, url)

		if false {
			for i, recipe := range result.Recipes {
				recipeDebugger.Println()
				recipeDebugger.Println("Recipe ", i, "=", recipe.String())
				recipeDebugger.Println()
			}
		}

		result.Recipes = append(result.Recipes, extracted.RecipeData.Recipes...)
		urls = append(urls, extracted.RecipeData.EmbeddedArticleUrls...)
	}

	return result
}
