package fuzzy_classifier

import (
	"errors"
	"fmt"
	"strings"
)

var builtInTags = []string{"@s", "@cap", "@twitter", "@url", "@email", "@integer", "@float", "@fraction", "@currency-number"}

func NewFuzzyClassifierFromString(definition string) *Classifier {
	classifier := NewFuzzyClassifier()
	classifier.AddOrPanic(definition)
	return classifier
}

func NewFuzzyClassifier() *Classifier {
	classifier := new(Classifier)
	classifier.multiVariantTags = make(map[string]bool)
	classifier.TagDefsByName = make(map[string]*TagDef)

	for _, tag := range builtInTags {
		classifier.findOrCreateTag(tag, true)
	}

	err := classifier.add(builtInRules, true)
	if err != nil {
		panic(err)
	}
	return classifier
}

func (classifier *Classifier) Add(definition string) error {
	return classifier.add(definition, false)
}

func (classifier *Classifier) AddOrPanic(definition string) {
	err := classifier.add(definition, false)
	if err != nil {
		panic(err)
	}
}

func (classifier *Classifier) add(definition string, builtIn bool) error {
	var currentCategory *Category
	var newCategories []*Category

	lines := strings.Split(definition, "\n")
	for lno, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}

		words := strings.Fields(line)

		if strings.HasPrefix(line, ":") {
			if !strings.HasPrefix(line, ":@") {
				return errors.New(fmt.Sprintf("line %d: expected @ after : in %#v", lno+1, line))
			}

			tag := strings.Replace(words[0], ":", "", 1)
			if len(tag) == 1 {
				return errors.New(fmt.Sprintf("line %d: expected a tag name after :@ in %#v", lno+1, line))
			}

			words = words[1:]

			currentCategory = new(Category)
			currentCategory.builtIn = builtIn
			currentCategory.tag = tag
			currentCategory.skipBefore = false
			currentCategory.skipAfter = false
			newCategories = append(newCategories, currentCategory)
			continue
		}

		if currentCategory == nil {
			return errors.New(fmt.Sprintf("line %d: expected a category header starting with :@, got %#v", lno+1, line))
		}

		if strings.HasPrefix(line, ".") {
			switch words[0] {
			case ".skip":
				for _, tag := range words[1:] {
					if tag == "-b" {
						currentCategory.skipBefore = false
						continue
					} else if tag == "+b" {
						currentCategory.skipBefore = true
						continue
					} else if tag == "-a" {
						currentCategory.skipAfter = false
						continue
					} else if tag == "+a" {
						currentCategory.skipAfter = true
						continue
					}
					if !strings.HasPrefix(tag, "@") {
						return errors.New(fmt.Sprintf("line %d: expected @ at the start of tag %#v in %#v", lno+1, tag, line))
					}
					currentCategory.skippableTags = append(currentCategory.skippableTags, tag)
				}

			case ".multi":
				classifier.multiVariantTags[currentCategory.tag] = true
			default:
				return errors.New(fmt.Sprintf("line %d: unknown instruction %#v in %#v", lno+1, words[0], line))
			}
			continue
		}

		reqs := wordsToRequirements(words)
		scheme := Scheme{requirements: reqs, lineNr: lno + 1}
		currentCategory.schemes = append(currentCategory.schemes, scheme)
	}

	for i := len(newCategories) - 1; i >= 0; i-- {
		category := newCategories[i]

		for _, scheme := range category.schemes {
			for _, req := range scheme.requirements {
				if err := verifyRequirement(classifier, category, scheme, req); err != nil {
					return err
				}
			}
		}

		category.tagDef = classifier.findOrCreateTag(category.tag, true)

		classifier.categories = append(classifier.categories, category)
	}

	return nil
}

func (classifier *Classifier) findOrCreateTag(tag string, create bool) *TagDef {
	def := classifier.TagDefsByName[tag]
	if (def == nil) && create {
		def = &TagDef{Tag: tag, Index: len(classifier.TagDefs)}
		classifier.TagDefs = append(classifier.TagDefs, def)
		classifier.TagDefsByName[tag] = def
	}
	return def
}

func wordsToRequirements(words []string) []Requirement {
	result := make([]Requirement, 0, len(words))
	for _, word := range words {
		result = append(result, wordToRequirement(word))
	}
	return result
}

func wordToRequirement(word string) Requirement {
	optional := false
	repeating := false
	if strings.HasPrefix(word, "?") {
		optional = true
		word = strings.Replace(word, "?", "", 1)
	}
	if strings.HasPrefix(word, "+") {
		repeating = true
		word = strings.Replace(word, "+", "", 1)
	}

	if strings.HasPrefix(word, "@") {
		return Requirement{typ: ReqTag, tag: word, optional: optional, repeating: repeating}
	} else if strings.HasPrefix(word, "$") {
		word = strings.Replace(word, "$", "", 1)
		return Requirement{typ: ReqLiteral, literal: word, stem: Normalize(word), optional: optional, repeating: repeating}
	} else {
		return Requirement{typ: ReqStem, literal: word, stem: Stem(word), optional: optional, repeating: repeating}
	}
}

func verifyRequirement(classifier *Classifier, category *Category, scheme Scheme, req Requirement) error {
	switch req.typ {
	case ReqLiteral:
		return nil
	case ReqStem:
		return nil

	case ReqTag:
		if classifier.findOrCreateTag(req.tag, false) == nil {
			return errors.New(fmt.Sprintf("line %d: unknown tag %#v in a rule for %#v", scheme.lineNr, req.tag, category.tag))
		}
		return nil

	default:
		panic("Unknown requirement type")
	}
}
