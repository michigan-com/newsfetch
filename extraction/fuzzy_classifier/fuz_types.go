package fuzzy_classifier

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
	TagsByPos  []map[string]Range
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

	r.TagsByPos[pos][tag] = rang

	ranges := append(r.TagsByName[tag], rang)
	r.TagsByName[tag] = ranges
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
			if rangeMap[tag] == rang {
				delete(rangeMap, tag)
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
