package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	mongoUri     string
	articleUrl   string
	siteStr      string
	sectionStr   string
	title        string
	apiKey       string
	output       bool
	timeit       bool
	body         bool
	verbose      bool
	includeTitle bool
	noprompt     bool
	startTime    time.Time
	VERSION      string
	loop         int
)

var NewsfetchCmd = &cobra.Command{
	Use: "newsfetch",
}

var url = "http://www.freep.com/story/news/local/michigan/2015/08/06/farid-fata-cancer-sentencing/31213475/"

func Execute(ver string) {
	VERSION = ver
	AddCommands()
	AddFlags()

	NewsfetchCmd.Execute()
}

func AddFlags() {
	NewsfetchCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	NewsfetchCmd.PersistentFlags().BoolVarP(&output, "output", "o", true, "Outputs results of command")
	NewsfetchCmd.PersistentFlags().BoolVarP(&timeit, "time", "m", false, "Outputs how long a command takes to finish")

	cmdGetArticles.Flags().StringVarP(&siteStr, "sites", "i", "all", "Comma separated list of Gannett sites to fetch articles from")
	cmdGetArticles.Flags().StringVarP(&sectionStr, "sections", "e", "all", "Comma separated list of article sections to fetch from")
	cmdGetArticles.Flags().StringVarP(&mongoUri, "save", "s", "", "Saves articles to mongodb server specified in this option, e.g. mongodb://localhost:27017/mapi")
	cmdGetArticles.Flags().BoolVarP(&body, "body", "b", false, "Fetches the article body content")

	cmdRemoveArticles.Flags().BoolVarP(&noprompt, "noprompt", "n", false, "Skips the confirmation prompt and automatically removes articles")
	cmdRemoveArticles.Flags().StringVarP(&mongoUri, "save", "s", "", "Saves articles to mongodb server specified in this option, e.g. mongodb://localhost:27017/mapi")

	cmdBody.Flags().StringVarP(&articleUrl, "url", "u", url, "URL of Gannett article")
	cmdBody.Flags().BoolVarP(&includeTitle, "title", "t", false, "Place title of article on the first line of output")

	cmdSummary.Flags().StringVarP(&title, "title", "t", "", "Title for article summarizer, if not supplied then the summarizer assumes first line is title")

	cmdTopPages.Flags().StringVarP(&apiKey, "apikey", "k", "", "Chartbeat API Key. Required for fetching chartbeat data")
	cmdTopPages.Flags().StringVarP(&mongoUri, "save", "s", "", "Saves articles to mongodb server specified in this option, e.g. mongodb://localhost:27017/mapi")
	cmdTopPages.Flags().IntVarP(&loop, "loop", "l", -1, "Specify the internval in seconds to loop the fetching of the toppages api")
}

func AddCommands() {
	cmdArticles.AddCommand(cmdGetArticles)
	cmdArticles.AddCommand(cmdRemoveArticles)

	cmdChartbeat.AddCommand(cmdTopPages)

	NewsfetchCmd.AddCommand(cmdBody, cmdArticles, cmdSummary, cmdVersion, cmdChartbeat)
}

func getElapsedTime(sTime *time.Time) {
	endTime := time.Now()
	fmt.Println("Total time to run: ", endTime.Sub(*sTime))
}
