package recipematcher

import (
	"strings"
	"testing"

	"github.com/michigan-com/newsfetch/util/fuzzy"
)

func TestIngredientClassifierSeparateForDebugging(t *testing.T) {
	otag(t, "@ingredient", `1 tablespoon lemon juice`)
}

func TestIngredientClassifierNewBatch(t *testing.T) {
	otag(t, "@ingredient", `
# 10 large basil leaves, divided
# 1/4 cup plus 4 teaspoons granulated sugar, divided
# 1/2 cup white rum
# 5 to 6 large russet potatoes, scrubbed
# 12 ounces bacon
# 2 cups sharp cheddar
# 1 can (14 ounces) artichoke hearts, drained
# 6 extra-large eggs
# 5 tablespoons freshly grated Parmesan cheese, divided
# 1 fully baked 9-inch pie shell (crumb, shortbread or flaky)
# 1 cup (2 sticks) butter, softened
# 4 ounces unsweetened chocolate, melted and cooled
# 2 ripe avocados, peeled, pitted, diced
# 2 ripe plum tomatoes, washed, cored, diced
# 6 hard-cooked eggs, peeled, quartered
# 2 cups fresh basil leaves
# 1/2 cup olive oil, divided
# 2 large, about 1 1/4 pounds, (or 4 small) skin-on, boneless chicken breasts, trimmed of excess fat
# 4 bone-in pork chops about 1-inch thick (about 2 pounds total)
# 2 strips orange zest, plus 1/2 cup orange juice and 2 oranges such as Cara Cara, peeled and with pith removed and flesh cut into segments
# 1 1/2 pounds asparagus (about 1 1/2 bunches), rinsed
# 1 small Yukon Gold potato (about 6 ounces), peeled and thinly sliced
# 1 pound fresh raw shrimp (size of choice)
# 1 pound fresh large bay scallops or sea scallops
# 1 cup shredded fresh cilantro, plus more for garnish
# 2 cups fresh raspberries and blueberries
# 4 cups Oat and Nut Granola with Dried Fruit (see recipe) or good-quality store-bought granola
# 4 sprigs fresh mint, for garnish
# 12 to 16 lemon or orange slices at least 1/4- to 1/3-inch thick
# 1/2 cup fat-free, reduced-sodium vegetable broth
# 1 cup fresh corn kernels cut from cob (about 2 ears)
# 1 tablespoon trans fat-free margarine

## too weird
# 2 strips orange zest, plus 1/2 cup orange juice and 2 oranges such as Cara Cara, peeled and with pith removed and flesh cut into segments
# Nonfat chicken broth instead of milk.
# Peeled, whole garlic cloves; mash them in with the potatoes.

	`)
}

func TestIngredientClassifier(t *testing.T) {
	otag(t, "@ingredient", `
1 garlic clove, peeled and crushed
1 stick salted butter, softened
1 stick softened butter, see note
8 ounces good-quality milk chocolate
5 ounces creamy or crunchy peanut butter, at room temperature
3 cups crispy puffed rice cereal
11/2 cups all-purpose flour, plus more for dusting
4 large cage-free egg yolks
1 to 11/2 tablespoons water, plus more as needed
4 bone-in pork chops about 1-inch thick (about 2 pounds total)
1/4 cup fresh parsley leaves, chopped
2 small heads cauliflower
1 clove garlic
2 tablespoons crumbled feta cheese
3 medium-large zucchini (about 1 1/2 pounds)
2 cups low-sodium chicken or vegetable broth, warmed to hot
1 tablespoon lemon juice
4 center-cut salmon fillets, about 5 ounces each
12 to 16 lemon or orange slices at least 1/4- to 1/3-inch thick
8 bone-in, skin-on chicken thighs (about 2 1/2 pounds)
2 tablespoons roughly chopped fresh chives or parsley (optional)
1 pound whole strawberries, washed, hulled
1 recipe whole wheat yeast dough
1 cup halved red or yellow cherry tomatoes, gently squeezed to remove seeds
1 tablespoon chopped fresh rosemary leaves
4 boneless strip steaks (about 9 to 11 ounces each) 1 inch thick, trimmed of excess fat
8 to 10 cups mixed salad greens
1/2 cup blue cheese crumbles
3 large egg whites
3 cups nonfat plain yogurt
2 teaspoons pure vanilla extract
1/2 cup orange juice
2 tablespoons maple syrup, optional
2 tablespoons butter, cut into pieces, optional

Roasted garlic.
Peeled and quartered rutabaga, celery root, turnip or parsnip.

Freshly grated Romano or Parmesan cheese, about 1/2 cup
Leaves from 1 bunch Thai basil or sweet basil
Lard or vegetable oil as needed to soften tortillas
Small bunch of parsley, minced
Herb sprigs (rosemary, thyme, parsley), optional
Leaves from 1 bunch of fresh cilantro
Candy corn
Tortilla chips or strips, optional
Caramelized onions, optional
Sour cream
Curls or julienne of lemon zest
Canola oil for the griddle
Fried rice or brown rice for serving, optional
Kernels from 6 ears of corn (about 3 cups)
Guacamole, sour cream, shredded lettuce, diced tomato
Floured baking spray
Nonstick vegetable oil spray
Zest and juice from 1 lemon
Olive oil-flavored nonstick cooking spray
Sushi grade Ahi Tuna
Blend of cracked peppercorns
Grilled baby romaine hearts
Chopped fresh Italian parsley, basil, or chives, for garnish
Oregano or parsley sprigs for garnish, optional
Fresh herbs such as basil or dill, optional
Flavored butter (see cook's note)
Vegetable oil cooking spray or floured baking spray
Corn kernels cut from 2 ears fresh corn, grilled
Orange rind (optional)
Finely grated zest of 1 lime
Fresh shredded herbs such as basil or tarragon, optional

All-purpose seasoning to taste
Sage leaves for garnish
Water to cover about 8 cups
Oil for the grill grates
Guacamole
Queso Chihauhua, shredded or queso fresco crumbled, optional
Vegetable oil for frying
Fresh rosemary sprigs
Creme fraiche for garnish, optional
A few sprigs fresh sage, very thinly sliced
Oil for deep-frying
Reserved juice from the grapefruit (about 1/2 cup)
Risotto or mashed potatoes for serving
Fresh chopped herbs, optional
Condiments as desired
Pita bread for serving
Watermelon, Avocado and Couscous salad, optional (see cook's note)
Maraschino cherries for garnish
About 14 cups water
Pinch of saffron
Pastry for a double pie crust, homemade or store-bought
Parchment paper
One (12- to 14-pound) fresh turkey (or frozen, thawed)
Leaves from 2 sprigs fresh thyme
Sauteed greens such as kale or mustard greens
Poppy seeds or sesame seeds, optional
Oil for frying
Chopped fresh tarragon
Bamboo or metal skewers
Pea tendrils, optional, for garnish
Thinly sliced baby golden beets
Baby pea shoots
Vegetable oil for cooking
Fresh berries, optional
Sprigs of thyme for garnish
Two large baby back pork ribs or two racks, 4 pounds each, spareribs, membranes removed
Favorite barbecue sauce
Oil for the grill
Tomato and cucumber slices if desired
Favorite mild vinaigrette, about 1/2 cup

Zest of one small lemon
Salt and pepper to taste
Nonstick floured baking spray
Pinch salt
Pinch nutmeg
Juice of half a lemon
Sugar for sprinkling
Nonstick cooking spray
Flour for work surface
Confectioners’ sugar for dusting, optional
Favorite homemade salted caramel sauce or store-bought
Salt to taste
Freshly cracked black pepper
Vegetable oil cooking spray
Fresh flat-leaf Italian parsley
Snipped fresh parsley (optional)
Squeeze of lemon juice
Kosher salt and freshly ground black pepper to taste
Fried onion rings, optional
Chopped fresh flat-leaf parsley, optional
Salt and pepper to taste, or favorite all-purpose seasoning
Fortified whipped frosting (see cook’s note)
Freshly ground black pepper to taste
Kosher salt to taste
Balsamic syrup, optional (see cook’s note)
Fresh chopped cilantro, optional
Butter, for greasing pan
Confectioner's sugar for dusting (optional)
A pinch of kosher salt
About 1/2 cup club soda
Lime wedges
Fresh seasonal berries
Corn on the cob
Dash cumin
About 1/2 cup heavy cream
Freshly ground black pepper
Flaky Maldon salt or sea salt and freshly ground black pepper to taste
Shredded lettuce
All-purpose flour or cornmeal for shaping the dough
Salt and pepper taste
Additional olive oil for brushing on top, optional
Pinch cayenne pepper
Pinch red pepper flakes (optional)
Olive oil
Sea salt or kosher salt to taste
Chopped green onion
Kosher salt and freshly ground pepper to taste
Canola oil
Salt and black pepper to taste
Paprika
Bourbon Barbecue Sauce (See note)
Salt and freshly ground black pepper to taste
Chopped parsley
Salt and ground pepper
Grated zest of 1 lemon
Salt and freshly ground black pepper, if desired
Pinch of red pepper flakes
Kosher salt
Freshly ground white pepper
Salt and pepper to taste or favorite all-purpose seasoning
Pinch of salt
Salsa, optional
Pinch of black pepper
Croutons for garnish
Decorating gel
Gummy worms
Salt and freshly ground pepper, to taste
Additional Parmesan cheese, if desired
Coarse salt and ground pepper
Lemon slices, for serving
Salad greens for serving, optional
Mexican blend shredded cheese, optional
Zest of 2 oranges
Freshly chopped chives or parsley
Whipped cream
Semisweet chocolate
Cooked basmati rice, optional
Salt and pepper, to taste
Favorite all-purpose seasoning mix or salt and pepper to taste
Sea salt and freshly ground black pepper to taste
Vegetable oil spray
Salt and freshly ground black pepper
Corn chips
Shredded cheese
Sliced green onions
Favorite salsa
Green onion whites from above
Confectioner's sugar, optional
Juice of 1 large lime
Pickled red onion slices (see note)
Chopped cilantro and onion mix
Kosher or coarse sea salt to taste
Olive oil cooking spray
Chopped parsley for garnish
Kosher salt and freshly ground black pepper
Zest of 3 medium limes
Zest of 1 medium lime
Fine sea salt
Ground white pepper to taste
Juice of one lime
Sliced green onion for garnish
Lime wedges for garnish, optional
Steamed rice or noodles for serving, optional
Sea salt and ground black pepper
Favorite lemon pepper or other citrus type of seasoning
Salt and white pepper to taste
Salt
Black pepper to taste
Freshly grated nutmeg
Parmesan cheese, freshly grated, for serving (about 1/2 cup)
Salt and ground black pepper to taste
Sea salt and black pepper to taste
Fresh basil for serving
Juice from half lemon
Sea salt and freshly ground black pepper, to taste
Nonstick cooking spray or canola oil
Crushed red pepper flakes to taste
Chopped fresh parsley, for serving
Salt and ground black pepper
Salt to taste, optional
Juice of 1 lime, plus lime wedges for serving
Kosher salt and ground pepper
Kosher salt and ground black pepper
Fresh lemon juice (optional)
White sanding sugar, for garnish
Fresh chopped herbs such as parsley, dill, sage and rosemary.
Extra-virgin olive oil
Cornmeal, for sprinkling
Coarse sea salt
Ground black pepper
Croutons
Juice of half a lime, plus more for serving
Pinch of kosher salt and freshly ground black pepper
Pinch of salt and freshly ground black pepper
Pinch of ground cayenne pepper
Sunchoke chips
Thinly sliced radish
Butter-flavored nonstick cooking spray
Pinch of kosher salt
Confectioners’ sugar, for dusting, optional
Fresh cilantro leaves, for garnish
Pinch of cream of tartar
Grated zest of 1/2 lemon
Pinch kosher salt
Pinch freshly ground white pepper
Pinch sugar
Generous pinch of white pepper, plus more to taste
Juice of 1/2 lemon
Kosher salt and freshly ground pepper
A few sprigs of flat-leaf parsley (or another herb) for garnish
Freshly grated Parmesan cheese
Ground black pepper to taste
Generous pinch of crushed red pepper flakes
Medium coarse salt and freshly ground black pepper to taste
A few drops favorite hot sauce
Juice of 2 lemons
Juice of 2 limes
Juice and zest from one lemon
Cinnamon sticks (optional)
Mixed greens with a light vinaigrette for serving
All-purpose flour, for dusting
Freshly grated Parmesan, for serving, optional
Sea salt and freshly ground black pepper
Juice of 1 lemon
Salt and black pepper, optional
Pinch of sugar
Several dashes of hot sauce
Sea salt and cracked black pepper
		`)
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
	//otag(t, c, "@ingredient", "")
}

func otag(t *testing.T, tag string, inputs string) {
	for _, input := range strings.Split(inputs, "\n") {
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}
		if strings.HasPrefix(input, "#") {
			continue
		}

		r := classifier.Process(input)
		actual, _ := r.GetTagMatchString(tag, fuzzy.Raw)

		inputAdjusted := fuzzy.CanonicalString(input)

		if actual != inputAdjusted {
			if actual != "" {
				t.Errorf("In %#v only matched %#v", input, actual)
			} else {
				t.Errorf("In %#v no match", input)
			}
			t.Logf("Result =\n%v", r.Description())
		}
	}
}
