package jargon

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestTechProse(t *testing.T) {
	text := `Hi! This is a test of tech terms.
It should consider F#, C++, .net, Node.JS and 3.141592 to be their own tokens. 
Similarly, #hashtag and @handle should work, as should an first.last+@example.com.
It should—wait for it—break on things like em-dashes and "quotes" and it ends.
It'd be great it it’ll handle apostrophes.
`
	r := strings.NewReader(text)
	got := collect(Tokenize(r))

	expected := []string{
		"Hi", "!",
		"F#", "C++", ".net", "Node.JS", "3.141592",
		"#hashtag", "@handle", "first.last+@example.com",
		"should", "—", "wait", "it", "break", "em-dashes", "quotes",
		"It'd", "it’ll", "apostrophes",
	}

	for _, e := range expected {
		if !contains(e, got) {
			t.Errorf("Expected to find token %q, but did not.", e)
		}
	}

	// Check that last .
	nextToLast := got[len(got)-2]
	if nextToLast.String() != "." {
		t.Errorf("The next-to-last token should be %q, got %q.", ".", nextToLast)
	}

	// Check that last \
	last := got[len(got)-1]
	if last.String() != "\n" {
		t.Errorf("The last token should be %q, got %q.", "\n", last)
	}

	// No trailing punctuation
	for _, token := range got {
		if utf8.RuneCountInString(token.String()) == 1 {
			// Skip actual (not trailing) punctuation
			continue
		}

		if strings.HasSuffix(token.String(), ",") || strings.HasSuffix(token.String(), ".") {
			t.Errorf("Found trailing punctuation in %q", token.String())
		}
	}
}

func BenchmarkProse(b *testing.B) {
	text := `Maximize, bleeding-edge cultivate engage B2B ecologies rich-clientAPIs viral embrace engage integrated systems morph, ecologies, scalable user-contributed. Redefine synergies, architectures redefine value; empower Cluetrain capture eyeballs optimize scale tag architect, addelivery. Out-of-the-box monetize; best-of-breed rich; webservices architectures impactful social exploit open-source mashups sexy morph, harness; killer initiatives benchmark e-business infomediaries channels enable e-business, eyeballs functionalities end-to-end transform remix applications empower reinvent. Action-items synergistic: networkeffects citizen-media platforms. Communities, intuitive, efficient tag synergistic post integrate weblogs networkeffects embedded impactful--deliverables peer-to-peer wikis strategize out-of-the-box, capture--syndicate engineer holistic! Monetize collaborative productize solutions enterprise webservices 24/7 data-driven blogospheres implement incentivize citizen-media value-added morph methodologies benchmark networkeffects, intuitive synergies user-contributed 24/365 granular sticky leverage. Extensible ROI convergence innovate turn-key; markets ubiquitous mindshare: networking, platforms."

	Global compelling metrics markets morph communities capture webservices leading-edge B2B iterate enhance. Best-of-breed exploit scalable, B2C synergies web services mindshare seamless semantic deliver, productize rich maximize! Partnerships implement mindshare dot-com engineer e-business, solutions. Front-end aggregate enterprise visionary communities incubate, ubiquitous addelivery innovative markets engage.
	
	Platforms technologies visionary--value, supply-chains B2B infomediaries repurpose monetize methodologies front-end, functionalities; share empower users. Real-time impactful reinvent share, grow orchestrate syndicate rich synergies, global convergence, "seize front-end, widgets data-driven one-to-one networks convergence value infrastructures clicks-and-mortar transition." Whiteboard transparent interactive blogging peer-to-peer deliver web-readiness communities integrate optimize compelling: action-items value leading-edge supply-chains design. Semantic magnetic maximize rich infrastructures value-added plug-and-play sexy networking implement enhance communities aggregate, frictionless. Optimize architect long-tail channels blogospheres, "vortals e-business functionalities." Mindshare e-business remix webservices user-contributed integrate holistic eyeballs vertical, 24/365 e-enable, niches, users niches reinvent recontextualize.
	
	Monetize, matrix bleeding-edge syndicate engineer drive e-business synthesize embrace revolutionary share podcasts repurpose: impactful convergence, dynamic vortals global harness. Initiatives integrate, strategize visualize; metrics seamless clicks-and-mortar customized strategize webservices user-centric impactful integrate mindshare.
	
	Platforms efficient supply-chains supply-chains, "life-hacks networking design integrated beta-test technologies e-commerce!" Matrix cross-platform proactive rss-capable portals B2C value-added reinvent strategic enable scalable. Webservices evolve share redefine empower cross-platform communities recontextualize webservices visualize, viral relationships widgets robust methodologies aggregate peer-to-peer integrated. Scalable robust web-enabled remix blogospheres magnetic: seamless revolutionize life-hacks bricks-and-clicks portals supply-chains, disintermediate mindshare. Reinvent channels, authentic citizen-media social leading-edge deliver, dynamic; niches action-items action-items, methodologies, cultivate revolutionize integrated standards-compliant, e-markets, peer-to-peer enterprise matrix mindshare. B2C architectures leverage seamless turn-key open-source, orchestrate syndicate seize! Harness, cultivate, integrate; networks, embedded, "networks expedite semantic; viral beta-test one-to-one evolve interactive real-time B2C."`

	r := strings.NewReader(text)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		consume(Tokenize(r))
	}
}

func TestURLs(t *testing.T) {
	// We mostly get lucky on URLs due to punct rules

	tests := map[string]string{
		`http://www.google.com`:                     `http://www.google.com`,                     // as-is
		`http://www.google.com.`:                    `http://www.google.com`,                     // "."" should be considered trailing punct
		`http://www.google.com/`:                    `http://www.google.com/`,                    // trailing slash OK
		`http://www.google.com/?`:                   `http://www.google.com/`,                    // "?" should be considered trailing punct
		`http://www.google.com/?foo=bar`:            `http://www.google.com/?foo=bar`,            // "?" is querystring
		`http://www.google.com/?foo=bar.`:           `http://www.google.com/?foo=bar`,            // trailing "."
		`http://www.google.com/?foo=bar&qaz=qux`:    `http://www.google.com/?foo=bar&qaz=qux`,    // "?" with &
		`http://www.google.com/?foo=bar&qaz=q%20ux`: `http://www.google.com/?foo=bar&qaz=q%20ux`, // with encoding
		`//www.google.com`:                          `//www.google.com`,                          // scheme-relative
		`/usr/local/bin/foo.bar`:                    `/usr/local/bin/foo.bar`,
		`c:\windows\notepad.exe`:                    `c:\windows\notepad.exe`,
	}

	for input, expected := range tests {
		r := strings.NewReader(input)
		got := <-Tokenize(r) // just take the first token

		if got.String() != expected {
			t.Errorf("Expected URL %s to result in %s, but got %s", input, expected, got)
		}
	}
}

func TestTechHTML(t *testing.T) {
	h := `
<html>
<p foo="bar">
Hi! Let's talk Ruby on Rails.
<!-- Ignore ASPNET MVC in comments -->
</p>
</html>
`
	r := strings.NewReader(h)
	got := collect(TokenizeHTML(r))

	expected := []string{
		`<p foo="bar">`, // tags kept whole
		"\n",            // whitespace preserved
		"Hi", "!",
		"Ruby", "on", "Rails", // make sure text node got tokenized
		"<!-- Ignore ASPNET MVC in comments -->", // make sure comment kept whole
		"</p>",
	}

	for _, e := range expected {
		if !contains(e, got) {
			t.Errorf("Expected to find token %q, but did not.", e)
		}
	}
}

func contains(value string, tokens []Token) bool {
	for _, t := range tokens {
		if t.String() == value {
			return true
		}
	}
	return false
}

// Checks that value, punct and space are equal for two slices of token; deliberately does not check lemma
func equals(a, b []Token) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].String() != b[i].String() {
			return false
		}
		if a[i].IsPunct() != b[i].IsPunct() {
			return false
		}
		if a[i].IsSpace() != b[i].IsSpace() {
			return false
		}
		// deliberately not checking for IsLemma(); use reflect.DeepEquals
	}

	return true
}
