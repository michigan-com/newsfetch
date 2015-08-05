package cli

import (
	"flag"
	"fmt"
	"log"

	"github.com/michigan-com/newsfetch/lib"
)

func main() {
	flag.Parse()

	url := flag.Arg(0)
	if url == "" {
		url = "http://detroitnews.com/story/news/local/detroit-city/2015/08/04/female-body-found-possible-hit-run-detroit/31094589/"
	}

	body, err := lib.ExtractBodyFromURL(url)
	if err != nil {
		panic(err)
	}

	log.Printf("Success.")
	fmt.Printf("Body:\n-----\n%s\n-----\n", body)
}
