package recipematcher

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

    # need some kind of assertion (like “starts a sentence” or “preceded by an adverb” or “after a comma”)
    :@confusing_action
    $mix
    $warm
    $top @subject
    $cream @subject
    gently $spoon
    $spoon @subject
    #$cut
    $simulatedconflictingword

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
