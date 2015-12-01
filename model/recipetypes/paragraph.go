package recipetypes

import (
	"github.com/michigan-com/newsfetch/model/bodytypes"
)

type Paragraph struct {
	bodytypes.Paragraph

	Role Role
	bodytypes.Match

	Fragments []Fragment
}
