// Code generated by "stringer -type=Role"; DO NOT EDIT

package recipetypes

import "fmt"

const _Role_name = "ConflictShitTailEndMarkerTitleIngredientsHeadingDirectionsHeadingIngredientSubsectionHeadingIngredientDirectionSpecialOutOfBandSpecialPossibleTitle"

var _Role_index = [...]uint8{0, 8, 16, 25, 30, 48, 65, 92, 102, 111, 118, 134, 147}

func (i Role) String() string {
	if i < 0 || i >= Role(len(_Role_index)-1) {
		return fmt.Sprintf("Role(%d)", i)
	}
	return _Role_name[_Role_index[i]:_Role_index[i+1]]
}
