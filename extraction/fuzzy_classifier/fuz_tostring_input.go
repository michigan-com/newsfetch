package fuzzy_classifier

import (
	"bytes"
)

func (c *Classifier) Description() string {
	var buf bytes.Buffer
	first := true
	for i := len(c.categories) - 1; i >= 0; i-- {
		category := c.categories[i]

		if category.builtIn {
			continue
		}
		if first {
			first = false
		} else {
			buf.WriteString("\n")
		}
		buf.WriteString("CAT ")
		buf.WriteString(category.tag)
		buf.WriteString("\n")
		for _, attr := range category.attributes {
			buf.WriteString("ATTR ")
			buf.WriteString(attr.name)
			buf.WriteString(" = ")
			buf.WriteString(attr.value)
			buf.WriteString("\n")
		}
		if len(category.skippableTags) != 0 {
			buf.WriteString("SKIP")
			for _, tag := range category.skippableTags {
				buf.WriteString(" ")
				buf.WriteString(tag)
			}
			buf.WriteString("\n")
		}
		for _, scheme := range category.schemes {
			buf.WriteString("SCHEME")
			for _, req := range scheme.requirements {
				buf.WriteString(" ")
				buf.WriteString(req.String())
			}
			buf.WriteString("\n")
		}
	}
	return buf.String()
}

func (r Requirement) String() string {
	prefix := ""
	if r.optional {
		prefix = "?"
	}
	if r.repeating {
		prefix = prefix + "+"
	}

	if len(r.tag) > 0 {
		return prefix + "<" + r.tag + ">"
	} else {
		return prefix + r.literal
	}
}
