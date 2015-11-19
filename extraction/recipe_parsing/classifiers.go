package recipe_parsing

import (
	fuz "github.com/michigan-com/newsfetch/extraction/fuzzy_classifier"
)

func NewIngredientClassifier() *fuz.Classifier {
	classifier := fuz.NewFuzzyClassifier()
	classifier.AddOrPanic(ingredientRules)
	return classifier
}

func NewDirectionClassifier() *fuz.Classifier {
	classifier := fuz.NewFuzzyClassifier()
	classifier.AddOrPanic(directionRules)
	return classifier
}

var directionRules = `
    :@direction
    @action
    need to sit
    @note
    @extra
    @nonaction

    :@action
    $add
    $stir
    $blend
    $combine
    $adjust seasoning
    $boil
    lightly $oil
    $place
    $serve
    $divide
    $pat
    $remove
    $preheat
    $whisk
    $slice
    $arrange
    $grind
    $set aside
    $bake
    $line
    $spray
    $use
    $season
    $return
    $reduce
    $put
    $brush
    $have ready
    $heat
    $pour
    $press
    $refrigerate
    $grill
    $beat
    $increase
    $continue
    $continuing
    $mix together
    $transfer
    $soak
    $allow
    $garnish with
    $reserve
    $crush
    $toss
    $roast
    $toast
    $coat
    $shave
    $cover
    $chill
    $trim
    $flip
    $fill
    $roll
    $prepare
    $cut out
    $cut through
    $cut @subject
    $scrape
    $turn
    $drop
    $sprinkle
    $check
    $grease
    $slide
    $sift
    $steam
    $rub
    $marinate
    $broil
    $spread
    $thaw
    $cook
    $break
    $melt
    $dry-age
    $fold
    $tie
    $make
    @confusing_action

    # need some kind of assertion (like “starts a sentance” or “preceded by an adverb” or “after a comma”)
    :@confusing_action
    $mix
    $warm
    $top @subject
    $cream @subject
    gently $spoon
    $spoon @subject
    #$cut

    :@nonaction
    in a small bowl

    :@note
    cook's note
    nutritional analysis
    analysis per
    analysis based on
    analysis is for
    analysis without
    if you can't find
    if you don't have
    in this recipe
    can cut this recipe
    this recipe
    is great with
    we used
    can be made
    keeps it
    for more of @a
    mix in well
    come together
    is available
    is a twist
    is great as is
    $allows
    $try
    can be served
    $look for
    can make
    day in advance
    day ahead
    can easily double
    great way
    simple to make
    $served as
    $served with
    can substitute
    do not substitute
    make for a nice
    $change up
    inspired by
    ideal for
    to brine
    grocery stores
    easiest to make
    can be prepared
    is the secret
    this happens when
    another thing I discovered

    :@extra
    adapted from
    from the
    from @cap
    # From "365 Ways to Cook Eggs"
    from @number
    from chef
    from brothers
    from a @cap
    is from
    this recipe screams
    executive chef
    have a question?
    # typos
    adapated

    :@subject
    @a
    the
    @pronoun
    each
    mixture
    @number

    :@a
    a
    an

    :@is_are
    is
    are

    :@pronoun
    it
    them
`

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
    @quantity_number_with_possible_range @measurement_clause ?@quantity_unit @lose_ingredient_subcomponent
    @quantity_number_with_possible_range ?@measurement_clause @ingredient_name_with_adjectives @quantity_unit ?@postfix_rep
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
    ?@squeeze @ingredient_name juice
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
    bamboo or metal skewers
    whole wheat yeast dough

    :@ingredient_name
    #.multi
    seasoning blend
    hot sauce
    other citrus type of seasoning
    vanilla extract
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
    maple syrup
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
    cauliflower
    ?@broth_kind_rep broth
    @fillet_type_rep fillet
    chicken
    chicken thigh

    :@greens
    lettuce
    ?mixed salad greens
    @green_sprigs ?sprigs
    sprigs of @green_sprigs
    herb ?sprigs ?@green_sprigs_rep
    @leaves_source ?leaves
    herbs
    romaine hearts
    ?fresh herbs such as @herbs_rep

    :@fillet_type_rep
    .skip @or_and
    +@fillet_type

    :@fillet_type
    salmon

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
    thyme
    rosemary
    oregano
    parsley

    :@leaves_source
    sage
    parsley
    rosemary

    :@vegetables
    avocado
    couscous
    tomato
    cucumber
    rutabaga
    celery root
    turnip
    parsnip
    zucchini

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
    strawberry

    :@dairy
    milk
    cream
    creme fraiche
    sour cream
    ?@cheese_kinds cheese ?crumbles
    @cheese_kind
    yogurt

    :@meat
    back pork ribs
    pork chops
    racks
    turkey
    strip steak

    :@fish
    ahi tuna
    tuna

    :@broth_kind_rep
    .skip @or_and
    +@broth_kind

    :@broth_kind
    chicken
    vegetable

    :@pepper_kind
    cayenne

    :@cheese_kinds
    .skip @or_and
    +@cheese_kind

    :@cheese_kind
    Parmesan
    Romano
    feta
    blue

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
    @adjective @or_and @adjective

    #############################
    :@adjective
    favorite
    homemade
    kosher
    additional
    butter-flavored
    from above
    whole
    trimmed of excess fat

    @size_adjective

    fine
    coarse

    heavy
    mild

    white
    black
    green
    red
    brown
    yellow

    sushi grade
    good-quality

    baby
    confectioners
    flat-leaf
    crispy
    cage-free
    bone-in
    boneless
    skin-on
    low-sodium
    center-cut
    nonfat
    plain
    pure
    ##
    sweet
    semisweet
    ##### 
    @processed_clause
    ##### countries
    french
    italian
    thai
    mexican blend

    :@postfix_rep
    .skip +b @or_and
    +@postfix

    :@postfix
    @postfix_adjective
    @processed_clause

    :@postfix_adjective
    at room temperature
    @postfix_measurement_clause

    :@postfix_measurement_clause
    ?@about @quantity_number_with_possible_range @measurement_unit @postfix_measurement_relationship

    :@postfix_measurement_relationship
    thick
    total
    each

    :@measurement_clause
    ?@about @quantity_number_with_possible_range @weight_unit

    :@processed_clause
    ?@processing_attribute_rep @processed ?@processing_goal

    :@processed
    $sliced
    $ground
    $grated
    $chopped
    $cracked
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
    $roasted
    $quartered
    $halved
    $crumbled
    $warmed
    $warmed to hot
    $frozen
    $thawed
    $washed
    $hulled
    $squeezed
    $cut into pieces

    :@processing_goal
    to remove seeds

    :@processing_attribute_rep
    .skip @or_and
    +@processing_attribute

    :@processing_attribute
    ?@processing_attribute_adj @processing_attribute_name

    :@processing_attribute_name
    $roughly
    $gently
    $freshly
    $fresh
    $thinly

    :@processing_attribute_adj
    very

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
    ?@about @standalone_quantity_number
    ?@about @quantity_number_with_possible_range @quantity_unit ?of
   
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
    bunch
    ears
    teaspoon
    tablespoon
    clove
    caps
    stick
    recipe
    ?@size_adjective heads
    @weight_unit

    :@measurement_unit
    inch
    $cm
    @weight_unit

    :@weight_unit
    pound
    ounce
    kilogram
    $kg

    :@size_adjective
    small
    medium
    medium-large
    large

    :@about
    about
    at least

`
