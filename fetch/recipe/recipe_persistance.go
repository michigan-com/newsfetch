package fetch

import (
	"gopkg.in/mgo.v2/bson"

	"github.com/michigan-com/newsfetch/lib"
	"github.com/michigan-com/newsfetch/model/recipetypes"
)

func SaveRecipes(mongoUri string, recipes []*recipetypes.Recipe) error {
	session := lib.DBConnect(mongoUri)
	defer lib.DBClose(session)

	collection := session.DB("").C("Recipe")

	totalUpdates := 0
	totalInserts := 0
	for _, recipe := range recipes {
		// make sure the recipe doesn't have _id before upsert to avoid strange bugs stuff
		temp := *recipe
		temp.Id = ""
		info, err := collection.Upsert(bson.M{"article_id": recipe.ArticleId}, temp)
		if err != nil {
			panic(err)
		}

		if info.UpsertedId != nil {
			recipe.Id = info.UpsertedId.(bson.ObjectId)
			totalInserts++
		} else {
			totalUpdates++
		}
	}
	recipeDebugger.Println(totalUpdates, "recipes updated,", totalInserts, "recipes added")

	return nil
}

func LoadAllRecipes(mongoUri string) ([]*recipetypes.Recipe, error) {
	session := lib.DBConnect(mongoUri)
	defer lib.DBClose(session)

	collection := session.DB("").C("Recipe")
	var result []*recipetypes.Recipe
	err := collection.Find(nil).All(&result)
	return result, err
}
