package orderedlist

import (
	"reflect"
	"testing"
)

func TestInsertString(t *testing.T) {
	var ordered []string

	ordered = InsertString(ordered, "foo", true)
	if ok := reflect.DeepEqual(ordered, []string{"foo"}); !ok {
		t.Errorf("Insertion into an empty list didn't work")
	}

	ordered = InsertString(ordered, "bar", true)
	if ok := reflect.DeepEqual(ordered, []string{"bar", "foo"}); !ok {
		t.Errorf("Insertion into a 1-element list didn't work, got %v", ordered)
	}

	ordered = InsertString(ordered, "boz", true)
	if ok := reflect.DeepEqual(ordered, []string{"bar", "boz", "foo"}); !ok {
		t.Errorf("Insertion into a 2-element list didn't work, got %v", ordered)
	}

	ordered = InsertString(ordered, "bar", true)
	if ok := reflect.DeepEqual(ordered, []string{"bar", "boz", "foo"}); !ok {
		t.Errorf("Inserting a duplicate should be a no-op")
	}
}

func TestContainsString(t *testing.T) {
	var ordered []string

	if found := ContainsString(ordered, "foo"); found != false {
		t.Errorf("Empty ordered shouldn't contain anything")
	}

	ordered = []string{"foo"}

	if found := ContainsString(ordered, "foo"); found != true {
		t.Errorf("1-element list should contain its element")
	}
	if found := ContainsString(ordered, "bar"); found != false {
		t.Errorf("1-element list shouldn't contain bogus elements")
	}

	ordered = []string{"bar", "foo"}

	if found := ContainsString(ordered, "foo"); found != true {
		t.Errorf("2-element list should contain its elements")
	}
	if found := ContainsString(ordered, "bar"); found != true {
		t.Errorf("2-element list should contain its elements")
	}
	if found := ContainsString(ordered, "boz"); found != false {
		t.Errorf("2-element list shouldn't contain bogus elements")
	}
}
