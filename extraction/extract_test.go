package extraction_test

import (
	"strings"
	"testing"

	"github.com/michigan-com/newsfetch/extraction"
	"github.com/michigan-com/newsfetch/lib"
	m "github.com/michigan-com/newsfetch/model"
)

var globalForbidden = []string{
	`Follow him on Twitter`,
	`Follow her on Twitter `,
	`Follow on Twitter`,
	`Listen to him at`,
	`Listen to her at`,
	`Check out our latest Tigers podcast`,
	`download our free Tigers Xtra app`,
	`Copyright 2015 The Associated Press.`,
	`All rights reserved`,
	`This material may not be published`,
}

type TestRec struct {
	url       string
	expected  []string
	forbidden []string
	html      string
}

func TestContactLineRemovalWithoutSpecialStrings(t *testing.T) {
	runTest(t, &TestRec{
		url:       "",
		expected:  []string{`Foo foo`, `Bar bar`},
		forbidden: []string{`Contact Robert Allen`},
		html: `
			<p>Foo foo.</p>
			<p>Bar bar.</p>
			<p><em>Contact Robert Allen at rallen@freepress.com or <a href="http://www.twitter.com/rallenMI">@rallenMI</a>. </em></p>`,
	})
}

func TestIntegrationOakland(t *testing.T) {
	runTest(t, &TestRec{
		url: "http://www.freep.com/story/news/local/michigan/oakland/2015/08/20/police-chase-troy-bloomfield-hills-warren-absconder-shooting/32056645/",
		expected: []string{`The Oakland County Sheriff's office is reviewing the fatal shooting of a parolee by police late Wednesday night following a high-speed chase that started in Macomb County and ended with a crash on Telegraph Road in Bloomfield Hills.`,
			`"We've got a lot of work to do," McCabe said.`},
		forbidden: []string{},
		html: `
			<p>The Oakland County Sheriff's office is reviewing the fatal shooting of a parolee by police late Wednesday night following a high-speed chase that started in Macomb County and ended with a crash on Telegraph Road in Bloomfield Hills.</p>
			<p><p>Investigators say police opened fire on Deviere Ernel Ransom, 24, of Detroit, after he ran behind a nearby drug store and then pulled a handgun on approaching officers. Police say he had fired shot at a Troy police officer during the car chase, which started in Warren.</p></p>
			<p><p>According to investigators, Ransom had sawed off his tether earlier Wednesday evening, and armed with a gun, made threats to kill someone. Warren police were notified and began chasing the suspect in a car just before midnight near 12 Mile and Mound roads.</p></p>
			<p>Police chased the vehicle into Troy, where Ransom reportedly fired shots at a Troy police officer who was joining the pursuit. Bloomfield Police officers were alerted  as the vehicle left the freeway and collided with another car at the corner of Woodward Avenue and Square Lake Road. Ransom ran from the car to behind a CVS drug store.</p>
			<p>As police approached, Ransom  pulled a handgun, prompting officers to open fire. Investigators believe both Bloomfield Hills and Warren officers fired their weapons.</p>
			<p>Those departments asked the Oakland County Sheriff's Office, the largest law enforcement organization in the county, to conduct an investigation, a common practice when police officers are involved in shootings.</p>
			<p>"Our investigators were out there all night," said Oakland County Undersheriff Michael McCabe.</p>
			<p>In addition to dashcam video from patrol cars, the investigators were interviewing witnesses, including the other motorist, who suffered minor injuries. Investigators were also checking with businesses in the area to see if they had surveillance cameras operating at the time.</p>
			<p>"We've got a lot of work to do," McCabe said.</p>
			<p>Ransom was being sought for absconding from parole and assault with a deadly weapon at the time he was killed. Autopsy results were not immediately available. He was paroled from prison last December for a 2013 offense of assault by strangulation or suffocation.</p>
			<p><em>Contact L.L. Brasier: 248-858-2262 or lbrasier@freepress.com</em></p>`,
	})
}

func TestIntegrationTigersManagerCandidates(t *testing.T) {
	runTest(t, &TestRec{
		url:       "http://www.freep.com/story/sports/mlb/tigers/2015/09/11/detroit-tigers-possible-manager-candidates/72051912/",
		expected:  []string{`First, it was the general manager.`, "Tigers general manager Al Avila today\u00a0said that no decision has been made."},
		forbidden: []string{`Contact Anthony Fenech`},
		html: `
			<p>First, it was the general manager.</p>
			<p>Next, all signs say it will be the manager.</p>
			<p>As the Detroit Tigers' rebooting rolls along – a coin termed by former president and general manager Dave Dombrowski, who was relieved of his duties Aug. 4 after his trade deadline deals starting the reboot – the first order of business for new general manager Al Avila will be finding a new manager.</p>
			<p>The Free Press learned late Thursday that longtime owner Mike Ilitch<a href="http://www.freep.com/story/sports/mlb/tigers/2015/09/11/detroit-tigers-brad-ausmus-fired/72048180/"> intends to fire second-year manager Brad Ausmus</a> after this season ends, according to a person with knowledge of the front office's plans. The person asked not to be identified because he was not authorized to discuss the situation publicly.</p>
			<p>Tigers general manager Al Avila today <a href="http://www.freep.com/story/sports/mlb/tigers/2015/09/11/brad-ausmus-detroit-tigers/72070960/">said that no decision has been made</a>.</p>
			<p>In Ausmus' two seasons with the Tigers, he is 154-148 (.510). Last season, the team went 90-72 and won the American League Central, falling in the first round of the playoffs to the Orioles. This season, they are 64-76 (.457) and in last place in the division.</p>
			<p>With the top man in the front office long gone and the top man in the clubhouse next out the door, Avila will look around the league to find the man to right the ship of the Tigers' reboot and will likely land one with extensive managerial experience and a successful track record. Here are a few candidates to be the 38th manager in team history:</p>
			<p><strong>■</strong>Ron Gardenhire: After a year off since being relieved of his duties with the Twins, Gardenhire, who turns 58 next month, could be ready for a new challenge within the American League Central. He compiled a 1,068-1,039 record in 13 seasons with the Twins, which included six division titles. He carries a strong resume of managerial experience – something that surely won't be overlooked as the organization's experience with the inexperienced Ausmus went south.</p>
			<p><strong>■</strong>Lloyd McClendon: The Mariners' manager was passed over for Ausmus in 2013 and hasn't found success in Seattle, but carries a strong reputation inside the Tigers' clubhouse – he served on Jim Leyland's staff in a number of roles. If the Mariners' new leadership elects to go in a different direction after the season, McClendon, 56, could carry the ties to make a strong run. He managed five seasons with the Pirates from 2001-05.</p>
			<p><strong>■</strong>Manny Acta: He spent three seasons with the Nationals from 2007-09 and then three with the Indians from 2010-12, never making the playoffs but never having the talent at his disposal to do so. He is the softer-spoken of a couple of Latino candidates the Tigers could consider, especially with the heavy influence of Latino players in the clubhouse. Acta, 46, hails from the Dominican Republic and has a 372-518 (.418) record as manager.</p>
			<p><strong>■</strong>Ozzie Guillen: If the front office went with a full-scale change to a firecracker of a name, the Tigers could consider Guillen, who has nine years of managerial experience, including a World Series to his name. Guillen, 51, is close friends with Miguel Cabrera – he also hails from Venezuela – and has one of the stronger personalities in the game. Guillen turned in five winning records in eight seasons with the White Sox from 2004-11 but flamed out after one season with the Marlins in 2012.</p>
			<p><strong>■</strong>A few familiar faces: It's highly unlikely that Tigers settle on another rookie manager, but first-base coach Omar Vizquel is well-liked in the clubhouse and profiles as a future manager in aspirations and baseball acumen. … Longtime teammates Alan Trammell and Kirk Gibson returned to the Tigers scene this season – Trammell as a special assistant to the general manager and Gibson as a broadcaster with Fox Sports Detroit – and both have managerial experience. Trammell was a Dombrowski hire and was handed a bad bunch, going 186-300 in three seasons before being fired after the 2005 season. Gibson, 58, managed the Diamondbacks to a National League West title in 2011 – his first full season – and checked in with two .500 seasons before getting fired late in the 2014 season with a 353-375 mark. With his recent diagnosis of Parkinson's disease, it's not known whether he is interested in managing again. … Jim Leyland, 70, is in his second season as a special assistant to the general manager after retiring from a 22-year career in the dugout and is likely done with the day-to-day grind.</p>
			<p><em>Contact Anthony Fenech: <a href="mailto:afenech@freepress.com">afenech@freepress.com</a>. Follow him on Twitter <a href="http://www.twitter.com/anthonyfenech">@anthonyfenech</a>. Check out our latest Tigers podcast at <a href="http://www.freep.com/tigerspodcast">freep.com/tigerspodcast</a> or on iTunes. And download our free Tigers Xtra app on <a href="https://itunes.apple.com/us/app/tigers-xtra/id962147637?mt=8">Apple</a> and <a href="https://play.google.com/store/apps/details?id=com.cincinnati.dolly.Tigers&amp;hl=en">Android</a>!</em></p>`,
	})
}

func TestIntegrationTigersMiguelComeback(t *testing.T) {
	runTest(t, &TestRec{
		url:       "http://www.freep.com/story/sports/mlb/tigers/2015/10/05/detroit-tigers-miguel-cabrera/73391112/",
		expected:  []string{`It could have been one of the most memorable comebacks in Los Angeles Angels history, but instead, it will go down as a mere footnote.`},
		forbidden: []string{},
		html: `
			<p>It could have been one of the most memorable comebacks in Los Angeles Angels history, but instead, it will go down as a mere footnote.</p>
			<p>When the Angels scored five runs in the ninth inning Saturday to beat the Texas Rangers, 11-10, they kept themselves in the postseason race for another day and prevented Texas from clinching the AL West. Then on Sunday, Texas beat the Angels, wrapping up the division and ensuring that Houston, not Los Angeles, would end up with a wild card.</p>
			<p>The sheer length of the baseball season makes it hard to tell the difference between a fleeting moment of glory and a true turning point. Fortunately, we now have the benefit of hindsight. So while that amazing ninth inning by the Angels ended up being fairly meaningless, here are four other moments that really did change the 2015 season:</p>
			<p><strong>May 17 — Jeff Banister shakes up his bullpen.</strong></p>
			<p>Banister, the Texas manager, told his relievers before a May 17 game against Cleveland that there were no set roles in the bullpen any more. The Rangers were 15-22 at that point, and closer Neftali Feliz already had blown three saves. Shawn Tolleson pitched the ninth for Texas that day, and the Rangers won, 5-1.</p>
			<p>A few days later, Tolleson earned the first save of his career. He would finish the season with 35 in 37 chances, adding stability to the late innings as the Rangers rallied to take the division by two games over Houston.</p>
			<p>Feliz was waived and then picked up<a href="http://www.freep.com/story/sports/mlb/tigers/2015/09/29/detroit-tigers-neftali-feliz/73016106/"> by the Detroit Tigers</a>.</p>
			<p><strong>May 21 — Jaime Garcia returns to the mound.</strong></p>
			<p>Garcia made only nine starts in 2013 and seven in 2014. He had thoracic outlet surgery in July 2014 to alleviate numbness and tingling in his pitching arm and hand. So it was fair to wonder what the St. Louis Cardinals could expect from him this year. But in his first game back, he allowed only two runs in seven innings against the New York Mets, an encouraging sign.</p>
			<p>Garcia ended up making 20 starts, going 10-6 with a 2.43 ERA. For a team that lost Adam Wainwright early on, it’s fair to suggest that Garcia’s performance was the difference between winning the NL Central and dropping to a wild card. The Cardinals won the division by two games.</p>
			<p><strong>July 3 — Miguel Cabrera <a href="http://www.freep.com/story/sports/mlb/tigers/2015/07/04/detroit-tigers-miguel-cabrera/29726763/">injures his left calf</a>.</strong></p>
			<p>The Detroit slugger would not play again until Aug. 14, and the team he came back to looked far different than the one he left. The Tigers went 15-20 in the interim and were in bad enough shape at the deadline that they traded stars David Price and Yoenis Cespedes. Price <a href="http://www.freep.com/story/sports/mlb/tigers/2015/10/01/toronto-blue-jays-david-price/73144536/">led Toronto to the AL East title</a>, and Cespedes <a href="http://www.freep.com/story/sports/mlb/tigers/2015/09/10/yoenis-cespedes-new-york-mets/72016556/">played a huge role for the New York Mets </a>in their NL East championship.</p>
			<p>Shortly after the deadline, the Tigers let general manager Dave Dombrowski go. He’s now <a href="http://www.freep.com/story/sports/columnists/drew-sharp/2015/08/20/detroit-tigers-dave-dombrowksi-drew-sharp/32023993/">running things in Boston</a>, so the butterfly effect from Cabrera’s injury could last awhile.</p>
			<p><strong>July 29 — The Mets don’t trade for Carlos Gomez.</strong></p>
			<p>They’ll be talking about this night in New York for years. Reports surfaced that Gomez was going to the Mets, and Wilmer Flores, who was expected to leave New York in the deal, was wiping tears from his eyes on the field during a game.</p>
			<p>The trade never was completed, though. Instead, the Mets kept Flores and traded for Cespedes. From July 31 on, Flores hit .296 with six homers. Cespedes hit .287 with 17 home runs and 44 RBIs in 57 games for New York. The Mets outlasted Washington in the NL East.</p>
			<p>Gomez, meanwhile, was traded to Houston and hit only .242 for the Astros, who made the playoffs but fell short in their bid for the AL West title.</p>`,
	})
}

func TestIntegrationShoplifter(t *testing.T) {
	runTest(t, &TestRec{
		url:       "http://www.freep.com/story/news/local/michigan/oakland/2015/10/06/cpl-holder-opens-fire-shoplifter-home-depot/73468588/",
		expected:  []string{},
		forbidden: []string{},
		html: `
			<p>A concealed-carry license holder is now cooperating with police after she opened fire on a shoplifter who was fleeing a Home Depot on Tuesday afternoon, Auburn Hills Police said.</p>
			<p>The shooting happened in the store’s parking lot at around 2 p.m., when Home Depot security was chasing a shoplifter in his 40s who jumped into a waiting dark SUV, said Lt. Jill McDonnell, an Auburn Hills police spokeswoman.</p>
			<p>But when the SUV began to pull away, a 48-year-old woman suddenly began firing shots at the fleeing vehicle. The vehicle escaped – but possibly has a flat tire, McDonnell said.</p>
			<p>The woman who fired the shots has a license to carry a firearm and is cooperating with police.</p>
			<p>It’s not clear whether the woman would face charges in the incident.</p>
			<p>The Home Depot, located on Joslyn, is part of a retail district with hundreds of stores in the area.</p>
			<p>The shooting comes just weeks after a bank customer in Warren opened fire on an armed robber, causing him to collapse from injuries. It’s not yet clear whether that man will face charges, either.</p>
			<p><em>Anyone with details on what happened can call the Auburn Hills Police Department at 248-370-9444.</em></p>
			<p><em>Contact Daniel Bethencourt: dbethencourt@freepress.com or 313-223-4531. Follow on Twitter at @_dbethencourt.</em></p>`,
	})
}

func TestIntegrationDingelHospital(t *testing.T) {
	runTest(t, &TestRec{
		url:       "http://www.freep.com/story/news/politics/2015/10/06/john-dingell-back-hospital/73464268/",
		expected:  []string{},
		forbidden: []string{`WASHINGTON`},
		html: `
			<p>WASHINGTON — Former U.S. Rep. John Dingell was admitted to Henry Ford Hospital on Monday and is expected to undergo a heart procedure, his wife's office said Tuesday.</p>
			<p>According to a brief statement from U.S. Rep. Debbie Dingell's office, her husband, who she replaced in Congress this year, went to the hospital in Detroit on Monday and is currently being evaluated. He is 89.</p>
			<p>Just before the congresswoman's office released the statement, her husband released his own statement via his Twitter feed, saying, "Back in the hospital. Being old sucks."</p>
			<p>No other information was immediately released about Dingell's health or circumstances, but Mrs. Dingell's office said he was "resting comfortably under doctor's care and is his usual feisty self." The office also said Debbie Dingell will remain in Michigan this week rather than traveling to Washington, though the House is in session.</p>
			<p>John Dingell retired from Congress early this year as its longest-serving member. Elected to replace his father, John Dingell Sr., in a special election in 1955, Dingell put together a career not only unparalleled in terms of longevity but helped write or pass much of the seminal legislation approved by Congress over the last 50 years.</p>
			<p>He was also one of its most powerful chairmen, sitting for years as the top Democrat atop the House Energy and Commerce Committee, a panel he helped expand into one of Congress' most influential.</p>
			<p>Dingell announced in early 2014 that he wouldn't run for a 30th full two-year term. Late last year, shortly after casting his last congressional votes, he spent three weeks in the hospital after suffering a hairline hip fracture.</p>
			<p><em>Contact Todd Spangler: 703-854-8947 or tspangler@freepress.com. Follow him on Twitter  @tsspangler.</em></p>`,
	})
}

func TestIntegrationBankRobbery(t *testing.T) {
	runTest(t, &TestRec{
		url:       "http://www.freep.com/story/news/local/michigan/wayne/2015/10/06/livonia-bank-robbery-arrest/73457536/",
		expected:  []string{},
		forbidden: []string{`Contact Robert Allen`},
		html: `
			<p>A thumbprint taken from a note demanding cash in a Livonia bank robbery led to the arrest and charging of a 29-year-old Romulus man, police said.</p>
			<p>Christopher Thomas Crowley was arraigned Monday on charges of bank robbery and cocaine possession after police found and arrested him in Inkster, according to a news release from Livonia Police Department.</p>
			<p>A robber Friday at the Citizens Bank on 31441 Plymouth Road passed the note to a teller. Described as a white man, about 30 years old, the robber fled the scene, leaving behind the note, police said. A latent print was lifted off the note and it was connected to Crowley, according to the news release.</p>
			<p>Charges include bank robbery, punishable by up to life in prison; two counts of possessing less than 25 grams of cocaine, punishable by 4 years in prison and possible fines up to $25,000, and habitual offender, a third-offense notice.</p>
			<p>Crowley's bond was set at $1-million cash or surety. His probable cause conference is set for Oct. 15, and his preliminary exam is set for Oct. 22.</p>
			<p><em>Contact Robert Allen at rallen@freepress.com or <a href="http://www.twitter.com/rallenMI">@rallenMI</a>. </em></p>`,
	})
}

/* template

func TestIntegrationXXXXX(t *testing.T) {
	runTest(t, &TestRec{
		url: "URL",
		expected: []string{},
		forbidden: []string{},
		html: ``,
	})
}
*/

func runTest(t *testing.T, rec *TestRec) {
	t.Logf("Testing URL: %v", rec.url)

	var extract *m.ExtractedBody
	if rec.html == "" {
		_, html, e, err := lib.ParseArticleAtURL(rec.url, true)
		extract = e
		if err != nil {
			t.Fatalf("Failed to parse article: %v", err)
		}

		println("Here's the HTML to embed for", rec.url)
		println(html)
	} else {
		extract = extraction.ExtractDataFromHTMLString(rec.html, rec.url, false)
	}

	text := extract.Text

	if text == "" {
		t.Errorf("Body extractor returned no text.")
	} else {
		for _, s := range rec.expected {
			if !strings.Contains(text, s) {
				t.Errorf("Expected body fragment not found: %#v", s)
			}
		}
		for _, s := range rec.forbidden {
			if strings.Contains(text, s) {
				t.Errorf("Forbidden body fragment found: %#v", s)
			}
		}
		for _, s := range globalForbidden {
			if strings.Contains(text, s) {
				t.Errorf("Globally forbidden body fragment found: %#v", s)
			}
		}
		t.Logf("in body: %#v", text)
	}
}
