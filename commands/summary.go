package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/urandom/text-summary/summarize"
)

var cmdSummary = &cobra.Command{
	Use:   "summary",
	Short: "Attempts to generate a summary based on an article body",
	Run: func(cmd *cobra.Command, args []string) {
		var summary summarize.Summarize

		if timeit {
			startTime = time.Now()
		}

		if title == "" {
			reader := bufio.NewReader(os.Stdin)
			contents, _ := ioutil.ReadAll(reader)
			lines := strings.Split(string(contents), "\n")

			summary = summarize.NewFromString(lines[0], strings.Join(lines[1:], "\n"))
		} else {
			summary = summarize.New(title, os.Stdin)
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
