package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/michigan-com/newsfetch/extraction"
	r "github.com/michigan-com/newsfetch/fetch/recipe"
	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model"
	"github.com/spf13/cobra"
)

var recipeDebugger = lib.NewCondLogger("recipes")

func printRecipies(articles []*m.Article) {
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

		articles, err := LoadArticles(globalConfig.MongoUrl)
		if err != nil {
			panic(err)
		}

		beforeCount := len(articles)
		articles = FilterArticlesForRecipeExtraction(articles)

		recipeDebugger.Printf("Loaded %d articles including %d in food subsection.", beforeCount, len(articles))

		if output {
			printRecipies(articles)
		}

		for _, article := range articles {
			err := r.DownloadAndSaveRecipesForArticle(globalConfig.MongoUrl, article)
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

			article, err := LoadArticleById(globalConfig.MongoUrl, articleId)
			if err != nil {
				panic(err)
			}

			err = r.DownloadAndSaveRecipesForArticle(globalConfig.MongoUrl, article)
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

		recipes := r.DownloadRecipesFromUrls(args)

		fmt.Printf("Found %d recipes.\n", len(recipes))
		for i, recipe := range recipes {
			fmt.Printf("Recipe #%d: %s\n", i, recipe.String())
		}

		if globalConfig.MongoUrl != "" {
			err := r.SaveRecipes(globalConfig.MongoUrl, recipes)
			if err != nil {
				panic(err)
			}
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}

var cmdExtractRecipiesFromSearch = &cobra.Command{
	Use:   "process-search",
	Short: "Extract recipes using the search API",
	Run: func(cmd *cobra.Command, args []string) {
		if timeit {
			startTime = time.Now()
		}

		page := 1
		total := 0
		for {
			println("Page", page)
			urls, err := extraction.ExtractArticleURLsFromSearchResults("recipe", page)
			if err != nil {
				panic(err)
			}

			for _, url := range urls {
				id := lib.GetArticleId(url)
				println("Found article", id, "at", url)
			}

			recipes := r.DownloadRecipesFromUrls(urls)
			fmt.Printf("Found %d recipes.\n", len(recipes))

			if globalConfig.MongoUrl != "" {
				err := r.SaveRecipes(globalConfig.MongoUrl, recipes)
				if err != nil {
					panic(err)
				}
				fmt.Printf("Saved %d recipes.\n", len(recipes))
			}

			total += len(urls)

			if len(urls) == 0 {
				break
			}
			page++
		}

		println("Total", total)

		// recipes := r.DownloadRecipesFromUrls(args)

		// fmt.Printf("Found %d recipes.\n", len(recipes))
		// for i, recipe := range recipes {
		// 	fmt.Printf("Recipe #%d: %s\n", i, recipe.String())
		// }

		// if globalConfig.MongoUrl != "" {
		// 	err := r.SaveRecipes(globalConfig.MongoUrl, recipes)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}

func LoadArticles(mongoUri string) ([]*m.Article, error) {
	session := lib.DBConnect(mongoUri)
	defer lib.DBClose(session)

	articleCol := session.DB("").C("Article")

	var result []*m.Article
	err := articleCol.Find(nil).All(&result)
	return result, err
}

func LoadArticleById(mongoUri string, articleId int) (*m.Article, error) {
	session := lib.DBConnect(mongoUri)
	defer lib.DBClose(session)

	articleCol := session.DB("").C("Article")

	var result *m.Article
	err := articleCol.Find(bson.M{"article_id": articleId}).One(&result)
	return result, err
}

func LoadRemoteArticles(url string) ([]*m.Article, error) {
	artDebugger.Println("Fetching ", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	artDebugger.Println(fmt.Sprintf("Successfully fetched %s", url))

	var response struct {
		Articles []m.Article `json:"articles"`
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&response)
	if err != nil {
		return nil, err
	}

	return sliceOfArticlesToSliceOfPointers(response.Articles), nil
}

func FilterArticlesBySubsection(articles []*m.Article, section string, subsection string) []*m.Article {
	result := make([]*m.Article, 0, len(articles))
	for _, el := range articles {
		if (el.Section == section) && (el.Subsection == subsection) {
			result = append(result, el)
		}
	}
	return result
}

func FilterArticlesForRecipeExtraction(articles []*m.Article) []*m.Article {
	return FilterArticlesBySubsection(articles, "life", "food")
}

func sliceOfArticlesToSliceOfPointers(articles []m.Article) []*m.Article {
	result := make([]*m.Article, 0, len(articles))
	for _, el := range articles {
		copy := el
		result = append(result, &copy)
	}
	return result
}
