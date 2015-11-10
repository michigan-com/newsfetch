package recipe_parsing

import (
	fuz "github.com/michigan-com/newsfetch/extraction/fuzzy_classifier"
)

func NewIngredientClassifier() *fuz.Classifier {
	classifier := fuz.NewFuzzyClassifier()
	classifier.AddOrPanic(ingredientRules)
	return classifier
}

var ingredientRules = `
    :@not_ingredient
    Healthy Table
    In Good Taste
    match made in heaven

    :@ingredient
    .skip +b +a @alternative @optional @purpose_clause @extra_reference
    @ingredient_list
    @ingredient_list @strict_quantity

    :@ingredient_list
    .skip +a @plus_more
    @lose_ingredient_alternative
    @sole_ingredient_component
    +@strict_ingredient_alternative @lose_ingredient_alternative
    +@strict_ingredient_alternative and @lose_ingredient_alternative ?@mix
    @lose_ingredient_alternative plus @lose_ingredient_component

    :@mix
    mix
    salad

    :@lose_ingredient_alternative
    @lose_ingredient_component
    @lose_ingredient_component or @lose_ingredient_component

    :@strict_ingredient_alternative
    @strict_ingredient_component
    @lose_ingredient_component or @strict_ingredient_component

    :@sole_ingredient_component
    .skip +a @adjective @adjective_or_adjective
    @lose_ingredient_component
    
    :@lose_ingredient_component
    .skip +a @to_taste @extra_qualities
    @quantity @lose_ingredient_subcomponent
    @quantity_number_with_possible_range @ingredient_name_with_adjectives @quantity_unit ?@postfix_rep
    @lose_ingredient_subcomponent
    @common_ingredient_component ?@postfix_rep
    
    :@strict_ingredient_component
    .skip +a @to_taste @extra_qualities
    @quantity @strict_ingredient_subcomponent
    @strict_ingredient_subcomponent
    @common_ingredient_component

    :@common_ingredient_component
    @portion_names @of_or_from ?@quantity @strict_ingredient_subcomponent
    corn kernels cut from ?@quantity @strict_ingredient_subcomponent
    ?@squeeze @strict_ingredient_subcomponent juice
    @strict_ingredient_subcomponent slices
    @seasoning

    :@lose_ingredient_subcomponent
    @ingredient_name_with_adjectives ?@postfix_rep

    # no suffixes
    :@strict_ingredient_subcomponent
    @ingredient_name_with_adjectives

    :@ingredient_name_with_adjectives
    .skip +b @adjective @adjective_or_adjective @the
    @ingredient_name
    @corner_case_ingredient_name

    :@alternative
    or store-bought
    or @alternative_with_adjectives
    or another herb
    @alternative_with_adjectives or

    :@alternative_with_adjectives
    @seasoning

    :@seasoning
    .skip +b @adjective
    all-purpose seasoning
    all-purpose seasoning mix
    all-purpose seasoning blend

    :@squeeze
    squeeze of

    :@of_or_from
    of
    from

    :@to_taste
    to taste
    # typo?
    taste
    as desired

    :@plus_more
    plus more to taste
    plus more for serving
    plus more for dusting
    plus more as needed

    :@purpose_clause
    for @purpose_name
    as needed to soften tortillas
    to adjust consistency
    to cover

    :@purpose_name
    $cooking
    $garnish
    $desting
    $dusting
    $serving
    $sprinkling
    $frying
    work $surface
    $greasing pan
    $shaping the dough
    $brushing on top
    griddle
    ?the $grill grates
    ?the $grill
    deep-frying
    ?the $griddle

    :@optional
    optional
    if desired

    :@extra_reference
    see cook's note
    see note

    :@portion_names
    .skip @or_and
    +@portion_name

    :@portion_name
    curls
    julienne
    ?reserved juice
    ?@zest_adjective zest
    kernel
    leaves
    blend

    :@zest_adjective
    ?finely grated

    :@corner_case_ingredient_name
    mixed $greens with a light vinaigrette
    sauteed greens such as kale or mustard greens
    Queso Chihauhua
    queso fresco
    Pastry for a double pie crust
    One 12 to 14-pound fresh turkey or frozen thawed
    bamboo or metal skewers

    :@ingredient_name
    #.multi
    seasoning blend
    hot sauce
    other citrus type of seasoning
    ?nonstick @oil_kinds oil ?cooking spray
    @oil_kinds oil-flavored ?nonstick ?cooking spray
    ?nonstick cooking spray
    nonstick spray
    ?nonstick ?floured baking spray
    @oil_kinds oil
    oil
    condiment
    saffron
    poppy seeds
    sesame seeds
    tarragon
    garlic
    ##
    cream of tartar
    radish
    sea salt
    salt
    sanding sugar
    sugar
    Flaky Maldon salt
    flour
    all-purpose flour
    cornmeal
    ?@pepper_kind pepper
    ?@pepper_kind pepper flakes
    caramel sauce
    barbecue sauce
    ##
    lemon zest
    lemon pepper
    orange rind
    lime wedges
    ##
    nutmeg
    chives
    onion
    onion rings
    caramelized onions
    frosting
    balsamic syrup
    cilantro
    cilantro leaves
    ?@butter_type_rep butter
    club soda
    basil
    corn
    corn on the cob
    corn chips
    peppercorn
    candy corn
    sunchoke chips
    tortilla chips
    tortilla chips or strips
    chips
    cumin
    tortilla
    paprika
    Bourbon Barbecue Sauce
    salsa
    croutons
    decorating gel
    gummy worms
    ?@rice_kind rice
    ?@rice_kind rice cereal
    ?@chocolate_type_rep chocolate
    cinnamon sticks
    noodles
    @fruit
    @vegetables
    @greens
    @dairy
    @meat
    @fish
    water
    guacamole
    mashed potatoes
    potatoes
    risotto
    bread
    pita bread
    parchment paper
    pea tendrils
    pea shoots
    golden beets
    vinaigrette
    ?@mushroom_type_rep mushroom
    egg
    egg yolk
    egg white

    :@greens
    lettuce
    salad greens
    @green_sprigs ?sprigs
    sprigs of @green_sprigs
    herb ?sprigs ?@green_sprigs_rep
    sage leaves
    sage
    herbs
    romaine hearts
    ?fresh herbs such as @herbs_rep

    :@herbs_rep
    .skip @or_and
    +@herbs

    :@herbs
    basil
    dill
    sage
    @green_sprigs

    :@green_sprigs_rep
    .skip @or_and
    +@green_sprigs

    :@green_sprigs
    parsley
    thyme
    rosemary
    oregano

    :@vegetables
    avocado
    couscous
    tomato
    cucumber

    :@fruit
    lemon
    lime
    orange
    grapefruit
    seasonal berries
    berries
    watermelon
    maraschino cherries
    cherries

    :@dairy
    milk
    cream
    creme fraiche
    sour cream
    ?@cheese_kinds cheese
    @cheese_kind

    :@meat
    back pork ribs
    racks

    :@fish
    ahi tuna
    tuna

    :@pepper_kind
    cayenne

    :@cheese_kinds
    .skip @or_and
    +@cheese_kind

    :@cheese_kind
    Parmesan
    Romano

    :@oil_kinds
    .skip @or_and
    +@oil_kind

    :@oil_kind
    vegetable
    lard
    canola
    ?@olive_oil_kind olive

    :@olive_oil_kind
    extra-virgin    

    :@rice_kind
    basmati
    puffed

    :@mushroom_type_rep
    .skip @or_and
    +@mushroom_type

    :@mushroom_type
    portabella

    :@chocolate_type_rep
    .skip @or_and
    +@chocolate_type

    :@chocolate_type
    $milk

    :@butter_type_rep
    .skip @or_and
    +@butter_type

    :@butter_type
    ?@peanut_butter_type_rep peanut
    flavored

    :@peanut_butter_type_rep
    .skip @or_and
    +@peanut_butter_type

    :@peanut_butter_type
    creamy
    crunchy

    :@adjective_or_adjective
    @adjective or @adjective

    #############################
    :@adjective
    favorite
    homemade
    kosher
    additional
    butter-flavored
    from above

    small
    medium
    large

    fine
    coarse

    heavy
    mild

    white
    green
    black
    white
    red
    brown

    sushi grade
    good-quality

    baby
    confectioners
    flat-leaf
    crispy
    cage-free
    ##
    sweet
    semisweet
    ##### 
    @processed
    ##### countries
    french
    italian
    thai
    mexican blend

    :@postfix_rep
    .skip @or_and
    +@postfix

    :@postfix
    @postfix_adjective
    @processed

    :@postfix_adjective
    at room temperature

    :@processed
    very $thinly sliced
    ?$thinly sliced
    ?$freshly $ground
    ?$freshly $grated
    ?$freshly $chopped
    $fresh $chopped
    ?$freshly $cracked
    $crushed
    $grilled
    $fresh
    $snipped
    $salted
    $fortified
    $whipped
    $shredded
    $fried
    $pickled
    $cooked
    $steamed
    $crumbled
    $minced
    $diced
    $peeled
    $softened

    :@extra_qualities
    spareribs
    membranes removed
    @quantity each

    :@or_and
    or
    and

    :@or
    or

    :@and
    and

    :@the
    the

    :@quantity
    @quantity_unit
    @strict_quantity

    :@strict_quantity
    ?about @standalone_quantity_number
    ?about @quantity_number_with_possible_range @quantity_unit ?of
   
    :@standalone_quantity_number
    ?a pinch ?of
    generous pinch ?of
    dash
    a few sprigs ?of
    a few drops ?of
    a few sprigs ?of
    several dashes ?of
    small bunch ?of
    @quantity_number_with_possible_range

    :@quantity_number_with_possible_range
    @quantity_number_with_possible_unit
    @quantity_range

    :@quantity_range
    @quantity_number_with_possible_unit to @quantity_number_with_possible_unit

    :@quantity_number_with_possible_unit
    one
    two
    half ?a
    @number

    :@quantity_unit
    cup
    sprig
    pounds
    bunch
    ears
    teaspoon
    tablespoon
    clove
    caps
    stick
    ounce
`
