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

type Messages struct {
	Errors   []string
	Warnings []string
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
