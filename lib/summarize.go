package lib

import (
	"github.com/neurosnap/sentences/data"
	"github.com/neurosnap/sentences/punkt"
	"github.com/urandom/text-summary/summarize"
)

type PunktTextSplitter struct {
	summarize.DefaultTextSplitter
}

func (p PunktTextSplitter) Sentences(text string) []string {
	b, err := data.Asset("data/english.json")
	if err != nil {
		panic(err)
	}
	training, err := punkt.LoadTraining(b)

	tokenizer := punkt.NewSentenceTokenizer(training)

	return tokenizer.Tokenize(text)
}

func NewPunktSummarizer(title, text string) summarize.Summarize {
	return summarize.Summarize{
		Title:             title,
		Text:              text,
		Language:          "en",
		StopWordsProvider: summarize.DefaultStopWords{},
		TextSplitter:      PunktTextSplitter{},
		IdealWordCount:    20,
	}
}
