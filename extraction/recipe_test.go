package extraction_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/michigan-com/newsfetch/extraction"
	a "github.com/michigan-com/newsfetch/fetch/article"
	m "github.com/michigan-com/newsfetch/model"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func eq(t *testing.T, comment string, actual, expected interface{}) {
	if actual != expected {
		t.Errorf("%s = %#v, expected %#v", comment, actual, expected)
	}
}

var skipFailingRecipes = true

func TestRecipeIntegration1(t *testing.T) {
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
		Title: "Broccoli Rabe and Mozzarella Calzones"
		Serving size: "4"
		Total time: "30 minutes"
		Prep time: "15 minutes"
		I: "All-purpose flour or cornmeal for shaping the dough"
		I: "1 tablespoon olive oil"
		I: "1 large broccoli rabe, chopped, leaves and tough ends separated"
		I: "4 cloves garlic, finely chopped"
		I: "1 teaspoon red pepper flakes or to taste"
		I: "2 anchovy fillets, finely chopped, optional"
		I: "Salt and pepper taste"
		I: "1 recipe favorite pizza dough"
		I: "2 cups shredded mozzarella cheese"
		I: "Additional olive oil for brushing on top, optional"
		D: "Preheat the oven to 450 to 500 degrees. Sprinkle a small amount of flour over a baking sheet and set it aside."
		D: "Place a large skillet over medium heat and add the olive oil. Once the oil is hot, add the tough ends of the broccoli rabe and cook for 2 minutes. Add the rest of the broccoli rabe, including the leafy part, along with the garlic, red pepper flakes and anchovies if using."
		D: "Cook, stirring occasionally, until the stems are tender, about 5 minutes. Season with salt and pepper and set the filling aside."
		D: "Sprinkle four on a clean work surface. Divide the pizza dough into 4 equal pieces and place one piece on the floured work surface. Roll out as you would for pizza until it’s quite thin. Pile a quarter of the broccoli rabe mixture and 1/2 cup of the mozzarella onto one side of the circle, leaving a lip around the edge."
		D: "Gather up the half of the dough without the filling and fold it over the filling to create a half-moon shape. Pinch the edges together. Place the calzone on the baking sheet. Repeat with remaining dough and filling."
		D: "Brush the tops with olive oil if desired."
		D: "Bake calzones until they are golden brown on the outside, about 6 to 8 minutes. Remove from oven and serve."
		D: "Cook’s note: If you have Italian sausage in your fridge, crumble it into the pan with the broccoli rabe and cook until it’s cooked through."
	`)
}

func TestRecipeIntegration2(t *testing.T) {
	testRecipes(t, "http://www.freep.com/story/life/food/recipes/2015/10/03/salty-broccoli-toast-recipes/73310358/", `
		<p><strong>Salty Broccoli on Toast</strong></p>
		<p><strong>Makes:</strong> 4 toasts / <strong>Preparation time: </strong>10 minutes / <strong>Total time: </strong>10 minutes</p>
		<p>Leanne Brown likes garlic and spice, but you can cut back the amounts of both in this recipe to taste.</p>
		<p>1 teaspoon olive oil</p>
		<p>3 cloves garlic, peeled, chopped</p>
		<p>1 teaspoon red pepper flakes or to taste</p>
		<p>1 anchovy fillet, chopped</p>
		<p>1 small head of broccoli with stem, chopped</p>
		<p>1/4 cup water</p>
		<p>4 slices hot toast</p>
		<p>Freshly grated Romano or Parmesan cheese, about 1/2 cup</p>
		<p>Salt and pepper to taste</p>
		<p>In a large skillet, warm up the oil over medium heat. Add the garlic and red pepper flakes and cook for 2 minutes, until they smell great but are not yet brown. Add the anchovy and cook for another minute. Add the broccoli and about 1/4 cup of water. Cover the pan, steam for 3 minutes, then toss and cook for 2 minutes, until the broccoli is tender and the water is gone. Spoon onto hot toast; top with a sprinkling of cheese, salt and pepper.</p>
		<p><strong>Variations: </strong><strong>Korean-style Spinach:</strong> In a large skillet, heat 1 teaspoon olive oil over medium heat. Add 4 cloves chopped garlic and cook 2 minutes or until fragrant. Add 1 bunch, washed spinach (remove tough ends) and 1 teaspoon soy sauce; cook 2 minutes or until spinach has wilted and shrunk. Turn off the heat and add 1/2 teaspoon toasted sesame oil and salt to taste. Mix together and taste; adjusting seasoning as needed. Remove the spinach from the pan and squeeze out any excess moisture. Serve over hot slices of toast. Sprinkle 1 teaspoon sesame seeds top.</p>
		<p><strong>Apple Cheddar Toast: </strong>Thinly slice 2 ounces Cheddar cheese and 1 small apple. Brown prefers to layer the apples on the hot toast and then slide small slices of Cheddar in between them, but do whatever makes sense to you. Season the toppings with salt and pepper to taste if desired.</p>
		<p>From “Good and Cheap: Eat Well on $4/Day” (Workman, $16.95). Tested by Susan Selasky for the Free Press Test Kitchen. Nutrition information not available.</p>
	`, `
		Title: "Salty Broccoli on Toast"
		Serving size: "4 toasts"
		Total time: "10 minutes"
		Prep time: "10 minutes"
		I: "1 teaspoon olive oil"
		I: "3 cloves garlic, peeled, chopped"
		I: "1 teaspoon red pepper flakes or to taste"
		I: "1 anchovy fillet, chopped"
		I: "1 small head of broccoli with stem, chopped"
		I: "1/4 cup water"
		I: "4 slices hot toast"
		I: "Freshly grated Romano or Parmesan cheese, about 1/2 cup"
		I: "Salt and pepper to taste"
		D: "Leanne Brown likes garlic and spice, but you can cut back the amounts of both in this recipe to taste."
		D: "In a large skillet, warm up the oil over medium heat. Add the garlic and red pepper flakes and cook for 2 minutes, until they smell great but are not yet brown. Add the anchovy and cook for another minute. Add the broccoli and about 1/4 cup of water. Cover the pan, steam for 3 minutes, then toss and cook for 2 minutes, until the broccoli is tender and the water is gone. Spoon onto hot toast; top with a sprinkling of cheese, salt and pepper."
		D: "Variations: Korean-style Spinach: In a large skillet, heat 1 teaspoon olive oil over medium heat. Add 4 cloves chopped garlic and cook 2 minutes or until fragrant. Add 1 bunch, washed spinach (remove tough ends) and 1 teaspoon soy sauce; cook 2 minutes or until spinach has wilted and shrunk. Turn off the heat and add 1/2 teaspoon toasted sesame oil and salt to taste. Mix together and taste; adjusting seasoning as needed. Remove the spinach from the pan and squeeze out any excess moisture. Serve over hot slices of toast. Sprinkle 1 teaspoon sesame seeds top."
		D: "Apple Cheddar Toast: Thinly slice 2 ounces Cheddar cheese and 1 small apple. Brown prefers to layer the apples on the hot toast and then slide small slices of Cheddar in between them, but do whatever makes sense to you. Season the toppings with salt and pepper to taste if desired."
	`)
}

func TestRecipeIntegrationLinks1(t *testing.T) {
	if skipFailingRecipes {
		t.SkipNow()
	}
	testRecipes(t, "http://www.freep.com/story/life/food/recipes/2015/09/23/apple-season-apple-pie-recipes/72691626/", `
		<p>Apples and pie go hand in hand.</p>
		<p>But a slice of pie isn’t the most convenient treat for those on the go. That’s where today’s Test Kitchen treats come in. Our recipes have all the ingredients of apple pie with a few unique twists — mini apple pies made in a muffin tin, apple pie bars drizzled with salted caramel, and apple pie cookies complete with a lattice “crust.”</p>
		<p><strong>Related:</strong><a href="http://www.freep.com/story/life/food/recipes/2015/09/24/free-press-apple-guide/72685434/">U-pick apple orchards and cider mills</a></p>
		<p>Now is the time to get baking, with Michigan apple season in full swing. At area apple orchards and cider mills you’ll find an abundance of apples. This year’s crop is expected to be larger than average: 24 million bushels (22.83 million bushels is the Michigan average).</p>
		<p><em>Contact Susan Selasky: 313-222-6872 or <a href="mailto:sselasky@freepress.com">sselasky@freepress.com</a>. Follow her on Twitter @SusanMariecooks.</em></p>`, `
		Linked article: http://www.freep.com/story/life/food/recipes/2015/09/23/apple-roses-recipe/72691592/
		Linked article: http://www.freep.com/story/life/food/recipes/2015/09/23/mini-apple-pies/72694116/
		Linked article: http://www.freep.com/story/life/food/recipes/2015/09/23/caramel-apple-pie-cookies/72691930/
		Linked article: http://www.freep.com/story/life/food/recipes/2015/09/23/salted-caramel-apple-pie-bars/72692424/
	`)
}

func TestRecipeIntegrationLinks2(t *testing.T) {
	if skipFailingRecipes {
		t.SkipNow()
	}
	testRecipes(t, "http://www.freep.com/story/life/food/recipes/2015/09/11/healthier-comfort-foods-fall/72016524/", `
		<p>Just around the corner is the official start of fall, the time of year comfort food cravings really settle in. There are creamy and hearty soups and chowders, cheesy lasagnas and casseroles drenched in cream-of-something soup.</p>
	`, `
		Linked article: http://www.freep.com/story/life/food/recipes/2015/09/11/roasted-vegetable-lasagna-recipe/72082180/
		Linked article: http://www.freep.com/story/life/food/recipes/2015/09/11/chicken-tamale-casserole/72083404/
		Linked article: http://www.freep.com/story/life/food/recipes/2015/09/11/butternut-squash-mac-cheese/72084214/
		Linked article: http://www.freep.com/story/life/food/recipes/2015/09/11/healthy-fall-meals-cabbage-soup/72082674/
	`)
}

func TestRecipeIntegration3(t *testing.T) {
	testRecipes(t, "http://www.freep.com/story/life/food/recipes/2015/07/01/herb-potato-salad/29524785/", `
		<p><b>Dorothy Sheets, 57, Farmington Hills</b></p>
		<p><b>Recipe:</b> Fresh Herb Potato Salad</p>
		<p>"It's great to take to Fourth of July potluck dinners because it has no mayo in it," Sheets says. "It's also nondairy and vegan."</p>
		<p><b>How it's made:</b> Sheets says she adapted this herbaceous potato salad recipe from a Free Press recipe several years ago. Cooked and cooled red skin potatoes are lightly coated with a mix of vinegar and olive oil. This salad gets a huge flavor boost from a blend of fresh herbs including parsley, basil, chives and oregano. Sheets notes that all fresh herb quantities can be adjusted to taste. "For basil and oregano, use the leaves only. For parsley, some tender stems are OK," Sheets says.</p>
		<p><b>Best aspect:</b> You can vary the herbs to your liking. Sheets has an herb garden and grabs whatever is good at the time. Also, not only does it work well for gatherings because it's safer, Sheets notes, "it's also lighter and healthier than one with mayo and eggs." Sheets suggests using a good-quality, fruity olive oil.</p>
		<p><b>Fresh Herb Potato Salad</b></p>
		<p><b>Serves:</b> 8 (1/2-cup servings) / <b>Preparation time:</b> 10 minutes / <b>Total time:</b> 40 minutes (plus chilling time)</p>
		<p><b>Ingredients</b></p>
		<p>4 cups red skin potatoes, cooked, cooled and cut into cubes</p>
		<p>1/4 cup cider vinegar</p>
		<p>2 tablespoons olive oil (use a good quality, fresh, fruity olive oil)</p>
		<p>1/3 to 1/2 cup chopped fresh parsley (curly or Italian are both fine)</p>
		<p>1/4 cup chopped fresh basil leaves</p>
		<p>2 tablespoons chopped fresh chives</p>
		<p>2 teaspoons chopped fresh oregano leaves</p>
		<p>2 tablespoons minced red onion</p>
		<p>2 garlic cloves, peeled, minced</p>
		<p>1/2 teaspoon salt</p>
		<p>1/4 teaspoon ground black pepper</p>
		<p>Pinch red pepper flakes (optional)</p>
		<p><b>Directions</b></p>
		<p>Place the cooked, cubed potatoes in a large bowl.</p>
		<p>In a measuring cup or small bowl, whisk together the vinegar, oil, parsley, basil, chives, oregano, onion, garlic, salt, black pepper and red pepper flakes, if using. Pour dressing over potatoes and stir gently. Cover and refrigerate at least 2 hours to allow flavors to blend. Serve.</p>
		<p><i>From Dorothy Sheets of Farmington Hills and tested by Susan Selasky for the Free Press Test Kitchen. Nutrition information not available. </i></p>	`, `
		Title: "Fresh Herb Potato Salad"
		Serving size: "8 (1/2-cup servings)"
		Total time: "40 minutes (plus chilling time)"
		Prep time: "10 minutes"
		I: "4 cups red skin potatoes, cooked, cooled and cut into cubes"
		I: "1/4 cup cider vinegar"
		I: "2 tablespoons olive oil (use a good quality, fresh, fruity olive oil)"
		I: "1/3 to 1/2 cup chopped fresh parsley (curly or Italian are both fine)"
		I: "1/4 cup chopped fresh basil leaves"
		I: "2 tablespoons chopped fresh chives"
		I: "2 teaspoons chopped fresh oregano leaves"
		I: "2 tablespoons minced red onion"
		I: "2 garlic cloves, peeled, minced"
		I: "1/2 teaspoon salt"
		I: "1/4 teaspoon ground black pepper"
		I: "Pinch red pepper flakes (optional)"
		D: "Place the cooked, cubed potatoes in a large bowl."
		D: "In a measuring cup or small bowl, whisk together the vinegar, oil, parsley, basil, chives, oregano, onion, garlic, salt, black pepper and red pepper flakes, if using. Pour dressing over potatoes and stir gently. Cover and refrigerate at least 2 hours to allow flavors to blend. Serve."
	`)
}

func TestRecipeIntegration4(t *testing.T) {
	testRecipes(t, "http://www.freep.com/story/life/food/recipes/2015/07/01/baked-potato-salad/29517797/", `
		<p><b>Theresa Makowski, 54, Canton </b></p>
		<p><b>Recipe:</b> Baked Potato Salad</p>
		<p>Makowski was intrigued when her sister put potatoes in the oven to bake for a potato salad she was going to make the next day. "She baked the potatoes but treated them just like a regular salad ingredient. I was disappointed. I really wanted to taste something different. So I came up with my own version."</p>
		<p><b>How it's made: </b>Makowski washes a few russets, rubs them with olive oil and lightly salts them with sea or kosher salt. The potatoes are wrapped in foil and bake until tender. Once done, the potatoes are cut up and mixed with extra-crispy bacon, green onions, cheese and some sour cream.</p>
		<p><b>Best aspects:</b> "It is awesome and goes great with a steak," Makowski says. "The really nice thing about serving the salad versus providing individual baked potatoes is twofold: There's very little waste (ever notice how people eat the top of the potato with the goodies but leave the bottom behind? not with this dish) and it's economical. You don't have to buy 10 potatoes for 10 guests. You can buy four and make a salad out of them and it goes much farther and people love it."</p>
		<p><b>Baked Potato Salad</b></p>
		<p><b>Serves:</b> 6 / <b>Preparation time:</b> 15 minutes / <b>Total time:</b> 1 hour, 30 minutes</p>
		<p><b>Ingredients</b></p>
		<p>5 to 6 large russet potatoes, scrubbed</p>
		<p>Olive oil</p>
		<p>Sea salt or kosher salt to taste</p>
		<p>12 ounces bacon</p>
		<p>16 ounces sour cream</p>
		<p>2 cups sharp cheddar</p>
		<p>Chopped green onion</p>
		<p><b>Directions</b></p>
		<p>Preheat the oven to 425 degrees.</p>
		<p>Rub the outside of the potatoes with olive oil and sprinkle with sea or kosher salt. Wrap the potatoes in foil and place in the oven. Bake until the potatoes are tender.</p>
		<p>Meanwhile, fry the bacon until it's extra-crispy. Crumble the bacon and set aside.</p>
		<p>When the potatoes are done and cool enough to handle, slice them in half and cut into bite-size chunks. Place the potatoes in a bowl and fold in the sour cream, bacon, sharp cheddar cheese and green onion.</p>
		<p><i>From Theresa Makowski of Canton and tested by Susan Selasky for the Free Press Test Kitchen. Nutrition information not available. </i></p>	`, `
		Title: "Baked Potato Salad"
		Serving size: "6"
		Total time: "1 hour, 30 minutes"
		Prep time: "15 minutes"
		I: "5 to 6 large russet potatoes, scrubbed"
		I: "Olive oil"
		I: "Sea salt or kosher salt to taste"
		I: "12 ounces bacon"
		I: "16 ounces sour cream"
		I: "2 cups sharp cheddar"
		I: "Chopped green onion"
		D: "Preheat the oven to 425 degrees."
		D: "Rub the outside of the potatoes with olive oil and sprinkle with sea or kosher salt. Wrap the potatoes in foil and place in the oven. Bake until the potatoes are tender."
		D: "Meanwhile, fry the bacon until it's extra-crispy. Crumble the bacon and set aside."
		D: "When the potatoes are done and cool enough to handle, slice them in half and cut into bite-size chunks. Place the potatoes in a bowl and fold in the sour cream, bacon, sharp cheddar cheese and green onion."
	`)
}

func TestRecipeIntegration5(t *testing.T) {
	testRecipes(t, "http://www.freep.com/story/life/food/recipes/2015/07/01/icebox-potato-salad-recipe/29515873/", `
		<p><b>Diana Balint, 60, Dearborn </b></p>
		<p><b>Recipe:</b> Icebox Potato Salad</p>
		<p>"This recipe is from my mother who was married in 1947 and lived in a small apartment in Detroit," Balint said. "She modified it from her mother who always kept space in the 'ice box' for the salad to sit for 24 hours." Balint also remembers her mother had a "patterned, cracked china cup that served as the measuring standard and it worked every time."</p>
		<p><b>How it's made:</b> Mayo- and egg-based, but with the addition of sour cream and mustard relish. "In 1947, there was no such thing as 'hot dog relish' so the ingredients would have been yellow mustard and dill relish," Balint says. This recipe does benefit from a good overnight chill.</p>
		<p><b>Best aspect:</b> The mustard relish lends tangy and vinegary tones. If you can't find mustard relish, use a blend of yellow mustard and dill relish.</p>
		<p><b>    <br>  </b></p>
		<p><b>    <br>  </b></p>
		<p><b>    <br>  </b></p>
		<p><b>    <br>  </b></p>
		<p><b>    <br>  </b></p>
		<p><b>    <b>Icebox Potato Salad</b>  </b></p>
		<p><b>Makes:</b> About 8 cups / <b>Preparation time:</b> 20 minutes / <b>Total time:</b> 1 hour (plus chilling times)</p>
		<p><b>Ingredients</b></p>
		<p>2 pounds red skin potatoes</p>
		<p>1/2 cup sour cream</p>
		<p>1 cup mayonnaise (or more for desired consistency)</p>
		<p>2 cups chopped celery or more as desired</p>
		<p>1 cup chopped green onions or more as desired</p>
		<p>1/3 cup hot dog relish (mustard and relish combined) or favorite relish</p>
		<p>6 hard-boiled eggs, 4 to chop and add to salad and 2 for decoration</p>
		<p>Salt to taste</p>
		<p>1/2  cup of chopped Spanish onion if more onion taste is desired.</p>
		<p><b>Directions</b></p>
		<p>Scrub the potatoes and then boil in salted water just until tender. Cool the potatoes in the refrigerator for 24 hours.</p>
		<p>Peel and slice cooked potatoes and add sour cream and mayonnaise as you are slicing potatoes or they will turn brown. Stir in the celery, green onion, relish, 4 chopped eggs, salt to taste, and, if using, chopped Spanish onion. Save two hard-boiled eggs to slice and place on top in a design. Cover with plastic wrap and place in refrigerator for 6-24 hours.</p>
		<p><i>From Diana Balint of Dearborn and tested by Susan Selasky for the Free Press Test Kitchen. Nutrition information not available. </i></p>
	`, `
		Title: "Icebox Potato Salad"
		Serving size: "About 8 cups"
		Total time: "1 hour (plus chilling times)"
		Prep time: "20 minutes"
		I: "2 pounds red skin potatoes"
		I: "1/2 cup sour cream"
		I: "1 cup mayonnaise (or more for desired consistency)"
		I: "2 cups chopped celery or more as desired"
		I: "1 cup chopped green onions or more as desired"
		I: "1/3 cup hot dog relish (mustard and relish combined) or favorite relish"
		I: "6 hard-boiled eggs, 4 to chop and add to salad and 2 for decoration"
		I: "Salt to taste"
		I: "1/2 cup of chopped Spanish onion if more onion taste is desired."
		D: "Scrub the potatoes and then boil in salted water just until tender. Cool the potatoes in the refrigerator for 24 hours."
		D: "Peel and slice cooked potatoes and add sour cream and mayonnaise as you are slicing potatoes or they will turn brown. Stir in the celery, green onion, relish, 4 chopped eggs, salt to taste, and, if using, chopped Spanish onion. Save two hard-boiled eggs to slice and place on top in a design. Cover with plastic wrap and place in refrigerator for 6-24 hours."
	`)
}

func TestRecipeIntegration6(t *testing.T) {
	if skipFailingRecipes {
		t.SkipNow()
	}
	// +intro
	testRecipes(t, "http://www.freep.com/story/life/food/recipes/2015/07/02/ribs-sweet-smokey-barbecuue/29622845/", `
		<p><b>Sweet and Smoky Bourbon Ribs</b></p>
		<p><b>Serves:</b> 12 / <b>Preparation time:</b> 1 hour (plus chilling times) / <b>Total time:</b> 4 hours (not all active time)</p>
		<p>This is the rib recipes that wins over my hungry guests consistently. Don't be put off by the length of this recipe, it will go together quickly. It takes a little patience in keeping the heat steady, but that's all.</p>
		<p><b>Ingredients</b></p>
		<p><b>RIBS</b></p>
		<p><ul>    <li>6 baby back rib racks about 1¼ pounds each (or equal weight of St. Louis-style ribs)</li>    <li>1 cup bourbon</li>  </ul></p>
		<p><b>RUB</b></p>
		<p><ul>    <li>3 tablespoons coarse salt</li>    <li>3 tablespoons packed brown sugar</li>    <li>3 tablespoons paprika</li>    <li>2 tablespoons ground black pepper</li>    <li>1 tablespoon garlic powder</li>    <li>1 teaspoon ground cumin</li>  </ul></p>
		<p>FOR THE GRILL</p>
		<p><ul>    <li>2 cups hickory wood chips</li>    <li>2 cups beer</li>  </ul></p>
		<p>SAUCE</p>
		<p><ul>    <li>Bourbon Barbecue Sauce (See note)</li>  </ul></p>
		<p><b>Directions</b></p>
		<p>Remove the thin membrane on the back (nonmeaty side) of the ribs, starting at one end of the rack and pulling toward the other. Unless you are using a rib rack, cut the ribs into four or five portions so that they will fit on the grill evenly. Place the ribs in a large roasting pan, and pour the bourbon over them. Chill for 30 minutes, turning the ribs often. Pour off and discard the bourbon.Meanwhile, whisk together all the rub ingredients in a small bowl. Sprinkle the rub mixture over both sides of the ribs. Refrigerate 1 hour.</p>
		<p>Place the wood chips in a medium bowl, and pour the beer over them. Let stand for 1 hour.</p>
		<p>Prepare a charcoal chimney starter full of briquettes. When they're gray and ash-covered, carefully turn out the hot coals onto one side of the grill for indirect cooking. If using a gas grill, follow the manufacturer's instructions for indirect cooking. This usually requires igniting all burners to preheat and then shutting off one burner if you have two or two burners if you have three.</p>
		<p>Remove 1 cup of the wood chips from the beer and drain. Scatter the chips over the coals. Fill a foil loaf pan halfway with water, and place opposite the coals. Place grill grate on grill.</p>
		<p>Arrange ribs on the grate above the loaf pan and away from the direct heat. Close the lid, positioning the top vent directly over the ribs. If you don't have a temperature gauge on grill, place the stem of a thermometer through this vent, with the gauge on the outside and the tip near the ribs, but not touching the meat or the grill rack; leave in place during cooking. Check the temperature after 5 minutes. It should register between 275 and 325 degrees. Adjust vents if needed by opening them wider to increase the heat or closing to decrease the heat. Leave any other vents closed.</p>
		<p>After 45 minutes, using the charcoal chimney starter set atop a nonflammable surface, heat 15 more charcoal briquettes until gray and ash-covered.</p>
		<p>When the temperature in the grill falls below 275 degrees, use oven mitts to lift off the upper rack with the ribs and place it on a nonflammable surface. Using tongs, add the additional charcoal briquettes to the grill. Drain the remaining 1 cup of wood chips, and sprinkle over the charcoal. Place the grate with the ribs back on the grill, cover and continue cooking until the ribs are very tender and the meat pulls away from the bones, about 45 minutes to 1 hour longer. During the last 15 minutes of grilling, brush on the sauce.</p>
		<p><b>Cook's note:</b> To prepare the Bourbon Barbecue Sauce, in a large heavy saucepan, mix together 2 cups ketchup, ½ cup mild-flavored molasses, 1⁄3 cup bourbon, 1/4 cup Dijon mustard, 3 tablespoons favorite hot pepper sauce, 2 tablespoons Worcestershire sauce, 2 teaspoons paprika, 1 teaspoon garlic powder and 1 teaspoon onion powder.</p>
		<p>Bring to a boil over medium heat, stirring occasionally. Reduce heat to medium-low, and simmer uncovered, stirring frequently, until the sauce thickens and the flavors blend, about 15 minutes. Makes about 2½ cups. This sauce can be made 1 week ahead, covered and chilled.</p>
		<p>Adapted from Bon Appetit. Tested by Susan M. Selasky for the Free Press Test Kitchen</p>
		<p>591 calories (54% from fat), 35 grams fat (13 grams sat. fat), 35 grams carbohydrate, 35 grams protein, 1,381 mg sodium, 137 mg cholesterol, 1 gram fiber.</p>
	`, `
		Title: "Sweet and Smoky Bourbon Ribs"
		Serving size: "12"
		Total time: "4 hours (not all active time)"
		Prep time: "1 hour (plus chilling times)"
		Nutrition data: "591 calories (54% from fat), 35 grams fat (13 grams sat. fat), 35 grams carbohydrate, 35 grams protein, 1,381 mg sodium, 137 mg cholesterol, 1 gram fiber."
		I: "RIBS"
		I: "6 baby back rib racks about 1¼ pounds each (or equal weight of St. Louis-style ribs)"
		I: "1 cup bourbon"
		I: "RUB"
		I: "3 tablespoons coarse salt"
		I: "3 tablespoons packed brown sugar"
		I: "3 tablespoons paprika"
		I: "2 tablespoons ground black pepper"
		I: "1 tablespoon garlic powder"
		I: "1 teaspoon ground cumin"
		I: "FOR THE GRILL"
		I: "2 cups hickory wood chips"
		I: "2 cups beer"
		I: "SAUCE"
		I: "Bourbon Barbecue Sauce (See note)"
		D: "This is the rib recipes that wins over my hungry guests consistently. Don't be put off by the length of this recipe, it will go together quickly. It takes a little patience in keeping the heat steady, but that's all."
		D: "Remove the thin membrane on the back (nonmeaty side) of the ribs, starting at one end of the rack and pulling toward the other. Unless you are using a rib rack, cut the ribs into four or five portions so that they will fit on the grill evenly. Place the ribs in a large roasting pan, and pour the bourbon over them. Chill for 30 minutes, turning the ribs often. Pour off and discard the bourbon.Meanwhile, whisk together all the rub ingredients in a small bowl. Sprinkle the rub mixture over both sides of the ribs. Refrigerate 1 hour."
		D: "Place the wood chips in a medium bowl, and pour the beer over them. Let stand for 1 hour."
		D: "Prepare a charcoal chimney starter full of briquettes. When they're gray and ash-covered, carefully turn out the hot coals onto one side of the grill for indirect cooking. If using a gas grill, follow the manufacturer's instructions for indirect cooking. This usually requires igniting all burners to preheat and then shutting off one burner if you have two or two burners if you have three."
		D: "Remove 1 cup of the wood chips from the beer and drain. Scatter the chips over the coals. Fill a foil loaf pan halfway with water, and place opposite the coals. Place grill grate on grill."
		D: "Arrange ribs on the grate above the loaf pan and away from the direct heat. Close the lid, positioning the top vent directly over the ribs. If you don't have a temperature gauge on grill, place the stem of a thermometer through this vent, with the gauge on the outside and the tip near the ribs, but not touching the meat or the grill rack; leave in place during cooking. Check the temperature after 5 minutes. It should register between 275 and 325 degrees. Adjust vents if needed by opening them wider to increase the heat or closing to decrease the heat. Leave any other vents closed."
		D: "After 45 minutes, using the charcoal chimney starter set atop a nonflammable surface, heat 15 more charcoal briquettes until gray and ash-covered."
		D: "When the temperature in the grill falls below 275 degrees, use oven mitts to lift off the upper rack with the ribs and place it on a nonflammable surface. Using tongs, add the additional charcoal briquettes to the grill. Drain the remaining 1 cup of wood chips, and sprinkle over the charcoal. Place the grate with the ribs back on the grill, cover and continue cooking until the ribs are very tender and the meat pulls away from the bones, about 45 minutes to 1 hour longer. During the last 15 minutes of grilling, brush on the sauce."
		D: "Cook's note: To prepare the Bourbon Barbecue Sauce, in a large heavy saucepan, mix together 2 cups ketchup, ½ cup mild-flavored molasses, 1⁄3 cup bourbon, 1/4 cup Dijon mustard, 3 tablespoons favorite hot pepper sauce, 2 tablespoons Worcestershire sauce, 2 teaspoons paprika, 1 teaspoon garlic powder and 1 teaspoon onion powder."
		D: "Bring to a boil over medium heat, stirring occasionally. Reduce heat to medium-low, and simmer uncovered, stirring frequently, until the sauce thickens and the flavors blend, about 15 minutes. Makes about 2½ cups. This sauce can be made 1 week ahead, covered and chilled."
	`)
}

func TestRecipeIntegration7(t *testing.T) {
	testRecipes(t, "http://www.freep.com/story/life/food/recipes/2014/11/20/cheese-soup-test-kitchen-recipe/19324313/", `
		<p>Since our flashback recipes from the 1930s "Home Makers' Handibook" by the Women's Service Bureau of the Detroit Free Press have been such a hit, we decided to offer recipes from another gem.</p>
		<p>The "Detroit Free Press Cookbook: A Collection of the Best Loved Recipes from the Free Press Tower Kitchen, " by Jeremy Iggers and the late Nettie Duffield, was published in 1984 and includes recipes for local favorites, staff favorites, lots of ethnic options and recipes featuring locally made products such as Vernors.</p>
		<p>We'll post a recipe from this cookbook each week on. Our hope is that you, too, will enjoy recipes that were enjoyed nearly 30 years ago.</p>
		<p>This recipe makes 2<sup>1</sup>/<sub>2</sub> quarts of soup.</p>
		<p><span class="-newsgate-element-cci-howto-begin"></span></p>
		<p><b>    <span class="-newsgate-paragraph-cci-howto-head">Canadian Cheese Soup</span>  </b></p>
		<p><b>Ingredients</b></p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span> 5 tablespoons butter or margarine</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span> 2 medium carrots, finely chopped</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span> 2 ribs celery, finely chopped</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span> 1 medium onion, finely chopped</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span><sup>1</sup>/<sub>2</sub> green pepper, seeded and finely chopped</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span> 4 to 5 mushrooms, cleaned and finely chopped</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span><sup>1</sup>/<sub>2</sub> cup cooked ham, finely chopped, if desired</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span><sup>1</sup>/<sub>2</sub> cup flour</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span> 2 tablespoons cornstarch</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span> 1 quart chicken broth</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span> 1 quart milk</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span><sup>1</sup>/<sub>2</sub> teaspoon paprika</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span><sup>1</sup>/<sub>4</sub> to <sup>1</sup>/<sub>2</sub> teaspoon cayenne</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span><sup>1</sup>/<sub>2</sub> teaspoon dry mustard, if desired</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span> 1 pound process sharp cheddar cheese, grated</p>
		<p><span class="-newsgate-macro-cci-z-sym-square-bullet">■</span> Salt and freshly ground black pepper, if desired</p>
		<p><b>Directions</b></p>
		<p>In large, heavy soup pot, melt butter or margarine; add carrots, celery, onion, green pepper, mushrooms and ham, if desired. Cook over medium heat until vegetables are crisp- tender, about 10 minutes, stirring occasionally. Do not brown.</p>
		<p>Add flour and cornstarch and cook, stirring constantly, about three minutes. Add broth and cook, stirring constantly, until slightly thickened. Add milk, paprika, cayenne and mustard, if desired. Add cheese gradually, stirring constantly, until cheese is melted. Do not allow soup to boil after cheese is added because it will curdle. Season to taste with salt and black pepper, if desired. Serve very hot.</p>
		<p><span class="-newsgate-element-cci-howto-end"></span></p>
	`, `
		Title: "Canadian Cheese Soup"
		Serving size: "2 1/2 quarts"
		I: "5 tablespoons butter or margarine"
		I: "2 medium carrots, finely chopped"
		I: "2 ribs celery, finely chopped"
		I: "1 medium onion, finely chopped"
		I: "1/2 green pepper, seeded and finely chopped"
		I: "4 to 5 mushrooms, cleaned and finely chopped"
		I: "1/2 cup cooked ham, finely chopped, if desired"
		I: "1/2 cup flour"
		I: "2 tablespoons cornstarch"
		I: "1 quart chicken broth"
		I: "1 quart milk"
		I: "1/2 teaspoon paprika"
		I: "1/4 to 1/2 teaspoon cayenne"
		I: "1/2 teaspoon dry mustard, if desired"
		I: "1 pound process sharp cheddar cheese, grated"
		I: "Salt and freshly ground black pepper, if desired"
		D: "In large, heavy soup pot, melt butter or margarine; add carrots, celery, onion, green pepper, mushrooms and ham, if desired. Cook over medium heat until vegetables are crisp- tender, about 10 minutes, stirring occasionally. Do not brown."
		D: "Add flour and cornstarch and cook, stirring constantly, about three minutes. Add broth and cook, stirring constantly, until slightly thickened. Add milk, paprika, cayenne and mustard, if desired. Add cheese gradually, stirring constantly, until cheese is melted. Do not allow soup to boil after cheese is added because it will curdle. Season to taste with salt and black pepper, if desired. Serve very hot."
	`)
}

func TestRecipeIntegration8(t *testing.T) {
	testRecipes(t, "http://www.freep.com/story/life/food/recipes/2014/09/25/wolfgang-puck-recipe/16157775/", `
		<p>If you try to cook fresh seasonal produce as much as possible, you might be puzzled to see artichokes in your market at this time of year. Aren’t they a spring vegetable?</p>
		<p>Well, yes, to give you the shortest answer. Springtime is when the biggest crops usually fill produce departments.</p>
		<p>But artichoke plants also typically deliver fall crops. In fact, depending on exactly where they’re grown (most come from various areas in California), they’re an almost year-round crop. So you could well find good artichokes in your market right now or in weeks to come.</p>
		<p>The vegetable’s robust size and shape, thorny petals and satisfyingly nutty flavor and smooth texture make artichokes feel like perfect things to eat in autumn. They’re also incomparably light and healthy — very low in fat, cholesterol-free, high in fiber and good sources of vitamin C and other nutrients. On top of that, they satisfy hunger because they take so long to eat.</p>
		<p>Of course, many  of those benefits can go out the window when they’re served in traditional ways. Typical accompaniments include lots of melted butter, vinaigrette, or mayonnaise for dipping, or bread crumb stuffings mixed with generous amounts of butter and cheese.</p>
		<p>That’s why I decided to take a fresher, lighter approach to cooking and serving whole artichokes in my recently published cookbook, “Wolfgang Puck Makes It Healthy.” The following recipe features whole steamed artichokes served with a version of Green Goddess dip made not with the usual mayonnaise and sour cream but with nonfat plain Greek yogurt.</p>
		<p>Whenever possible, I cook whole artichokes in a pressure cooker, which steams them extra quickly — just 10 minutes under pressure. But if you want to cook them on the stovetop instead, simply put them with the other ingredients into a large nonreactive pot, bring to a boil, and simmer until tender enough for a leaf from the center to pull out easily after 30 to 45 minutes, depending on their size.</p>
		<p>If you’ve never eaten an artichoke before, start by pulling a petal from near the bottom. Dunk its fleshy base into the dip and then scrape off the flesh between your teeth before discarding the petal. Continue working round and round and up the artichoke. When you get to the fuzzy “choke” at the center, scrape it out with a spoon and discard it. Then, use a fork and knife to cut up and eat the artichoke’s heart, dipping each bite into the dressing.</p>
		<p>That may sound like a lot of work to eat a vegetable. But you’ll find it surprisingly enjoyable. And the results are so delicious, and so light, that you’ll want to go on making this recipe again and again, whatever the season.</p>
		<p><b>Pressure-cooker Steamed Whole Artichokes with Green Goddess Dip</b></p>
		<p><b>Serves:</b> 4</p>
		<p>4 large whole fresh artichokes</p>
		<p>1 lemon, cut into 4 center slices each ¼-inch thick, remaining cut ends reserved</p>
		<p>1 cup dry white wine</p>
		<p>1 cup organic vegetable broth, or water</p>
		<p>½ tablespoon whole coriander seeds</p>
		<p>½ tablespoon whole black peppercorns</p>
		<p>1 bay leaf</p>
		<p>1 cup low-fat Green Goddess Dressing</p>
		<p>First, trim the artichokes: With a sharp serrated knife, cut off the stem ends and a little bit of the base to form a flat bottom on each artichoke. Steadying each artichoke on its side on a cutting board, slice off the top third of the narrower petal end. With kitchen shears, snip off any remaining sharp petal tips. With the reserved lemon ends, gently rub all the cut edges on each artichoke to keep them from oxidizing.</p>
		<p>In a pressure cooker, combine the wine, broth or water, coriander, peppercorns and bay leaf. Stand the artichokes upright, side by side, inside the cooker. Place a lemon slice on top of each artichoke.</p>
		<p>Secure the lid on the pressure cooker. Bring to full pressure and, once pressure has been reached, cook the artichokes for 10 minutes, until tender.</p>
		<p>Release the pressure using the quick-release valve. Using a long-handled spoon, carefully remove the artichokes to a platter or individual plates and leave them to cool slightly.</p>
		<p>Serve the artichokes hot, warm, or chilled, accompanied by the Low-Fat Green Goddess Dressing in individual ramekins or bowls for dipping. Be sure to provide side plates or bowls, or a large communal one, for discarding petals as eaten.</p>
		<p>From and tested by Wolfgang Puck. Nutrition information not available.</p>
		<p><b>Low-fat Green Goddess Dressing</b></p>
		<p><b>Makes: </b>About 2½ cups</p>
		<p>1¼ cups nonfat plain Greek yogurt</p>
		<p>½ cup packed baby spinach leaves</p>
		<p>2 tablespoons packed chopped fresh flat-leaf parsley leaves</p>
		<p>2 tablespoons packed chopped fresh basil leaves</p>
		<p>2 tablespoons packed chopped fresh chives</p>
		<p>2 tablespoons packed chopped fresh chervil leaves</p>
		<p>2 tablespoons fresh lemon juice</p>
		<p>½ ripe  avocado, pitted</p>
		<p>1 garlic clove, coarsely chopped</p>
		<p>Kosher salt</p>
		<p>Freshly ground white pepper</p>
		<p>In a blender, combine the yogurt, spinach, parsley, basil, chives, chervil and lemon juice. With a tablespoon, scoop the avocado flesh out of the skin into the blender. Add the garlic and a little salt and pepper to taste.</p>
		<p>Blend the ingredients, pulsing the machine on and off and stopping as necessary to scrape down the bowl with a spatula, until a smooth dressing forms. Taste the mixture and, if necessary, pulse in a little more salt and pepper to taste.</p>
		<p>Transfer the dressing to a nonreactive container. Cover and refrigerate until ready to use. Serve within 3 to 4 days.</p>
		<p>From and tested by Wolfgang Puck. Nutrition information not available.</p>
	`, `
		Title: "Pressure-cooker Steamed Whole Artichokes with Green Goddess Dip"
		Serving size: "4"
		I: "4 large whole fresh artichokes"
		I: "1 lemon, cut into 4 center slices each ¼-inch thick, remaining cut ends reserved"
		I: "1 cup dry white wine"
		I: "1 cup organic vegetable broth, or water"
		I: "½ tablespoon whole coriander seeds"
		I: "½ tablespoon whole black peppercorns"
		I: "1 bay leaf"
		I: "1 cup low-fat Green Goddess Dressing"
		D: "First, trim the artichokes: With a sharp serrated knife, cut off the stem ends and a little bit of the base to form a flat bottom on each artichoke. Steadying each artichoke on its side on a cutting board, slice off the top third of the narrower petal end. With kitchen shears, snip off any remaining sharp petal tips. With the reserved lemon ends, gently rub all the cut edges on each artichoke to keep them from oxidizing."
		D: "In a pressure cooker, combine the wine, broth or water, coriander, peppercorns and bay leaf. Stand the artichokes upright, side by side, inside the cooker. Place a lemon slice on top of each artichoke."
		D: "Secure the lid on the pressure cooker. Bring to full pressure and, once pressure has been reached, cook the artichokes for 10 minutes, until tender."
		D: "Release the pressure using the quick-release valve. Using a long-handled spoon, carefully remove the artichokes to a platter or individual plates and leave them to cool slightly."
		D: "Serve the artichokes hot, warm, or chilled, accompanied by the Low-Fat Green Goddess Dressing in individual ramekins or bowls for dipping. Be sure to provide side plates or bowls, or a large communal one, for discarding petals as eaten."

		Title: "Low-fat Green Goddess Dressing"
		Serving size: "About 2½ cups"
		I: "1¼ cups nonfat plain Greek yogurt"
		I: "½ cup packed baby spinach leaves"
		I: "2 tablespoons packed chopped fresh flat-leaf parsley leaves"
		I: "2 tablespoons packed chopped fresh basil leaves"
		I: "2 tablespoons packed chopped fresh chives"
		I: "2 tablespoons packed chopped fresh chervil leaves"
		I: "2 tablespoons fresh lemon juice"
		I: "½ ripe avocado, pitted"
		I: "1 garlic clove, coarsely chopped"
		I: "Kosher salt"
		I: "Freshly ground white pepper"
		D: "In a blender, combine the yogurt, spinach, parsley, basil, chives, chervil and lemon juice. With a tablespoon, scoop the avocado flesh out of the skin into the blender. Add the garlic and a little salt and pepper to taste."
		D: "Blend the ingredients, pulsing the machine on and off and stopping as necessary to scrape down the bowl with a spatula, until a smooth dressing forms. Taste the mixture and, if necessary, pulse in a little more salt and pepper to taste."
		D: "Transfer the dressing to a nonreactive container. Cover and refrigerate until ready to use. Serve within 3 to 4 days."
	`)
}

func TestRecipeIntegration9(t *testing.T) {
	testRecipes(t, "http://www.freep.com/story/life/food/recipes/2015/06/25/let-others-worry-about-the-burgers-you-focus-on-the-sangria/29161075/", `
		<p>Want to be the hero of your Fourth of July gathering? Leave the burgers and dogs to somebody else. Ditto for the potato and pasta salads. What you want to bring is the sangria. Because it's hard to go wrong at an outdoor summer party when you're the one toting the pitcher cocktail.</p>
		<p>Still, I'm not a big believer in working hard for my cocktail. So this recipe is a breeze to assemble. Just dump and stir in the morning, then let it chill for a few hours before serving. Whatever you do, don't add ice until it's in the glass, and even then keep it to one or two cubes at most. Nobody wants a watered-down cocktail.</p>
		<p>For this recipe, I call for cava — the sparkling wine of Spain — but feel free to substitute the bubbles of your choice. Or if you'd rather cut the alcohol a little (can't imagine why, but whatever), ginger beer or a lemon-lime soda are fine substitutes.</p>
		<p><b>Peach And Raspberry Sangria With Cava And Strawberry Ice<br></b></p>
		<p><b>Servings:</b> 10 / <b>Preparation time: </b>10 minutes active / <b>Total time:</b>  2 to 4 hours chilling</p>
		<p>1 cup brandy</p>
		<p>1 cup peach juice</p>
		<p>1/2 cup simple syrup or agave syrup</p>
		<p>1 bottle, 750-milliliters  dry red wine (such as rioja)</p>
		<p>6 ounces fresh raspberries</p>
		<p>2 oranges, thinly sliced</p>
		<p>2 limes, thinly sliced</p>
		<p>1 bag, (16-ounces)  frozen strawberries</p>
		<p>3/4 cup orange juice</p>
		<p>1/4 cup sugar</p>
		<p>1 bottle, (750-milliliter bottle) cava (or other sparkling wine)</p>
		<p>In a large pitcher, stir together the brandy, peach juice and syrup until the syrup is dissolved. Add the wine and stir again. Stir in the raspberries, oranges and limes, then cover and refrigerate for 2 to 4 hours.</p>
		<p>Meanwhile, in a blender combine the strawberries, orange juice and sugar. Puree until very smooth. Pour into 2 ice cube trays, then freeze for 2 to 4 hours, or until solid.</p>
		<p>When ready to serve, slowly pour the cava into the pitcher. Stir once or twice gently just to mix. Pour into serving glasses, then add 1 to 2 frozen strawberry cubes to each glass.</p>
		<p><i>From and tested by the Associated Press. Nutrition information not available. </i></p>
	`, `
		Title: "Peach And Raspberry Sangria With Cava And Strawberry Ice"
		Total time: "2 to 4 hours chilling"
		Prep time: "10 minutes active"
		I: "1 cup brandy"
		I: "1 cup peach juice"
		I: "1/2 cup simple syrup or agave syrup"
		I: "1 bottle, 750-milliliters dry red wine (such as rioja)"
		I: "6 ounces fresh raspberries"
		I: "2 oranges, thinly sliced"
		I: "2 limes, thinly sliced"
		I: "1 bag, (16-ounces) frozen strawberries"
		I: "3/4 cup orange juice"
		I: "1/4 cup sugar"
		I: "1 bottle, (750-milliliter bottle) cava (or other sparkling wine)"
		D: "In a large pitcher, stir together the brandy, peach juice and syrup until the syrup is dissolved. Add the wine and stir again. Stir in the raspberries, oranges and limes, then cover and refrigerate for 2 to 4 hours."
		D: "Meanwhile, in a blender combine the strawberries, orange juice and sugar. Puree until very smooth. Pour into 2 ice cube trays, then freeze for 2 to 4 hours, or until solid."
		D: "When ready to serve, slowly pour the cava into the pitcher. Stir once or twice gently just to mix. Pour into serving glasses, then add 1 to 2 frozen strawberry cubes to each glass."
	`)
}

func TestRecipeIntegration10(t *testing.T) {
	// +intro
	testRecipes(t, "http://www.freep.com/story/life/food/recipes/2015/08/28/slow-cooker-sriracha-beans/71353992/", `
		<p><span class="-newsgate-macro-cci-drop-initial-"></span>Eating a healthy diet that includes tomato products is one step men can take to reduce prostate cancer, according to the American Institute for Cancer Research. That’s good news, since prostate cancer is the second most common cancer among American men.</p>
		<p>One disease-fighting agent in tomato products that has received a lot of attention is a phytochemical (beneficial compound in plants) called lycopene. It acts as an antioxidant, protecting against cancer by preventing cell damage in the body. Lycopene is the pigment that gives tomatoes, pink grapefruit, watermelon, guavas and papayas their red color.</p>
		<p>While research results of this potential link are mixed, some studies show a reduced prostate cancer risk with higher intakes of dietary lycopene.</p>
		<p>Interestingly, the lycopene from processed tomatoes (such as canned, sauce, paste, soup, juice, and ketchup) is more readily absorbed by the body than the lycopene from raw tomatoes. Today’s recipe features lycopene-rich ketchup and tomato paste.</p>
		<p>Lycopene is fat-soluble, meaning you absorb more of it when you add a little fat. So if you are enjoying our lycopene-filled slow-cooker beans, the canola oil used to sauté the onions will actually help you absorb more of the lycopene found in tomato products.</p>
		<p><span class="-newsgate-paragraph-cci-endnote-contact-">Darlene Zimmerman is a registered dietitian in Henry Ford Hospital’s Heart &amp; Vascular Institute. For questions about today’s recipe, call 313-972-1920.</span></p>
		<p><span class="-newsgate-element-cci-howto--begin"></span></p>
		<p><span class="-newsgate-paragraph-cci-howto-head-">    <b>Slow-Cooker Sriracha Beans</b>  </span></p>
		<p><b>Makes:</b> 11 servings (<sup>1</sup>/<sub>2</sub> cup each) / <b>Preparation time:</b> 20 minutes  (plus overnight soaking of beans / <b>Total time:</b> 5 hours, 30 minutes using a slow-cooker</p>
		<p>For best results it’s best to soak dry beans in water overnight according to package directions before simmering and then adding them to the slow cooker.</p>
		<p><span class="-newsgate-paragraph-cci-howto-components-">2 cups dry white beans (Great Northern or navy)</span></p>
		<p><span class="-newsgate-paragraph-cci-howto-components-">2 cups chopped onion</span></p>
		<p><span class="-newsgate-paragraph-cci-howto-components-">    <sup>1</sup>/<sub>4</sub> cup real bacon bits</span></p>
		<p><span class="-newsgate-paragraph-cci-howto-components-">1 tablespoon canola oil</span></p>
		<p><span class="-newsgate-paragraph-cci-howto-components-">2 cups water</span></p>
		<p><span class="-newsgate-paragraph-cci-howto-components-">    <sup>1</sup>/<sub>2</sub> cup ketchup</span></p>
		<p><span class="-newsgate-paragraph-cci-howto-components-">6 tablespoons tomato paste</span></p>
		<p><span class="-newsgate-paragraph-cci-howto-components-">    <sup>1</sup>/<sub>3</sub> cup packed brown sugar</span></p>
		<p><span class="-newsgate-paragraph-cci-howto-components-">    <sup>1</sup>/<sub>3</sub> cup maple syrup</span></p>
		<p><span class="-newsgate-paragraph-cci-howto-components-">3 tablespoons Worcestershire sauce, divided</span></p>
		<p><span class="-newsgate-paragraph-cci-howto-components-">2 teaspoons dry mustard</span></p>
		<p><span class="-newsgate-paragraph-cci-howto-components-">3 tablespoons cider vinegar</span></p>
		<p><span class="-newsgate-paragraph-cci-howto-components-">1<sup>1</sup>/<sub>2</sub> teaspoons sriracha  sauce</span></p>
		<p>After soaking beans overnight, drain liquid. Return beans to pot, add enough water to cover beans by an inch. Simmer 30 minutes.</p>
		<p>While beans are cooking, sauté onion and bacon bits in oil, in a skillet, over medium heat. Drain beans and add to a 5-quart slow cooker along with onion mixture, water, ketchup, tomato paste, brown sugar, maple syrup, 2 tablespoons Worcestershire sauce and dry mustard. Cover and cook on high for 4 to 5 hours.</p>
		<p>When cooking is complete, stir in remaining Worcestershire sauce, vinegar and sriracha sauce.</p>
		<p><span class="-newsgate-paragraph-cci-howto-volume-">Created by Darlene Zimmerman, MS, RD, and tested by Susan Selasky for the Free Press Test Kitchen. Analysis per <sup>1</sup>/<sub>2</sub> cup serving.</span></p>
		<p><span class="-newsgate-paragraph-cci-howto-volume-">217 <b>calories</b> (9% from <b>fat</b>), 2 grams <b>fat</b> (0 grams <b>sat. fat</b>, 0 grams <b>trans fat</b>), 40 grams <b>carbohydrates</b>, 9 grams <b>protein</b>, 266 mg <b>sodium</b>, 2 mg <b>cholesterol</b>, 70 mg <b>calcium</b>, 8 grams <b>fiber</b>. Food exchanges: 2 starch, 2 vegetable.</span></p>
		<p><span class="-newsgate-element-cci-howto--end"></span></p>
	`, `
		Title: "Slow-Cooker Sriracha Beans"
		Serving size: "11 servings (1/2 cup each)"
		Total time: "5 hours, 30 minutes using a slow-cooker"
		Prep time: "20 minutes (plus overnight soaking of beans"
		Nutrition data: "217 calories (9% from fat), 2 grams fat (0 grams sat. fat, 0 grams trans fat), 40 grams carbohydrates, 9 grams protein, 266 mg sodium, 2 mg cholesterol, 70 mg calcium, 8 grams fiber. Food exchanges: 2 starch, 2 vegetable."
		I: "2 cups dry white beans (Great Northern or navy)"
		I: "2 cups chopped onion"
		I: "1/4 cup real bacon bits"
		I: "1 tablespoon canola oil"
		I: "2 cups water"
		I: "1/2 cup ketchup"
		I: "6 tablespoons tomato paste"
		I: "1/3 cup packed brown sugar"
		I: "1/3 cup maple syrup"
		I: "3 tablespoons Worcestershire sauce, divided"
		I: "2 teaspoons dry mustard"
		I: "3 tablespoons cider vinegar"
		I: "1 1/2 teaspoons sriracha sauce"
		D: "For best results it’s best to soak dry beans in water overnight according to package directions before simmering and then adding them to the slow cooker."
		D: "After soaking beans overnight, drain liquid. Return beans to pot, add enough water to cover beans by an inch. Simmer 30 minutes."
		D: "While beans are cooking, sauté onion and bacon bits in oil, in a skillet, over medium heat. Drain beans and add to a 5-quart slow cooker along with onion mixture, water, ketchup, tomato paste, brown sugar, maple syrup, 2 tablespoons Worcestershire sauce and dry mustard. Cover and cook on high for 4 to 5 hours."
		D: "When cooking is complete, stir in remaining Worcestershire sauce, vinegar and sriracha sauce."
	`)
}

/* template

func TestRecipeIntegration10(t *testing.T) {
	testRecipes(t, "URL", `
	`, `
		XXX
	`)
}
*/

func testRecipes(t *testing.T, url string, html string, expected string) {
	t.Logf("Testing URL: %v", url)

	var extract *m.ExtractedBody
	if strings.TrimSpace(html) == "" {

		processor := a.ParseArticleAtURL(url, true)
		extract = processor.ExtractedBody
		if processor.Err != nil {
			t.Fatalf("Failed to parse article: %v", processor.Err)
		}

		println("Here's the HTML to embed for", url)
		println(processor.Html)
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
		ds := diff(expected, actual)
		ds = strings.Join(addPrefixToLines(fixLineColons(strings.Split(ds, "\n")), "> "), "\n")
		t.Errorf("Got (diff'd against expected value):\n%s", ds)
	}
}

func trimLines(input []string) []string {
	result := make([]string, 0, len(input))
	for _, el := range input {
		result = append(result, strings.TrimSpace(el))
	}
	return result
}

func addPrefixToLines(input []string, prefix string) []string {
	result := make([]string, 0, len(input))
	for _, el := range input {
		result = append(result, prefix+el)
	}
	return result
}

func fixLineColons(input []string) []string {
	result := make([]string, 0, len(input))
	for _, el := range input {
		el = strings.Replace(el, ":", " →", 1)
		result = append(result, el)
	}
	return result
}

func diff(a, b string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(a, b, true)
	if len(diffs) > 2 {
		diffs = dmp.DiffCleanupSemantic(diffs)
		diffs = dmp.DiffCleanupEfficiency(diffs)
	}
	return diffsToString(diffs)
}

func diffsToString(diffs []diffmatchpatch.Diff) string {
	var buff bytes.Buffer
	for _, diff := range diffs {
		text := diff.Text
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			buff.WriteString("(++")
			buff.WriteString(text)
			buff.WriteString("++)")
		case diffmatchpatch.DiffDelete:
			buff.WriteString("(~~")
			buff.WriteString(text)
			buff.WriteString("~~)")
		case diffmatchpatch.DiffEqual:
			buff.WriteString(text)
		}
	}
	return buff.String()
}
