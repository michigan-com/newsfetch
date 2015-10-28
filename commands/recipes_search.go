package commands

import (
	"fmt"
	"time"

	"github.com/michigan-com/newsfetch/extraction"
	r "github.com/michigan-com/newsfetch/fetch/recipe"
	m "github.com/michigan-com/newsfetch/model"
	"github.com/spf13/cobra"
)

var searchOpts struct {
	pages   int
	onlyNew bool
}

type RecipeStats struct {
	RecipeCount int

	ArticlesWithoutRecipesCount        int
	articlesProcessedCount             int
	articlesIgnoredAsDuplicatesCount   int
	articlesIgnoredAsExistingCount     int
	articlesIgnoredInWrongSectionCount int

	URLsWithoutRecipes []string
}

func (lhs *RecipeStats) merge(rhs RecipeStats) {
	lhs.RecipeCount += rhs.RecipeCount
	lhs.ArticlesWithoutRecipesCount += rhs.ArticlesWithoutRecipesCount
	lhs.articlesProcessedCount += rhs.articlesProcessedCount
	lhs.articlesIgnoredAsDuplicatesCount += rhs.articlesIgnoredAsDuplicatesCount
	lhs.articlesIgnoredAsExistingCount += rhs.articlesIgnoredAsExistingCount
	lhs.articlesIgnoredInWrongSectionCount += rhs.articlesIgnoredInWrongSectionCount
	lhs.URLsWithoutRecipes = append(lhs.URLsWithoutRecipes, rhs.URLsWithoutRecipes...)
}

func (r RecipeStats) String() string {
	return fmt.Sprintf("%d recipes in %d articles, plus %d articles without recipes; total %d articles processed, %d ignored (%d dups, %d existing, %d wrong section)", r.RecipeCount, (r.articlesProcessedCount - r.ArticlesWithoutRecipesCount), r.ArticlesWithoutRecipesCount, r.articlesProcessedCount, r.articlesIgnoredAsDuplicatesCount+r.articlesIgnoredAsExistingCount+r.articlesIgnoredInWrongSectionCount, r.articlesIgnoredAsDuplicatesCount, r.articlesIgnoredAsExistingCount, r.articlesIgnoredInWrongSectionCount)
}

var cmdExtractRecipiesFromSearch = &cobra.Command{
	Use:   "process-search",
	Short: "Extract recipes using the search API",
	Run: func(cmd *cobra.Command, args []string) {
		startTime = time.Now()
		overallStats := RecipeStats{}

		page := 1
		total := 0
		processedURLsTable := make(map[string]bool, searchOpts.pages*10)
		for {
			urls, err := extraction.ExtractArticleURLsFromSearchResults("recipe", page)
			if err != nil {
				panic(err)
			}

			filteredURLs := m.FilterArticleURLsBySection(urls, "life")
			unprocessedFilteredURLs := filterUnprocessed(filteredURLs, processedURLsTable)
			processableURLs := unprocessedFilteredURLs

			if searchOpts.onlyNew {
				if globalConfig.MongoUrl == "" {
					panic("Need a MongoDB URI to run with --only-new")
				}

				existingURLs, err := CheckRecipeURLs(globalConfig.MongoUrl, unprocessedFilteredURLs)
				if err != nil {
					panic(err)
				}

				existingURLsTable := makeTableFromStrings(existingURLs)

				processableURLs = make([]string, 0, len(unprocessedFilteredURLs))
				for _, url := range unprocessedFilteredURLs {
					if !existingURLsTable[url] {
						processableURLs = append(processableURLs, url)
					}
				}

				fmt.Printf("Found existing recipes for %d of %d URLs. Will process %d URLs.\n", len(existingURLs), len(unprocessedFilteredURLs), len(processableURLs))
				for _, url := range existingURLs {
					fmt.Printf("Existing:     %s\n", url)
				}
				for _, url := range processableURLs {
					fmt.Printf("Will process: %s\n", url)
				}
			}

			result := r.DownloadRecipesFromUrls(processableURLs)

			pageStats := RecipeStats{}
			pageStats.RecipeCount = len(result.Recipes)
			pageStats.ArticlesWithoutRecipesCount = len(result.URLsWithoutRecipes)
			pageStats.articlesProcessedCount = len(result.URLs)
			pageStats.articlesIgnoredAsDuplicatesCount = len(filteredURLs) - len(unprocessedFilteredURLs)
			pageStats.articlesIgnoredAsExistingCount = len(unprocessedFilteredURLs) - len(processableURLs)
			pageStats.articlesIgnoredInWrongSectionCount = len(urls) - len(filteredURLs)
			pageStats.URLsWithoutRecipes = result.URLsWithoutRecipes

			overallStats.merge(pageStats)

			for _, url := range unprocessedFilteredURLs {
				processedURLsTable[url] = true
			}

			fmt.Printf("\nPage %d: %s\n", page, pageStats)
			fmt.Printf("Totals: %s\n", overallStats)

			fmt.Printf("\nSo far, got %d URLs without recipes:\n", len(overallStats.URLsWithoutRecipes))
			for idx, url := range overallStats.URLsWithoutRecipes {
				fmt.Printf("%03d) %s\n", idx+1, url)
			}

			if globalConfig.MongoUrl != "" {
				err := r.SaveRecipes(globalConfig.MongoUrl, result.Recipes)
				if err != nil {
					panic(err)
				}
				fmt.Printf("Saved %d recipes.\n", len(result.Recipes))
			}

			total += len(urls)

			if len(urls) == 0 {
				break
			}

			page++
			if page > searchOpts.pages {
				break
			}
		}

		getElapsedTime(&startTime)
	},
}

func init() {
	cmdExtractRecipiesFromSearch.Flags().IntVarP(&searchOpts.pages, "pages", "n", -1, "Specify the number of pages to process")
	cmdExtractRecipiesFromSearch.Flags().BoolVarP(&searchOpts.onlyNew, "only-new", "O", false, "Only process articles not already in the database")
}

func makeTableFromStrings(items []string) map[string]bool {
	table := make(map[string]bool, len(items))

	for _, item := range items {
		table[item] = true
	}

	return table
}
