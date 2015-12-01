package commands

import (
	"fmt"
	"os"
	"time"

	fr "github.com/michigan-com/newsfetch/fetch/recipe"
	// m "github.com/michigan-com/newsfetch/model"
	"github.com/michigan-com/newsfetch/model/recipetypes"
	"github.com/spf13/cobra"
)

var cmdVerifyRecipies = &cobra.Command{
	Use:   "verify",
	Short: "Load all recipes from the DB and verify them using the fuzzy classifier",
	Run: func(cmd *cobra.Command, args []string) {
		startTime = time.Now()
		stats := recipetypes.NewStats()

		recipes, err := fr.LoadAllRecipes(globalConfig.MongoUrl)
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(os.Stderr, "Loaded %d recipes.\n", len(recipes))

		for _, recipe := range recipes {
			_ = recipe
		}

		// if false {
		// 	for _, s := range incorrect.List {
		// 		fmt.Printf("~~~ %s\n", s)
		// 	}
		// }

		// if false {
		// 	for _, s := range unmatched.List {
		// 		fmt.Printf("! %s\n", s)
		// 	}
		// }

		// if false {
		// 	for _, s := range partial.List {
		// 		fmt.Printf("%s\n", s)
		// 	}
		// }

		// if true {
		// 	fmt.Printf("----------------------------------------\n")
		// 	for _, s := range problematic.List {
		// 		fmt.Printf("%s\n", s)
		// 	}
		// 	fmt.Printf("----------------------------------------\n")
		// }

		// if true {
		// 	fmt.Printf("----------------------------------------\n")
		// 	for _, s := range problematicDirections.List {
		// 		fmt.Printf("\n%s\n", s)
		// 	}
		// 	fmt.Printf("----------------------------------------\n")
		// }

		// sort.Strings(matched.List)

		// file, err := os.Create("/tmp/matched.txt")
		// if err != nil {
		// 	panic(err)
		// }
		// _, err = file.WriteString(strings.Join(matched.List, ""))
		// if err != nil {
		// 	panic(err)
		// }
		// err = file.Close()
		// if err != nil {
		// 	panic(err)
		// }

		format := "%-20s %d\n"
		fmt.Printf(format, "Total recipes:", stats.RecipeTotal)
		fmt.Printf(format, "Fully matched:", stats.RecipeFullyMatched)
		fmt.Printf(format, "Less than 3 unmatched:", stats.Recipe3ToGo)
		fmt.Printf(format, "Others:", stats.RecipeTotal-stats.RecipeFullyMatched-stats.Recipe3ToGo)

		getElapsedTime(&startTime)
	},
}
