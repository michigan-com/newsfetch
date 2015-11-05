package chartbeat

// func (r *Recent) Run(session *mgo.Session, apiKey string) {
//   chartbeatDebugger.Printf("Recent")

//   urls, err := FormatChartbeatUrls("live/recent/v3", lib.Sites, apiKey)
//   if err != nil {
//     chartbeatDebugger.Println("%v", err)
//     return
//   }

//   recents := f.FetchRecent(urls)

//   if session != nil {
//     chartbeatDebugger.Printf("Saving recents...")

//     f.SaveRecents(recents, session)

//     // Update mapi
//     if !noUpdate {
//       resp, err := http.Get("https://api.michigan.com/recent/")
//       if err != nil {
//         chartbeatDebugger.Printf("%v", err)
//       } else {
//         defer resp.Body.Close()
//         chartbeatDebugger.Printf("Updated recent snapshot at %v", time.Now())
//       }
//     }
//   } else {
//     chartbeatDebugger.Printf("Variable 'mongoUri' not specified, no data will be saved")
//   }
// }