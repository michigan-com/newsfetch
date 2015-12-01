package recipetypes

import (
	"bytes"

	"github.com/michigan-com/newsfetch/model/bodytypes"
)

//go:generate $GOPATH/bin/stringer -type=Role
type Role int

const (
	Conflict Role = iota

	ShitTail
	EndMarker

	// recipes + others?
	Title

	IngredientsHeading
	DirectionsHeading

	IngredientSubsectionHeading
	Ingredient
	Direction

	Special
	OutOfBandSpecial // special data that's located outside of the recipe itself

	// additional roles for parsing
	PossibleTitle
)

type RoleMatch struct {
	Role Role
	bodytypes.Match
}

func (rm RoleMatch) String() string {
	var buf bytes.Buffer
	buf.WriteString(rm.Confidence.String())
	buf.WriteString(" ")
	buf.WriteString(rm.Role.String())
	return buf.String()
}
