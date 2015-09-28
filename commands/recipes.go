package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/michigan-com/newsfetch/lib"
	"github.com/spf13/cobra"
)

func printRecipies(articles []*lib.Article) {
	for _, article := range articles {
		fmt.Printf("%s/%s/%s - %s - %s\n", article.Source, article.Section, article.Subsection, article.Headline, article.Url)
	}
}

var cmdRecipes = &cobra.Command{
	Use:   "recipes",
	Short: "Command for Gannett recipe articles",
}

var cmdReprocessRecipies = &cobra.Command{
	Use:   "reprocess-all",
	Short: "Re-extract recipes from all articles saved in Mongo",
	Run: func(cmd *cobra.Command, args []string) {
		if timeit {
			startTime = time.Now()
		}

		articles, err := lib.LoadArticles(globalConfig.MongoUrl)
		if err != nil {
			panic(err)
		}

		beforeCount := len(articles)
		articles = lib.FilterArticlesForRecipeExtraction(articles)

		println("Loaded", beforeCount, "articles including", len(articles), "in food subsection.")

		if output {
			printRecipies(articles)
		}

		for _, article := range articles {
			err := lib.DownloadAndSaveRecipesForArticle(globalConfig.MongoUrl, article)
			if err != nil {
				panic(err)
			}
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}

var cmdReprocessRecipeById = &cobra.Command{
	Use:   "reprocess-id",
	Short: "Re-process recipes with given article IDs (8-digit ints) from Mongo",
	Run: func(cmd *cobra.Command, args []string) {
		if timeit {
			startTime = time.Now()
		}

		for _, arg := range args {
			articleId, err := strconv.Atoi(arg)
			if err != nil {
				panic(err)
			}

			article, err := lib.LoadArticleById(globalConfig.MongoUrl, articleId)
			if err != nil {
				panic(err)
			}

			err = lib.DownloadAndSaveRecipesForArticle(globalConfig.MongoUrl, article)
			if err != nil {
				panic(err)
			}
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}

var cmdExtractRecipiesFromUrl = &cobra.Command{
	Use:   "process-url",
	Short: "Extract recipes from the given URL",
	Run: func(cmd *cobra.Command, args []string) {
		if timeit {
			startTime = time.Now()
		}

		recipes := lib.DownloadRecipesFromUrls(args)

		fmt.Printf("Found %d recipes.\n", len(recipes))
		for i, recipe := range recipes {
			fmt.Printf("Recipe #%d: %s\n", i, recipe.String())
		}

		if globalConfig.MongoUrl != "" {
			err := lib.SaveRecipes(globalConfig.MongoUrl, recipes)
			if err != nil {
				panic(err)
			}
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}
