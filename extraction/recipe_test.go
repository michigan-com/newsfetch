package extraction_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/michigan-com/newsfetch/extraction"
	a "github.com/michigan-com/newsfetch/fetch/article"
	m "github.com/michigan-com/newsfetch/model"
)

func eq(t *testing.T, comment string, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("%s = %#v, expected %#v", comment, actual, expected)
	}
}

func TestRecipeIntegration1(t *testing.T) {
	t.Skip("Fails currently")
	testRecipes(t, "http://freep.com/story/life/food/recipes/2015/10/03/broccoli-rabe-calzones/73222300/", `
		<p><strong>Broccoli Rabe and Mozzarella Calzones</strong></p>
		<p><strong>Makes: </strong>4 / <strong>Preparation time: </strong>15 minutes / <strong>Total time: </strong>30 minutes</p>
		<p>All-purpose flour or cornmeal for shaping the dough</p>
		<p>1 tablespoon olive oil</p>
		<p>1 large broccoli rabe, chopped, leaves and tough ends separated</p>
		<p>4 cloves garlic, finely chopped</p>
		<p>1 teaspoon red pepper flakes or to taste</p>
		<p>2 anchovy fillets, finely chopped, optional</p>
		<p>Salt and pepper taste</p>
		<p>1 recipe favorite pizza dough</p>
		<p>2 cups shredded mozzarella cheese</p>
		<p>Additional olive oil for brushing on top, optional</p>
		<p>Preheat the oven to 450 to 500 degrees. Sprinkle a small amount of flour over a baking sheet and set it aside.</p>
		<p>Place a large skillet over medium heat and add the olive oil. Once the oil is hot, add the tough ends of the broccoli rabe and cook for 2 minutes. Add the rest of the broccoli rabe, including the leafy part, along with the garlic, red pepper flakes and anchovies if using.</p>
		<p>Cook, stirring occasionally, until the stems are tender, about 5 minutes. Season with salt and pepper and set the filling aside.</p>
		<p>Sprinkle four on a clean work surface. Divide the pizza dough into 4 equal pieces and place one piece on the floured work surface. Roll out as you would for pizza until it’s quite thin. Pile a quarter of the broccoli rabe mixture and 1/2 cup of the mozzarella onto one side of the circle, leaving a lip around the edge.</p>
		<p>Gather up the half of the dough without the filling and fold it over the filling to create a half-moon shape. Pinch the edges together. Place the calzone on the baking sheet. Repeat with remaining dough and filling.</p>
		<p>Brush the tops with olive oil if desired.</p>
		<p>Bake calzones until they are golden brown on the outside, about 6 to 8 minutes. Remove from oven and serve.</p>
		<p><strong>Cook’s note: </strong>If you have Italian sausage in your fridge, crumble it into the pan with the broccoli rabe and cook until it’s cooked through.</p>
		<p><em>From “Good and Cheap: Eat Well on $4/Day” (Workman, $16.95). Tested by Susan Selasky for the Free Press Test Kitchen. Nutrition information not available.</em></p>
	`, `
		TODO
	`)
}

/* template

func TestRecipeIntegration1(t *testing.T) {
	testRecipes(t, "URL", `
	`, `
	`)
}
*/

func testRecipes(t *testing.T, url string, html string, expected string) {
	t.Logf("Testing URL: %v", url)

	var extract *m.ExtractedBody
	if strings.TrimSpace(html) == "" {
		_, html, e, err := a.ParseArticleAtURL(url, true)
		extract = e
		if err != nil {
			t.Fatalf("Failed to parse article: %v", err)
		}

		println("Here's the HTML to embed for", url)
		println(html)
	} else {
		extract = extraction.ExtractDataFromHTMLString(html, url, false)
	}

	var lines []string
	for _, embeddedURL := range extract.RecipeData.EmbeddedArticleUrls {
		lines = append(lines, fmt.Sprintf("Linked article: %s", embeddedURL))
	}

	for _, recipe := range extract.RecipeData.Recipes {
		lines = append(lines, "")
		lines = append(lines, recipe.PlainLines()...)
	}

	actual := strings.TrimSpace(strings.Join(lines, "\n"))
	expected = strings.TrimSpace(strings.Join(trimLines(strings.Split(expected, "\n")), "\n"))

	if actual != expected {
		t.Errorf("Got:\n%s\n\nExpected:\n%s", actual, expected)
	}
}

func trimLines(input []string) []string {
	result := make([]string, 0, len(input))
	for _, el := range input {
		result = append(result, strings.TrimSpace(el))
	}
	return result
}
