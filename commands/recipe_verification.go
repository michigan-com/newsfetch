package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"
	"unicode"

	fuz "github.com/michigan-com/newsfetch/extraction/fuzzy_classifier"
	rp "github.com/michigan-com/newsfetch/extraction/recipe_parsing"
	fr "github.com/michigan-com/newsfetch/fetch/recipe"
	"github.com/spf13/cobra"
)

type VerificationStats struct {
	RecipeTotal        int
	RecipeFullyMatched int
	Recipe3ToGo        int
}

var cmdVerifyRecipies = &cobra.Command{
	Use:   "verify",
	Short: "Load all recipes from the DB and verify them using the fuzzy classifier",
	Run: func(cmd *cobra.Command, args []string) {
		startTime = time.Now()
		stats := VerificationStats{}

		recipes, err := fr.LoadAllRecipes(globalConfig.MongoUrl)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(os.Stderr, "Loaded %d recipes.\n", len(recipes))

		classifier := rp.NewIngredientClassifier()
		dc := rp.NewDirectionClassifier()

		incorrect := new(UniqueStringList)
		matched := new(UniqueStringList)
		partial := new(UniqueStringList)
		unmatched := new(UniqueStringList)
		problematic := new(UniqueStringList)
		problematicDirections := new(UniqueStringList)

		for _, recipe := range recipes {
			totalIngredients := 0
			matchedIngredients := 0

			var problematicHere []string
			for _, ingred := range recipe.Ingredients {
				s := strings.TrimSpace(ingred.Text)

				s = strings.Replace(s, "½", "1/2", -1)
				s = strings.Replace(s, "¼", "1/4", -1)
				s = strings.Replace(s, "¾", "3/4", -1)

				runes := []rune(s)

				adjusted := s
				adjusted = strings.Replace(adjusted, "(optional)", "", 1)

				if IsEntirelyUppercase(adjusted) {
					continue
				}

				_ = runes
				if unicode.IsDigit(runes[0]) {
					continue
				}

				r := classifier.Process(s)

				if _, incor := r.GetTagMatchString("@not_ingredient", fuz.Raw); incor {
					incorrect.Add(s)
					continue
				}

				totalIngredients++

				ms, _ := r.GetTagMatchString("@ingredient", fuz.Raw)
				if ms == s {
					matchedIngredients++
					matched.Add(s + "\n")
				} else {
					if ms != "" {
						partial.Add(fmt.Sprintf("± %s\n  → %s", ms, s))
					} else {
						unmatched.Add(s)
					}
					problematicHere = append(problematicHere, s)
				}

			}

			stats.RecipeTotal++

			unmatchedIngredients := totalIngredients - matchedIngredients
			if unmatchedIngredients == 0 {
				stats.RecipeFullyMatched++
			} else if unmatchedIngredients <= 3 {
				stats.Recipe3ToGo++
			}

			if unmatchedIngredients <= 3 {
				problematic.AddList(problematicHere)
			}

			for _, dir := range recipe.Instructions {
				s := strings.TrimSpace(dir.Text)
				r := dc.Process(s)
				_, ok := r.GetTagMatchString("@direction", fuz.Raw)
				if !ok {
					problematicDirections.Add(s)
				}
			}
		}

		if false {
			for _, s := range incorrect.List {
				fmt.Printf("~~~ %s\n", s)
			}
		}

		if false {
			for _, s := range unmatched.List {
				fmt.Printf("! %s\n", s)
			}
		}

		if false {
			for _, s := range partial.List {
				fmt.Printf("%s\n", s)
			}
		}

		if true {
			fmt.Printf("----------------------------------------\n")
			for _, s := range problematic.List {
				fmt.Printf("%s\n", s)
			}
			fmt.Printf("----------------------------------------\n")
		}

		if true {
			fmt.Printf("----------------------------------------\n")
			for _, s := range problematicDirections.List {
				fmt.Printf("\n%s\n", s)
			}
			fmt.Printf("----------------------------------------\n")
		}

		sort.Strings(matched.List)

		file, err := os.Create("/tmp/matched.txt")
		if err != nil {
			panic(err)
		}
		_, err = file.WriteString(strings.Join(matched.List, ""))
		if err != nil {
			panic(err)
		}
		err = file.Close()
		if err != nil {
			panic(err)
		}

		format := "%-20s %d\n"
		fmt.Printf(format, "Total recipes:", stats.RecipeTotal)
		fmt.Printf(format, "Fully matched:", stats.RecipeFullyMatched)
		fmt.Printf(format, "Less than 3 unmatched:", stats.Recipe3ToGo)
		fmt.Printf(format, "Others:", stats.RecipeTotal-stats.RecipeFullyMatched-stats.Recipe3ToGo)

		getElapsedTime(&startTime)
	},
}

func IsEntirelyUppercase(s string) bool {
	return s == strings.ToUpper(s) && strings.IndexFunc(s, unicode.IsLetter) != -1
}

type UniqueStringList struct {
	List []string
	Set  map[string]bool
}

func (u *UniqueStringList) Add(s string) {
	if !u.Set[s] {
		if u.Set == nil {
			u.Set = make(map[string]bool)
		}
		u.List = append(u.List, s)
		u.Set[s] = true
	}
}

func (u *UniqueStringList) AddList(items []string) {
	for _, item := range items {
		u.Add(item)
	}
}
