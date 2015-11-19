package commands

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

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

		matcher := rp.NewMatcher()

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
				s := fuz.CanonicalString(ingred.Text)

				subheadm := matcher.MatchIngredientSubhead(s)
				im := matcher.MatchIngredient(s)
				dm := matcher.MatchDirection(s)

				if subheadm.Confidence >= rp.Likely {
					continue
				}

				switch im.Confidence {
				case rp.Negative:
					incorrect.Add(s)
				case rp.Perfect:
					matchedIngredients++
					matched.Add(s + "\n")
					totalIngredients++
				case rp.Likely:
					totalIngredients++
				case rp.Possible:
					partial.Add(fmt.Sprintf("± %s\n  → %s", im.Rationale, s))
					problematicHere = append(problematicHere, s)
					totalIngredients++
				case rp.None:
					unmatched.Add(s)
					problematicHere = append(problematicHere, s)
					totalIngredients++
				}

				if dm.Confidence >= rp.Likely {
					// TODO: report conflict
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
				s := fuz.CanonicalString(dir.Text)
				dm := matcher.MatchDirection(s)
				if dm.Confidence < rp.Likely {
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
