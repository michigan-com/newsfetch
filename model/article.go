package model

import (
	"fmt"
	"time"

	"github.com/michigan-com/newsfetch/lib"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var Debugger = lib.NewCondLogger("article-model")

var articleIdIndex = mgo.Index{
	Key:    []string{"article_id"},
	Unique: true,
}

type Article struct {
	Id          bson.ObjectId  `bson:"_id,omitempty" json:"_id"`
	ArticleId   int            `bson:"article_id" json:"article_id"`
	Headline    string         `bson:"headline" json:"headline`
	Subheadline string         `bson:"subheadline" json:"subheadline"`
	Section     string         `bson:"section" json:"section"`
	Subsection  string         `bson:"subsection" json:"subsection"`
	Source      string         `bson:"source" json:"source"`
	Created_at  time.Time      `bson:"created_at" json:"created_at"`
	Updated_at  time.Time      `bson:"updated_at" json:"updated_at"`
	Timestamp   time.Time      `bson:"timestamp" json:"timestamp"`
	Url         string         `bson:"url" json:"url"`
	Photo       *Photo         `bson:"photo" json:"photo"`
	BodyText    string         `bson:"body" json:"body"`
	Visits      []TimeInterval `body:"visits" json:"visits"`
}

type PhotoInfo struct {
	Url    string `bson:"url"`
	Width  int    `bson:"width"`
	Height int    `bson:"height"`
}

type Photo struct {
	Caption   string    `bson:"caption"`
	Credit    string    `bson:"credit"`
	Full      PhotoInfo `bson:"full"`
	Thumbnail PhotoInfo `bson:"thumbnail"`
}

type TimeInterval struct {
	Max       int       `bson:"max"`
	Timestamp time.Time `bson:"timestamp"`
}

func (a *Article) String() string {
	return fmt.Sprintf("<Article Id: %d, Headline: %s, Url: %s>", a.ArticleId, a.Headline, a.Url)
}

func (article *Article) Save(session *mgo.Session) (bool, error) {
	isNew := false
	articleCol := session.DB("").C("Article")
	err := articleCol.EnsureIndex(articleIdIndex)
	if err != nil {
		lib.Logger.Println("Article ensure article_id is unique failed: ", err)
		return isNew, err
	}

	update := bson.M{
		"$set": bson.M{
			"headline":    article.Headline,
			"subheadline": article.Subheadline,
			"section":     article.Section,
			"subsection":  article.Subsection,
			"source":      article.Source,
			"updated_at":  article.Updated_at,
			"timestamp":   article.Timestamp,
			"url":         article.Url,
			"photo":       article.Photo,
			"body":        article.BodyText,
		},
		"$setOnInsert": bson.M{"created_at": article.Created_at},
	}

	info, err := articleCol.Upsert(bson.M{"article_id": article.ArticleId}, update)
	if err != nil {
		return isNew, err
	}

	if info.UpsertedId != nil {
		isNew = true
	}

	return isNew, nil
}
