package chartbeat

// func (r *Referrers) Run(session *mgo.Session, apiKey string) {
//   chartbeatDebugger.Printf("Referrers")

//   urls, err := FormatChartbeatUrls("live/referrers/v3", lib.Sites, apiKey)
//   if err != nil {
//     chartbeatDebugger.Println("%v", err)
//     return
//   }

//   referrers := f.FetchReferrers(urls)

//   if session != nil {
//     chartbeatDebugger.Printf("Saving referrers...")

//     f.SaveReferrers(referrers, session)

//     // Update mapi
//     if !noUpdate {
//       resp, err := http.Get("https://api.michigan.com/referrers/")
//       if err != nil {
//         chartbeatDebugger.Printf("%v", err)
//       } else {
//         defer resp.Body.Close()
//         chartbeatDebugger.Printf("Updated referrers snapshot at Mapi at %v", time.Now())
//       }
//     }
//   } else {
//     chartbeatDebugger.Printf("Variable 'mongoUri' not specified, no data will be saved")
//   }
// }