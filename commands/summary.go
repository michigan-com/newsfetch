package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/michigan-com/newsfetch/lib"
	"github.com/neurosnap/text-summary/summarize"
	"github.com/spf13/cobra"
)

var cmdSummary = &cobra.Command{
	Use:   "summary",
	Short: "Attempts to generate a summary based on an article body",
	Run: func(cmd *cobra.Command, args []string) {
		var summary summarize.Summarize
		tokenizer := lib.LoadTokenizer()

		if timeit {
			startTime = time.Now()
		}

		if title == "" {
			reader := bufio.NewReader(os.Stdin)
			contents, _ := ioutil.ReadAll(reader)
			lines := strings.Split(string(contents), "\n")

			summary = lib.NewPunktSummarizer(lines[0], strings.Join(lines[1:], "\n"), tokenizer)
		} else {
			text, _ := ioutil.ReadAll(os.Stdin)
			summary = lib.NewPunktSummarizer(title, string(text), tokenizer)
		}

		if output {
			for _, point := range summary.KeyPoints() {
				fmt.Println(point)
			}
		}

		if timeit {
			getElapsedTime(&startTime)
		}
	},
}
