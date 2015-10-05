package lib

import (
	"gopkg.in/mgo.v2/bson"
)

func SaveRecipes(mongoUri string, recipes []*Recipe) error {
	session := DBConnect(mongoUri)
	defer DBClose(session)

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
	Debugger.Println(totalUpdates, "recipes updated,", totalInserts, "recipes added")

	return nil
}
