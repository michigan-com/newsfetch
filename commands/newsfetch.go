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
	includeTitle bool
	noprompt     bool
	startTime    time.Time
	VERSION      string
	COMMITHASH   string
	loop         int
	noUpdate     bool
	timeLogger   = lib.NewCondLogger("timer")
	NewsfetchCmd = &cobra.Command{Use: "newsfetch"}
	url          = "http://www.freep.com/story/news/local/michigan/2015/08/06/farid-fata-cancer-sentencing/31213475/"
	//w            = new(tabwriter.Writer)
)

func Execute(ver, commit string) {
	VERSION = ver
	COMMITHASH = commit
	loadConfig()
	AddCommands()
	AddFlags()

	NewsfetchCmd.Execute()
}

func AddFlags() {
	cmdArticle.Flags().StringVarP(&articleUrl, "url", "u", url, "URL of Gannett article")

	cmdGetArticles.Flags().StringVarP(&siteStr, "sites", "i", "all", "Comma separated list of Gannett sites to fetch articles from")
	cmdGetArticles.Flags().StringVarP(&sectionStr, "sections", "e", "all", "Comma separated list of article sections to fetch from")

	cmdBody.Flags().StringVarP(&articleUrl, "url", "u", url, "URL of Gannett article")
	cmdBody.Flags().BoolVarP(&includeTitle, "title", "t", false, "Place title of article on the first line of output")

	cmdChartbeat.PersistentFlags().IntVarP(&loop, "loop", "l", -1, "Specify the internval in seconds to loop the fetching of the toppages api")
	cmdChartbeat.PersistentFlags().BoolVarP(&noUpdate, "no-update", "n", false, "If present, mapi will not be updated")
}

func AddCommands() {
	cmdArticles.AddCommand(cmdGetArticles)
	cmdArticles.AddCommand(cmdCopyArticles)

	cmdChartbeat.AddCommand(cmdTopPages)
	cmdChartbeat.AddCommand(cmdQuickStats)
	cmdChartbeat.AddCommand(cmdTopGeo)
	cmdChartbeat.AddCommand(cmdReferrers)
	cmdChartbeat.AddCommand(cmdRecent)
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
	timeLogger.Println("Total time to run: ", endTime.Sub(*sTime))
}
