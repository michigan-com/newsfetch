package commands

import (
	"time"
  "sort"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/spf13/cobra"

  "github.com/michigan-com/newsfetch/lib"
  m "github.com/michigan-com/newsfetch/model/chartbeat"
  h "github.com/michigan-com/newsfetch/model/historical"
)

var histDebugger = lib.NewCondLogger("newsfetch:commands:historical")

type ByVisits []h.Referrer

func (a ByVisits) Len() int { return len(a) }
func (a ByVisits) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByVisits) Less(i, j int) bool { return a[i].TotalViewers > a[j].TotalViewers }

var MIN_VISITS = 5.0

var cmdHistorical = &cobra.Command{
	Use:   "historical",
	Short: "Compile a history data to keep track of", Run: func(cmd *cobra.Command, args []string) {
		startTime = time.Now()
    referrerSnapshot := m.ReferrersSnapshot{}
    var newRecord = h.HistoricalEntry{}
    var referrerMap = make(map[string]h.Referrer)
  	var session *mgo.Session
  	if globalConfig.MongoUrl != "" {
  		session = lib.DBConnect(globalConfig.MongoUrl)
  		defer lib.DBClose(session)
  	}
    referrersCollection := session.DB("").C("Referrers")
    err := referrersCollection.Find(bson.M{}).One(&referrerSnapshot)

    if err != nil {
      histDebugger.Printf("%v", err)
      return
    }

    for _, sourceRef := range referrerSnapshot.Referrers {
      source := sourceRef.Source
			histDebugger.Printf("snapshot is %s", referrerSnapshot)

      for site, value := range sourceRef.Referrers {
        numValue := value.(float64)

				histDebugger.Printf("%s->%s: %d", source, site, numValue)

        if _, ok := referrerMap[site]; !ok {
          pubViews := make ([]bson.M, 0,  len(referrerSnapshot.Referrers))
          referrerMap[site] = h.Referrer{
            Site: site,
            TotalViewers: numValue,
            PublicationsCount: append(pubViews, bson.M{
              "source": source,
              "viewers": numValue,
            }),
          }
        } else {
          _ref := referrerMap[site]
          _ref.TotalViewers += numValue
          _ref.PublicationsCount = append(_ref.PublicationsCount, bson.M{
            "source": source,
            "viewers": numValue,
          })

					referrerMap[site] = _ref
        }
      }
    }

    referrers := make([]h.Referrer, 0, len(referrerMap))
    for _, ref := range referrerMap {
      if ref.TotalViewers >= MIN_VISITS {
        referrers = append(referrers, ref)
      }
    }
    sort.Sort(ByVisits(referrers))

    newRecord.Timestamp = time.Now()
    newRecord.Referrers = referrers

    saveCollection := session.DB("").C("History")

    saveCollection.Insert(newRecord)

    if err != nil {
      return
    }
		getElapsedTime(&startTime)
	},
}
