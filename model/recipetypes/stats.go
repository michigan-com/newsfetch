package recipetypes

type Stats struct {
	RecipeTotal        int
	RecipeFullyMatched int
	Recipe3ToGo        int

	NumMatchedIngredients int

	PerfectIngredients []string
	MatchedIngredients []string

	PerfectDirections []string
	MatchedDirections []string

	MisreportedIngredients []string
	MisreportedDirections  []string

	Conflicting []string

	PartialIngredients []string

	AssignedIngredients []string
	AssignedDirections  []string

	IngredientRanking map[string]int
}

func NewStats() *Stats {
	stats := new(Stats)

	stats.IngredientRanking = make(map[string]int, 100)

	return stats
}

func (stats *Stats) AddIngredient(text string, rank int) {
	if orank, ok := stats.IngredientRanking[text]; ok {
		if orank < rank {
			rank = orank
		} else if orank == rank {
			return // no changes
		}
	}
	stats.IngredientRanking[text] = rank
}
