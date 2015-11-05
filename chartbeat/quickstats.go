package chartbeat



// func (q *QuickStats) Run(session *mgo.Session, apiKey string) {
//   chartbeatDebugger.Printf("Quickstats")

//   quickStats := f.FetchQuickStats(urls)

//   if session != nil {
//     chartbeatDebugger.Printf("Saving quickstats...")

//     f.SaveQuickStats(quickStats, session)

//     // Update mapi
//     if !noUpdate {
//       resp, err := http.Get("https://api.michigan.com/quickstats/")
//       if err != nil {
//         chartbeatDebugger.Printf("%v", err)
//       } else {
//         defer resp.Body.Close()
//         chartbeatDebugger.Printf("Updated quickstats snapshot at Mapi at %v", time.Now())
//       }
//     }
//   } else {
//     chartbeatDebugger.Printf("Variable 'mongoUri' not specified, no data will be saved")
//     chartbeatDebugger.Printf("%v", quickStats)
//   }
// }
