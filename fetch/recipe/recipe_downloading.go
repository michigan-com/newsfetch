package fetch

import (
	"github.com/michigan-com/newsfetch/extraction"
	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model"
)

var recipeDebugger = lib.NewCondLogger("fetch-recipe")

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

func DownloadRecipesForArticle(article *m.Article) []*m.Recipe {
	return DownloadRecipesFromUrls([]string{article.Url})
}

func DownloadRecipesFromUrls(urls []string) []*m.Recipe {
	recipes := make([]*m.Recipe, 0)

	visited := make(map[string]bool)
	for len(urls) > 0 {
		url := urls[0]
		urls = urls[1:]

		if visited[url] {
			continue
		}
		visited[url] = true

		recipeDebugger.Println("Recipe extraction for URL", url)

		articleId := lib.GetArticleId(url)
		if articleId < 1 {
			recipeDebugger.Println("Skipped, cannot determine article ID")
			continue
		}

		extracted := extraction.ExtractDataFromHTMLAtURL(url, false)

		for _, recipe := range extracted.RecipeData.Recipes {
			recipe.ArticleId = articleId
		}

		recipeDebugger.Println("  found", len(extracted.RecipeData.Recipes), "recipes")

		if false {
			for i, recipe := range recipes {
				recipeDebugger.Println()
				recipeDebugger.Println("Recipe ", i, "=", recipe.String())
				recipeDebugger.Println()
			}
		}

		recipes = append(recipes, extracted.RecipeData.Recipes...)
		urls = append(urls, extracted.RecipeData.EmbeddedArticleUrls...)
	}

	return recipes
}
