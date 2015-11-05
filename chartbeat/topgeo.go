package chartbeat

// func (t *TopGeo) Run(session *mgo.Session, apiKey string) {
//   chartbeatDebugger.Printf("Topgeo")

//   urls, err := FormatChartbeatUrls("live/top_geo/v1", lib.Sites, apiKey)
//   if err != nil {
//     chartbeatDebugger.Println("ERROR: %v", err)
//     return
//   }

//   topGeo := f.FetchTopGeo(urls)

//   if session != nil {
//     chartbeatDebugger.Printf("Saving topgeo...")

//     f.SaveTopGeo(topGeo, session)

//     // Update mapi
//     if !noUpdate {
//       resp, err := http.Get("https://api.michigan.com/topgeo/")
//       if err != nil {
//         chartbeatDebugger.Printf("%v", err)
//       } else {
//         defer resp.Body.Close()
//         chartbeatDebugger.Printf("Updated topgeo snapshot at Mapi at %v", time.Now())
//       }
//     }
//   } else {
//     chartbeatDebugger.Printf("Variable 'mongoUri' not specified, no data will be saved")
//     chartbeatDebugger.Printf("%v", topGeo)
//   }
// }