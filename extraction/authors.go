package extraction

import (
	"regexp"
	"strings"
)

func ParseAuthors(authors []string) []string {
	parsedAuthors := make([]string, 0, len(authors))

	for _, author := range authors {
		parsedAuthors = append(parsedAuthors, ParseAuthor(author)...)
	}

	return parsedAuthors
}

func ParseAuthor(author string) []string {
	splitAuthors := strings.Split(author, " and ")
	authors := make([]string, 0, len(splitAuthors))

	for _, testAuthor := range splitAuthors {
		// Parse out "by ..." and "and by..."
		regex := regexp.MustCompile(`(and )?by `)
		testAuthor = regex.ReplaceAllString(testAuthor, "")

		if testAuthor == "" {
			continue
		}
		authors = append(authors, testAuthor)
	}

	return authors
}
