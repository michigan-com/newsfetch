package model

import (
	"time"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Article struct {
	Id          bson.ObjectId  `bson:"_id,omitempty" json:"_id"`
	ArticleId   int            `bson:"article_id" json:"article_id"`
	Headline    string         `bson:"headline" json:"headline`
	Subheadline string         `bson:"subheadline" json:"subheadline"`
	Section     string         `bson:"section" json:"section"`
	Subsection  string         `bson:"subsection" json:"subsection"`
	Source      string         `bson:"source" json:"source"`
	Summary     interface{}    `bson:"summary" json:"summary"`
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

func (article *Article) Save(session *mgo.Session) error {
	// Save the snapshot
	articleCol := session.DB("").C("Article")
	art := Article{}
	err := articleCol.
		Find(bson.M{"article_id": article.ArticleId}).
		Select(bson.M{"_id": 1, "created_at": 1}).
		One(&art)
	if err == nil {
		article.Created_at = art.Created_at
		articleCol.Update(bson.M{"_id": art.Id}, article)
		// Debugger.Println("Article updated: ", article)
	} else {
		//bulk.Insert(article)
		articleCol.Insert(article)
		// Debugger.Println("Article added: ", article)
	}

	return nil
}