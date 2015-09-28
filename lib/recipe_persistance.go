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
		existing := Recipe{}

		err := collection.Find(bson.M{"article_id": recipe.ArticleId}).Select(bson.M{"_id": 1}).
			One(&existing)
		if err == nil {
			collection.Update(bson.M{"_id": existing.Id}, recipe)
			totalUpdates++
		} else {
			collection.Insert(recipe)
			totalInserts++
		}
	}
	Debugger.Println(totalUpdates, "recipes updated,", totalInserts, "recipes added")

	return nil
}
