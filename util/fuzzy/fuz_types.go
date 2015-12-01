package fuzzy

import (
	"strings"
)

type Classifier struct {
	categories       []*Category
	multiVariantTags map[string]bool

	TagDefs       []*TagDef
	TagDefsByName map[string]*TagDef
}

type TagDef struct {
	Tag   string
	Index int
}

type Category struct {
	tag           string
	tagDef        *TagDef
	schemes       []Scheme
	attributes    []Attribute
	skippableTags []string
	builtIn       bool

	skipBefore bool
	skipAfter  bool

	leadingWordSet map[string]bool
	minLen         int
	maxLen         int
}

type Attribute struct {
	name  string
	value string
}

type Scheme struct {
	requirements []Requirement

	lineNr int
}

type ReqType int

const (
	ReqLiteral ReqType = iota
	ReqStem
	ReqTag
)

type Requirement struct {
	typ       ReqType
	literal   string
	stem      string
	tag       string
	optional  bool
	repeating bool
}

type Result struct {
	Words      []Word
	TagsByPos  []map[string][]Range
	TagsByName map[string][]Range

	TagDefs       []*TagDef
	TagDefsByName map[string]*TagDef

	multiVariantTags map[string]bool
}

type Word struct {
	Raw        string
	Trimmed    string
	Normalized string
	Stem       string
}

type Range struct {
	Pos int
	Len int
}

type WordFormat int

const (
	Raw WordFormat = iota
	Trimmed
)

func (r *Result) RangeCoversEntireInput(rang Range) bool {
	return rang.Pos == 0 && rang.Len == len(r.Words)
}

func (r *Result) GetWordString(i int, format WordFormat) string {
	switch format {
	case Raw:
		return r.Words[i].Raw
	case Trimmed:
		return r.Words[i].Trimmed
	default:
		panic("unknown format")
	}
}

func (r *Result) GetRangeString(rang Range, format WordFormat) string {
	var result []string
	for i := rang.Pos; i < rang.Pos+rang.Len; i++ {
		result = append(result, r.GetWordString(i, format))
	}
	return strings.Join(result, " ")
}

func (r *Result) GetTagMatchString(tag string, format WordFormat) (string, bool) {
	for _, rang := range r.TagsByName[tag] {
		return r.GetRangeString(rang, format), true
	}
	return "", false
}

func (r *Result) GetTagMatch(tag string) (Range, bool) {
	for _, rang := range r.TagsByName[tag] {
		return rang, true
	}
	return Range{Pos: -1, Len: 0}, false
}

func (r *Result) GetAllTagMatchStrings(tag string, format WordFormat) []string {
	result := make([]string, 0, len(r.TagsByName[tag]))
	for _, rang := range r.TagsByName[tag] {
		result = append(result, r.GetRangeString(rang, format))
	}
	return result
}

func (r *Result) HasTag(tag string) bool {
	return len(r.TagsByName[tag]) > 0
}

func (r *Result) HasTagAt(tag string, pos int) bool {
	if pos >= len(r.TagsByPos) {
		return false
	}
	return len(r.TagsByPos[pos][tag]) > 0
}

func (r *Result) AddTag(tag string, pos int, length int) {
	if !strings.HasPrefix(tag, "@") {
		panic("Tags must start with @")
	}

	if r.TagDefsByName[tag] == nil {
		println("All", len(r.TagDefsByName), "tags:")
		for tag, _ := range r.TagDefsByName {
			println(tag)
		}
		panic("Attemp to add unknown tag " + tag)
	}

	rang := Range{Pos: pos, Len: length}

	if !r.multiVariantTags[tag] {
		if r.IsTagCoveringRange(tag, rang) {
			return
		}

		r.RemoveTagInstancesCoveredByRange(tag, rang)
	}

	rangeMap := r.TagsByPos[pos]
	if findRangeInList(rang, rangeMap[tag]) < 0 {
		rangeMap[tag] = append(rangeMap[tag], rang)
	}
	if findRangeInList(rang, r.TagsByName[tag]) < 0 {
		r.TagsByName[tag] = append(r.TagsByName[tag], rang)
	}
}

func (r *Result) RemoveTagInstancesCoveredByRange(tag string, limit Range) {
	found := false
	ranges := r.TagsByName[tag]
	for _, rang := range ranges {
		if rang.Pos >= limit.Pos && (rang.Pos+rang.Len) <= (limit.Pos+limit.Len) {
			found = true
			break
		}
	}

	if !found {
		return
	}

	remaining := make([]Range, 0, len(ranges))
	for _, rang := range ranges {
		if rang.Pos >= limit.Pos && (rang.Pos+rang.Len) <= (limit.Pos+limit.Len) {
			// delete from r.TagsByPos too
			rangeMap := r.TagsByPos[rang.Pos]
			if list := rangeMap[tag]; list != nil {
				if newList, found := removeRangeFromList(rang, list); found {
					if len(newList) == 0 {
						delete(rangeMap, tag)
					} else {
						rangeMap[tag] = newList
					}
				}
			}
		} else {
			// keep
			remaining = append(remaining, rang)
		}
	}
	r.TagsByName[tag] = remaining
}

func (r *Result) IsTagCoveringRange(tag string, covered Range) bool {
	for _, rang := range r.TagsByName[tag] {
		if rang.Pos <= covered.Pos && (rang.Pos+rang.Len) >= (covered.Pos+covered.Len) {
			return true
		}
	}
	return false
}

func findRangeInList(rang Range, ranges []Range) int {
	for idx, r := range ranges {
		if r.Pos == rang.Pos && r.Len == rang.Len {
			return idx
		}
	}
	return -1
}

func removeRangeFromList(rang Range, ranges []Range) ([]Range, bool) {
	idx := findRangeInList(rang, ranges)
	if idx < 0 {
		return ranges, false
	}

	result := make([]Range, 0, len(ranges))
	for i, r := range ranges {
		if i != idx {
			result = append(result, r)
		}
	}
	return result, true
}
