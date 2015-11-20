package model

import (
	"fmt"
)

type RecipeExtractionResult struct {
	Recipes             []*Recipe
	UnusedParagraphs    []RecipeFragment
	EmbeddedArticleUrls []string
}

type ExtractedBody struct {
	Text       string
	RecipeData RecipeExtractionResult
	Messages   *Messages
}

// TODO only use the Sections array for this, don't have section/subsection distinction
type ExtractedSection struct {
	Section    string
	Subsection string
	Sections   []string
}

func (e *ExtractedBody) String() string {
	return fmt.Sprintf("<ExtractedBody %s\n %s>\n", e.Text, e.Messages)
}

type Messages struct {
	Errors   []string
	Warnings []string
}

func (m *Messages) String() string {
	return fmt.Sprintf("<Messages Errors: %v\n Warnings: %v>\n", m.Errors, m.Warnings)
}

func (msg *Messages) AddMessages(context string, other *Messages) {
	var prefix string
	if context == "" {
		prefix = ""
	} else {
		prefix = context + ": "
	}

	for _, message := range other.Errors {
		msg.Errors = append(msg.Errors, prefix+message)
	}
	for _, message := range other.Warnings {
		msg.Warnings = append(msg.Warnings, prefix+message)
	}
}

func (msg *Messages) AddError(message string) {
	msg.Errors = append(msg.Errors, message)
}

func (msg *Messages) AddWarning(message string) {
	msg.Warnings = append(msg.Warnings, message)
}

func (msg *Messages) AddErrorf(format string, args ...interface{}) {
	msg.AddError(fmt.Sprintf(format, args...))
}

func (msg *Messages) AddWarningf(format string, args ...interface{}) {
	msg.AddWarning(fmt.Sprintf(format, args...))
}
