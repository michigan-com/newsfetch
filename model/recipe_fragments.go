package model

type RecipeFragmentTag int

const (
	NoTag RecipeFragmentTag = iota
	TitleTag
	PossibleTitleTag

	ParagraphTag
	ShortParagraphTag // paragraph that might as well be an ingredient
	InstructionTag
	IngredientsSubtitleTag
	DirectionsSubtitleTag

	PossibleIngredientSubdivisionTag // normal paragraph unless proven otherwise

	TimingTag
	ServingSizeAltTimingTag // “This recipe serves ...”, found in the body of the article
	IngredientTag
	IngredientSubdivisionTag
	PossibleIngredientTag
	NutritionDataTag

	SignatureTag
	CopyrightTag
	EndMarkerTag
)

///////////////////////////////////////////////////////////////////////////////

type RecipeFragment interface {
	Tag() RecipeFragmentTag
	Mark(tag RecipeFragmentTag)
	AddToRecipe(recipe *Recipe)
}

///////////////////////////////////////////////////////////////////////////////

type RecipeMarkerFragment struct {
	TagF RecipeFragmentTag
}

func (f RecipeMarkerFragment) Tag() RecipeFragmentTag {
	return f.TagF
}
func (f RecipeMarkerFragment) Mark(tag RecipeFragmentTag) {
	panic("Unsupported")
}

func (f RecipeMarkerFragment) AddToRecipe(recipe *Recipe) {
}

///////////////////////////////////////////////////////////////////////////////

type RecipeTimingFragment struct {
	TagF            RecipeFragmentTag
	ServingSize     string
	PreparationTime *RecipeDuration
	TotalTime       *RecipeDuration
}

func (f RecipeTimingFragment) Tag() RecipeFragmentTag {
	return f.TagF
}
func (f RecipeTimingFragment) Mark(tag RecipeFragmentTag) {
	panic("Unsupported")
}

func (f RecipeTimingFragment) AddToRecipe(r *Recipe) {
	if f.ServingSize != "" {
		r.ServingSize = f.ServingSize
	}
	if f.PreparationTime != nil {
		r.PreparationTime = f.PreparationTime
	}
	if f.TotalTime != nil {
		r.TotalTime = f.TotalTime
	}
}

///////////////////////////////////////////////////////////////////////////////

type ParagraphFragment struct {
	RawHtml string
	Text    string
	TagF    RecipeFragmentTag
}

func (f ParagraphFragment) Tag() RecipeFragmentTag {
	return f.TagF
}
func (f *ParagraphFragment) Mark(tag RecipeFragmentTag) {
	f.TagF = tag
}

func (f ParagraphFragment) AddToRecipe(r *Recipe) {
	switch f.TagF {
	case TitleTag, PossibleTitleTag:
		if r.Title == "" {
			r.Title = f.Text
		}
	case InstructionTag, ParagraphTag, ShortParagraphTag, PossibleIngredientSubdivisionTag:
		instruction := RecipeInstruction{RawHtml: f.RawHtml, Text: f.Text}
		r.Instructions = append(r.Instructions, instruction)
	case IngredientTag, PossibleIngredientTag, IngredientSubdivisionTag:
		ingredient := RecipeIngredient{Text: f.Text}
		r.Ingredients = append(r.Ingredients, ingredient)
	}
}

///////////////////////////////////////////////////////////////////////////////

func (f NutritionData) Tag() RecipeFragmentTag {
	return NutritionDataTag
}
func (f NutritionData) Mark(tag RecipeFragmentTag) {
	panic("Unsupported")
}

func (f NutritionData) AddToRecipe(r *Recipe) {
	r.Nutrition = &f
}
