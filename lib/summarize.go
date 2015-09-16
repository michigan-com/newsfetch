package lib

import (
	//"fmt"
	"github.com/neurosnap/sentences/data"
	"github.com/neurosnap/sentences/punkt"
	"github.com/urandom/text-summary/summarize"
)

func LoadTokenizer() *SentenceTokenizer {
	b, err := data.Asset("data/english.json")
	if err != nil {
		panic(err)
	}

	training, _ := punkt.LoadTraining(b)

	tokenizer := &SentenceTokenizer{
		&punkt.DefaultSentenceTokenizer{
			Base:        punkt.NewBase(),
			Punctuation: punkt.Punctuation,
		},
	}

	tokenizer.Storage = training
	tokenizer.SentenceTokenizer = tokenizer
	return tokenizer
}

type PunktTextSplitter struct {
	summarize.DefaultTextSplitter
	*SentenceTokenizer
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
	return punkt.Tokenize(text, p.SentenceTokenizer)
}

func NewPunktSummarizer(title, text string, tokenizer *SentenceTokenizer) summarize.Summarize {
	return summarize.Summarize{
		Title:             title,
		Text:              text,
		Language:          "en",
		StopWordsProvider: summarize.DefaultStopWords{},
		TextSplitter:      PunktTextSplitter{SentenceTokenizer: tokenizer},
		IdealWordCount:    20,
	}
}
