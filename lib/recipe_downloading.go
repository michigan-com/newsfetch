package lib

func DownloadAndSaveRecipesForArticles(mongoUrl string, articles []*Article) error {
	for _, article := range articles {
		err := DownloadAndSaveRecipesForArticle(mongoUrl, article)
		if err != nil {
			return err
		}
	}
	return nil
}

func DownloadAndSaveRecipesForArticle(mongoUrl string, article *Article) error {
	recipes := DownloadRecipesForArticle(article)

	if mongoUrl != "" {
		err := SaveRecipes(mongoUrl, recipes)
		return err
	} else {
		return nil
	}
}

func DownloadRecipesForArticle(article *Article) []*Recipe {
	return DownloadRecipesFromUrls([]string{article.Url})
}

func DownloadRecipesFromUrls(urls []string) []*Recipe {
	recipes := make([]*Recipe, 0)

	visited := make(map[string]bool)
	for len(urls) > 0 {
		url := urls[0]
		urls = urls[1:]

		if visited[url] {
			continue
		}
		visited[url] = true

		Debugger.Println("Recipe extraction for URL", url)

		articleId := GetArticleId(url)
		if articleId < 1 {
			Debugger.Println("Skipped, cannot determine article ID")
			continue
		}

		extracted := ExtractBodyFromURLDirectly(url, false)

		for _, recipe := range extracted.RecipeData.Recipes {
			recipe.ArticleId = articleId
		}

		Debugger.Println("  found", len(extracted.RecipeData.Recipes), "recipes")

		if false {
			for i, recipe := range recipes {
				Debugger.Println()
				Debugger.Println("Recipe ", i, "=", recipe.String())
				Debugger.Println()
			}
		}

		recipes = append(recipes, extracted.RecipeData.Recipes...)
		urls = append(urls, extracted.RecipeData.EmbeddedArticleUrls...)
	}

	return recipes
}
