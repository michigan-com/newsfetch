package commands

import (
	//"text/tabwriter"
	"time"

	"github.com/michigan-com/newsfetch/lib"
	"github.com/spf13/cobra"
)

var (
	articleUrl   string
	siteStr      string
	sectionStr   string
	title        string
	output       bool
	timeit       bool
	body         bool
	verbose      bool
	includeTitle bool
	noprompt     bool
	startTime    time.Time
	VERSION      string
	COMMITHASH   string
	loop         int
	//w            = new(tabwriter.Writer)
)

var NewsfetchCmd = &cobra.Command{
	Use: "newsfetch",
}

var url = "http://www.freep.com/story/news/local/michigan/2015/08/06/farid-fata-cancer-sentencing/31213475/"

func Execute(ver, commit string) {
	VERSION = ver
	COMMITHASH = commit
	loadConfig()
	AddCommands()
	AddFlags()

	NewsfetchCmd.Execute()
}

func AddFlags() {
	NewsfetchCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	NewsfetchCmd.PersistentFlags().BoolVarP(&output, "output", "o", true, "Outputs results of command")
	NewsfetchCmd.PersistentFlags().BoolVarP(&timeit, "time", "m", false, "Outputs how long a command takes to finish")

	cmdArticle.Flags().StringVarP(&articleUrl, "url", "u", url, "URL of Gannett article")
	cmdArticle.Flags().BoolVarP(&body, "body", "b", false, "Fetches the article body content")

	cmdGetArticles.Flags().StringVarP(&siteStr, "sites", "i", "all", "Comma separated list of Gannett sites to fetch articles from")
	cmdGetArticles.Flags().StringVarP(&sectionStr, "sections", "e", "all", "Comma separated list of article sections to fetch from")
	cmdGetArticles.Flags().BoolVarP(&body, "body", "b", false, "Fetches the article body content")

	cmdBody.Flags().StringVarP(&articleUrl, "url", "u", url, "URL of Gannett article")
	cmdBody.Flags().BoolVarP(&includeTitle, "title", "t", false, "Place title of article on the first line of output")

	cmdChartbeat.PersistentFlags().IntVarP(&loop, "loop", "l", -1, "Specify the internval in seconds to loop the fetching of the toppages api")
}

func AddCommands() {
	cmdArticles.AddCommand(cmdGetArticles)
	cmdArticles.AddCommand(cmdCopyArticles)

	cmdChartbeat.AddCommand(cmdTopPages)
	cmdChartbeat.AddCommand(cmdQuickStats)
	cmdChartbeat.AddCommand(cmdTopGeo)
	cmdChartbeat.AddCommand(cmdAllBeats)

	NewsfetchCmd.AddCommand(
		cmdArticle,
		cmdBody,
		cmdArticles,
		cmdVersion,
		cmdChartbeat,
	)

	cmdRecipes.AddCommand(cmdReprocessRecipies)
	cmdRecipes.AddCommand(cmdReprocessRecipeById)
	cmdRecipes.AddCommand(cmdExtractRecipiesFromUrl)
	cmdRecipes.AddCommand(cmdExtractRecipiesFromSearch)
	NewsfetchCmd.AddCommand(cmdRecipes)
}

/*func printArticleBrief(w *tabwriter.Writer, article *m.Article) {
	fmt.Fprintf(
		w, "%s\t%s\t%s\t%s\t%s\n", article.Source, article.Section,
		article.Headline, article.Url, article.Timestamp,
	)
	w.Flush()
}*/

func getElapsedTime(sTime *time.Time) {
	endTime := time.Now()
	lib.Logger.Println("Total time to run: ", endTime.Sub(*sTime))
}
