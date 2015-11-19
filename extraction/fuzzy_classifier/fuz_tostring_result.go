package fuzzy_classifier

import (
	"bytes"
	// "strconv"
)

func (r *Result) Description() string {
	buf := new(bytes.Buffer)

	for pos, word := range r.Words {
		taggingsByTag := r.TagsByPos[pos]

		var singleWordTags []string

		for i := len(r.TagDefs) - 1; i >= 0; i-- {
			tag := r.TagDefs[i].Tag

			for _, tagging := range taggingsByTag[tag] {
				if tagging.Len == 1 {
					singleWordTags = append(singleWordTags, tag)
					continue
				}
				buf.WriteString(tag)

				end := tagging.Pos + tagging.Len
				for p2 := tagging.Pos; p2 < end; p2++ {
					buf.WriteString(" ")
					buf.WriteString(r.Words[p2].Raw)
				}
				buf.WriteString("\n")
			}
		}

		for _, tag := range singleWordTags {
			buf.WriteString(tag)
			buf.WriteString(" ")
		}
		buf.WriteString(word.Raw)
		buf.WriteString("\n")
	}

	return buf.String()
}

// func (r *Result) descriptionInto(buf *bytes.Buffer, indentation string) {
// }
