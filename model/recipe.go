package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type RecipeDuration struct {
	Text string `json:"text" bson:"text"`
}

type NutritionData struct {
	Text string `json:"text" bson:"text"`
}

type RecipeIngredient struct {
	Text string `json:"text" bson:"text"`
}

type RecipeInstruction struct {
	RawHtml string `json:"raw_html" bson:"raw_html,omitempty"`
	Text    string `json:"text" bson:"text"`
}

type RecipeImage struct {
	Url string `json:"url" bson:"url"`
}

type RecipePhoto struct {
	FullSizeImage *RecipeImage `json:"full" bson:"full,omitempty"`
	SmallImage    *RecipeImage `json:"small" bson:"small,omitempty"`
}

type Recipe struct {
	Id        bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	ArticleId int           `json:"article_id" bson:"article_id,omitempty"`
	Url       string        `json:"url" bson:"url,omitempty"`

	Title string `json:"title" bson:"title"`

	Photo *RecipePhoto `json:"photo" bson:"photo,omitempty"`

	ServingSize     string          `json:"serving_size" bson:"serving_size,omitempty"`
	PreparationTime *RecipeDuration `json:"prep_time" bson:"prep_time,omitempty"`
	TotalTime       *RecipeDuration `json:"total_time" bson:"total_time,omitempty"`

	Ingredients  []RecipeIngredient  `json:"ingredients" bson:"ingredients"`
	Instructions []RecipeInstruction `json:"instructions" bson:"instructions"`
	Nutrition    *NutritionData      `json:"nutrition" bson:"nutrition,omitempty"`
}

func NewRecipe() *Recipe {
	recipe := new(Recipe)
	recipe.Ingredients = make([]RecipeIngredient, 0)
	recipe.Instructions = make([]RecipeInstruction, 0)
	return recipe
}

func (r Recipe) String() string {
	b, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	out := new(bytes.Buffer)
	json.Indent(out, b, "", "  ")
	return out.String()
}

func (r Recipe) PlainLines() []string {
	var lines []string
	lines = append(lines, fmt.Sprintf("Title: %#v", r.Title))

	if r.ServingSize != "" {
		lines = append(lines, fmt.Sprintf("Serving size: %#v", r.ServingSize))
	}
	if r.TotalTime != nil {
		lines = append(lines, fmt.Sprintf("Total time: %#v", r.TotalTime.Text))
	}
	if r.PreparationTime != nil {
		lines = append(lines, fmt.Sprintf("Prep time: %#v", r.PreparationTime.Text))
	}
	for _, item := range r.Ingredients {
		lines = append(lines, fmt.Sprintf("I: %#v", item.Text))
	}
	for _, item := range r.Instructions {
		lines = append(lines, fmt.Sprintf("D: %#v", item.Text))
	}

	return lines
}
