package lib

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
	tag RecipeFragmentTag
}

func (f RecipeMarkerFragment) Tag() RecipeFragmentTag {
	return f.tag
}
func (f RecipeMarkerFragment) Mark(tag RecipeFragmentTag) {
	panic("Unsupported")
}

func (f RecipeMarkerFragment) AddToRecipe(recipe *Recipe) {
}

///////////////////////////////////////////////////////////////////////////////

type RecipeTimingFragment struct {
	ServingSize     string
	PreparationTime *RecipeDuration
	TotalTime       *RecipeDuration
}

func (f RecipeTimingFragment) Tag() RecipeFragmentTag {
	return TimingTag
}
func (f RecipeTimingFragment) Mark(tag RecipeFragmentTag) {
	panic("Unsupported")
}

func (f RecipeTimingFragment) AddToRecipe(r *Recipe) {
	r.ServingSize = f.ServingSize
	r.PreparationTime = f.PreparationTime
	r.TotalTime = f.TotalTime
}

///////////////////////////////////////////////////////////////////////////////

func (f RecipeIngredient) Tag() RecipeFragmentTag {
	return IngredientTag
}
func (f RecipeIngredient) Mark(tag RecipeFragmentTag) {
	panic("Unsupported")
}

func (f RecipeIngredient) AddToRecipe(r *Recipe) {
	r.Ingredients = append(r.Ingredients, f)
}

///////////////////////////////////////////////////////////////////////////////

type ParagraphFragment struct {
	RawHtml string
	Text    string
	tag     RecipeFragmentTag
}

func (f ParagraphFragment) Tag() RecipeFragmentTag {
	return f.tag
}
func (f *ParagraphFragment) Mark(tag RecipeFragmentTag) {
	f.tag = tag
}

func (f ParagraphFragment) AddToRecipe(r *Recipe) {
	switch f.tag {
	case TitleTag, PossibleTitleTag:
		r.Title = f.Text
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
