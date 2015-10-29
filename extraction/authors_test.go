package extraction

import (
	"fmt"
	"testing"
)

type AuthorTestCase struct {
	Input    string
	Expected []string
}

type AuthorsTestCase struct {
	Input    []string
	Expected []string
}

func TestParseAuthor(t *testing.T) {
	testCases := []*AuthorTestCase{
		&AuthorTestCase{"greg gardner and brent snavely", []string{"greg gardner", "brent snavely"}},
		&AuthorTestCase{"test", []string{"test"}},
		&AuthorTestCase{"that test", []string{"that test"}},
		&AuthorTestCase{"by test", []string{"test"}},
		&AuthorTestCase{"by that test", []string{"that test"}},
		&AuthorTestCase{"and by test", []string{"test"}},
		&AuthorTestCase{"and by that test", []string{"that test"}},
		&AuthorTestCase{"by test and by test2", []string{"test", "test2"}},
		&AuthorTestCase{"test and test2", []string{"test", "test2"}},
		&AuthorTestCase{"test and by test2 and test 4 and by test 5", []string{"test", "test2", "test 4", "test 5"}},
		&AuthorTestCase{"really long author name", []string{"really long author name"}},
		&AuthorTestCase{"long name has the word \"and\" in it", []string{"long name has the word \"and\" in it"}},
		&AuthorTestCase{"long name here and a long name there and by a long name over there", []string{"long name here", "a long name there", "a long name over there"}},
	}

	for _, test := range testCases {
		result := ParseAuthor(test.Input)
		err := checkAuthorTestCase(test.Expected, result)

		if err != "" {
			t.Fatal(err)
		}

	}
}

func TestParseAuthors(t *testing.T) {
	testCases := []*AuthorsTestCase{
		&AuthorsTestCase{
			[]string{"greg gardner and brent snavely", "alisa priddle"},
			[]string{"greg gardner", "brent snavely", "alisa priddle"},
		},
	}

	for _, test := range testCases {
		result := ParseAuthors(test.Input)
		err := checkAuthorTestCase(test.Expected, result)

		if err != "" {
			t.Fatal(err)
		}
	}
}

func checkAuthorTestCase(expected []string, result []string) string {

	if len(result) != len(expected) {
		return fmt.Sprintf("Expected %v, got %v", expected, result)
	}

	for i := 0; i < len(result); i++ {
		e := expected[i]
		r := result[i]

		if e != r {
			return fmt.Sprintf("Author %d should have been '%s', was '%s'", i, e, r)
		}
	}

	return ""
}
