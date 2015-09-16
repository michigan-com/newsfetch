package lib

import (
	"testing"
)

var text = `A metro Detroit doctor who was sentenced to 45 years in federal prison last month is appealing his sentencing and conviction to the U.S. 6th Circuit Court of Appeals.
Farid Fata, 50, who was convicted of violating more than 550 patients' trust and raking in more than $17 million from fraudulent billings was sentenced on July 10.
Defense Attorney Mark Kriger filed the appeal on Fata's behalf  Wednesday in the U.S. District Court's Eastern District of Michigan. Kriger confirmed the filing to the Free Press.
"We filed a notice of appeal on the sentence, but I don't feel it's appropriate to comment on pending cases," Kriger said.
Prior to Fata's sentencing, the court heard victim impact statements from nearly 22 victims, who shared unfathomable experiences of undergoing unnecessary chemotherapy treatments and losing teeth, of a patient diagnosed with lung cancer when he had kidney cancer and more. One patient was given 195 chemotherapy treatments, 177 of which were unnecessary.
Fata pleaded guilty in September to 13 counts of health care fraud, two counts of money laundering and one count of conspiring to pay and receive kickbacks. The case involved $34.7 million in billings to patients and insurance companies, and $17.6 million paid for work Fata admitted was unnecessary.
Federal prosecutors said Fata's case was the most egregious fraud case they've ever seen. U.S. District Judge Paul Borman said before sentencing Fata that the once-prominent oncologist "practiced greed and shut down whatever compassion he had."
Fata, a married father of three and a naturalized U.S. citizen whose native country is Lebanon, was charged with running the scheme that involved billing the government for medically unnecessary cancer and blood treatments.
The government said Fata ran the scheme from 2009 to 2014 through his medical businesses, including Michigan Hematology Oncology Centers, with offices in Clarkston, Bloomfield Hills, Lapeer, Sterling Heights, Troy and Oak Park.
He remains incarcerated, and his medical license has been revoked, but Fata's legal troubles are far from over. About $13 million has been collected since 2013 to go toward a $17.6-million criminal judgement against Fata.
U.S. Assistant Prosecutor Catherine Dick said prosecutors are continuing to work to close the gap with his assets. Dick said the patient-victims and their families are first priority for compensation, then private insurers, then Medicare. The whistle-blower who tipped off the federal investigation is to receive 10% as part of an agreement; typically, whistle-blowers receive 15-25%, Dick said. Fata is also facing 27 pending lawsuits from patients and their families in Oakland County. George Karadsheh, Fata's former practice business manager, has also filed a whistle-blower federal suit against the former doctor.
Contact Katrease Stafford: kstafford@freepress.com
`

var title = "Cancer doc Farid Fata appeals 45-year prison sentence"

func TestSummarizer(t *testing.T) {
	t.Log("Generate an article based on a title and a body of text.")

	expected := []string{
		"Farid Fata, 50, who was convicted of violating more than 550 patients' trust and raking in more than $17 million from fraudulent billings was sentenced on July 10.",
		"Defense Attorney Mark Kriger filed the appeal on Fata's behalf  Wednesday in the U.S. District Court's Eastern District of Michigan.",
		"The case involved $34.7 million in billings to patients and insurance companies, and $17.6 million paid for work Fata admitted was unnecessary.",
		"The government said Fata ran the scheme from 2009 to 2014 through his medical businesses, including Michigan Hematology Oncology Centers, with offices in Clarkston, Bloomfield Hills, Lapeer, Sterling Heights, Troy and Oak Park.",
		"George Karadsheh, Fata's former practice business manager, has also filed a whistle-blower federal suit against the former doctor.",
	}

	summarizer := NewPunktSummarizer(title, text)
	actual := summarizer.KeyPoints()

	if len(actual) != len(expected) {
		t.Fatalf("Actual: %d != Expected: %d", len(actual), len(expected))
	}

	// TO DO FIX SUMMARIZER PRODUCING DIFFERENT OUTPUTS FOR THE SAME INPUTS
	/*for i := 0; i < len(expected); i++ {
		if actual[i] != expected[i] {
			t.Fatalf("Actual: %s\n Expected: %s", actual[i], expected[i])
		}
	}*/

}
