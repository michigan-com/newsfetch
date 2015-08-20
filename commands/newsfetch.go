package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/michigan-com/newsfetch/lib"
	"github.com/op/go-logging"
	"github.com/spf13/cobra"
)

var (
	mongoUri     string
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
)

var logger = lib.GetLogger()

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

func Verbose(logLevel string) {
	level := logging.INFO
	var err error

	if logLevel != "" {
		level, err = logging.LogLevel(logLevel)
		if err != nil {
			logger.Error("Log level %s not found", logLevel)
		}
	}

	//env var trumps everything
	levelEnv := os.Getenv("LOGLEVEL")
	if levelEnv != "" {
		level, err = logging.LogLevel(levelEnv)
		if err != nil {
			logger.Error("Log level %s not found", logLevel)
		}
	}

	logging.SetLevel(level, "newsfetch")
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
}

func AddCommands() {
	cmdArticles.AddCommand(cmdGetArticles)
	cmdArticles.AddCommand(cmdRemoveArticles)

	NewsfetchCmd.AddCommand(cmdBody, cmdArticles, cmdSummary, cmdVersion)
}

func getElapsedTime(sTime *time.Time) {
	endTime := time.Now()
	fmt.Println("Total time to run: ", endTime.Sub(*sTime))
}
