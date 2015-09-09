package lib

import (
	//"fmt"
	"github.com/neurosnap/sentences/data"
	"github.com/neurosnap/sentences/punkt"
	"github.com/urandom/text-summary/summarize"
)

type PunktTextSplitter struct {
	summarize.DefaultTextSplitter
}

type SentenceTokenizer struct {
	*punkt.DefaultSentenceTokenizer
}

func (s *SentenceTokenizer) AnnotateTokens(tokens []*punkt.DefaultToken) []*punkt.DefaultToken {
	tokens = s.AnnotateFirstPass(tokens)
	tokens = s.AnnotateSecondPass(tokens)
	//fmt.Println("HI I ACTUALLY GOT HIT\n------------")
	return tokens
}

func (p PunktTextSplitter) Sentences(text string) []string {
	b, err := data.Asset("data/english.json")
	if err != nil {
		panic(err)
	}

	training, err := punkt.LoadTraining(b)

	tokenizer := &SentenceTokenizer{
		&punkt.DefaultSentenceTokenizer{
			Base:        punkt.NewBase(),
			Punctuation: punkt.Punctuation,
		},
	}

	tokenizer.Storage = training
	tokenizer.SentenceTokenizer = tokenizer

	return punkt.Tokenize(text, tokenizer)
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
