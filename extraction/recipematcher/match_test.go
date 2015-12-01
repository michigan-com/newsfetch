package recipematcher

import (
	"testing"

	"github.com/michigan-com/newsfetch/model/bodytypes"
	"github.com/michigan-com/newsfetch/model/recipetypes"
	"github.com/michigan-com/newsfetch/util/messages"
)

func TestMatch(t *testing.T) {
	// nothings
	if a, r, e := o("hello world", "distantly Ingredient"); a != e {
		t.Errorf("Expected %#v, got %#v, rationale: %v", e, a, r)
	}
	if a, r, e := o("wumblas", "distantly Ingredient"); a != e {
		t.Errorf("Expected %#v, got %#v, rationale: %v", e, a, r)
	}
	if a, r, e := o("A very long paragraph that does not seem like an ingredient at all, because it's so very long.", "distantly Direction"); a != e {
		t.Errorf("Expected %#v, got %#v, rationale: %v", e, a, r)
	}
	if a, r, e := o("Several sentences. Probably not an ingredient.", "distantly Direction"); a != e {
		t.Errorf("Expected %#v, got %#v, rationale: %v", e, a, r)
	}

	// ingredients
	if a, r, e := o("3 apples", "perfect Ingredient"); a != e {
		t.Errorf("Expected %#v, got %#v, rationale: %v", e, a, r)
	}
	if a, r, e := o("roasted garlic", "perfect Ingredient"); a != e {
		t.Errorf("Expected %#v, got %#v, rationale: %v", e, a, r)
	}
	if a, r, e := o("3 wumblas", "likely Ingredient"); a != e {
		t.Errorf("Expected %#v, got %#v, rationale: %v", e, a, r)
	}
	if a, r, e := o("bla bla bla apple bla bla bla", "possible Ingredient"); a != e {
		t.Errorf("Expected %#v, got %#v, rationale: %v", e, a, r)
	}

	// directions
	if a, r, e := o("add some wumblas", "likely Direction"); a != e {
		t.Errorf("Expected %#v, got %#v, rationale: %v", e, a, r)
	}
	if a, r, e := o("cut it into pieces", "likely Direction"); a != e {
		t.Errorf("Expected %#v, got %#v, rationale: %v", e, a, r)
	}

	// conflicts
	if a, r, e := o("3 simulatedconflictingword apples", "perfect Conflict"); a != e {
		t.Errorf("Expected %#v, got %#v, rationale: %v", e, a, r)
	}
	if a, r, e := o("3 simulatedconflictingword wumblas", "likely Conflict"); a != e {
		t.Errorf("Expected %#v, got %#v, rationale: %v", e, a, r)
	}

	// others
	if a, r, e := o("SIRUP", "likely IngredientSubsectionHeading"); a != e {
		t.Errorf("Expected %#v, got %#v, rationale: %v", e, a, r)
	}
}

func o(input string, expected string) (string, string, string) {
	input = CanonicalString(input)
	r, m, _ := Match(bodytypes.Paragraph{input, nil}, new(messages.Messages))
	rm := recipetypes.RoleMatch{r, m}
	return rm.String(), rm.Rationale, expected
}
