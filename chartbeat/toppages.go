package chartbeat

// func (t *TopPages) Run(session *mgo.Session, apiKey string) {
//   chartbeatDebugger.Println("Fetching toppages")
//   urls, err := FormatChartbeatUrls("live/toppages/v3", lib.Sites, apiKey)
//   urls = f.AddUrlParams(urls, "loyalty=1")

//   if err != nil {
//     chartbeatDebugger.Printf("ERROR: %v", err)
//     return
//   }

//   snapshot := f.FetchTopPages(urls)

//   if session != nil {
//     chartbeatDebugger.Println("Saving toppages snapshot")
//     err := f.SaveTopPagesSnapshot(snapshot, session)
//     if err != nil {
//       chartbeatDebugger.Printf("ERROR: %v", err)
//       return
//     }

//     f.CalculateTimeInterval(snapshot, session)

//     // Update mapi to let it know that a new snapshot has been saved
//     if !noUpdate {
//       resp, err := http.Get("https://api.michigan.com/popular/")
//       if err != nil {
//         chartbeatDebugger.Printf("%v", err)
//       } else {
//         defer resp.Body.Close()
//         now := time.Now()
//         chartbeatDebugger.Printf("Updated toppages snapshot at Mapi at %v", now)
//       }
//     }
//   } else {
//     chartbeatDebugger.Printf("Variable 'mongoUri' not specified, no data will be saved")
//   }
// }