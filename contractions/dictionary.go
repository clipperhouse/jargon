package contractions

type dictionary struct{}

// Dictionary for expanding common contractions into distinct words. Examples:
//	don't → does not
//	we've → we have
//	she's -> she is
// Caveats:
// - only lower case right now; TODO: add support for title case and all uppercase?
// - returns expanded words as a single string with a space in it; caller might wish to re-tokenize
var Dictionary = &dictionary{}

// Lookup attempts to convert single-token contractions to non-contracted version.  Examples:
//	don't → does not
//	we've → we have
//	she's -> she is
// Caveats:
// - only lower case right now
// - returns expanded words as a single token with a space in it; caller might wish to re-tokenize
func (d *dictionary) Lookup(s []string) (string, bool) {
	if len(s) != 1 {
		return "", false
	}

	canonical, ok := contractions[s[0]]
	if !ok {
		return "", false
	}

	return canonical, true
}

// TODO: preserve case, i.e. "You'll" → "You will"

var contractions = map[string]string{
	// prefer a map of explicit lookups, vs some logic/loop to generalize
	// downside: to handle case consistently, it gets verbose
	"i'll":    "i will",
	"i’ll":    "i will",
	"you'll":  "you will",
	"you’ll":  "you will",
	"she'll":  "she will",
	"she’ll":  "she will",
	"he'll":   "he will",
	"he’ll":   "he will",
	"we'll":   "we will",
	"we’ll":   "we will",
	"they'll": "they will",
	"they’ll": "they will",

	"i'm":     "i am",
	"i’m":     "i am",
	"you're":  "you are",
	"you’re":  "you are",
	"she's":   "she is", // arguably 'she has', would need parts-of-speech to determine
	"she’s":   "she is",
	"he's":    "he is",
	"he’s":    "he is",
	"we're":   "we are",
	"we’re":   "we are",
	"they're": "they are",
	"they’re": "they are",

	"i've":    "i have",
	"i’ve":    "i have",
	"you've":  "you have",
	"you’ve":  "you have",
	"we've":   "we have",
	"we’ve":   "we have",
	"they've": "we have",
	"they’ve": "we have",

	"i'd":    "i would", // arguably "i had"
	"i’d":    "i would",
	"you'd":  "you would",
	"you’d":  "you would",
	"she'd":  "she would",
	"she’d":  "she would",
	"he'd":   "he would",
	"he’d":   "he would",
	"we'd":   "we would",
	"we’d":   "we would",
	"they'd": "we would",
	"they’d": "we would",

	"isn't":  "is not",
	"isn’t":  "is not",
	"aren't": "are not",
	"aren’t": "are not",
	"wasn't": "was not",
	"wasn’t": "was not",

	"don't":   "do not",
	"don’t":   "do not",
	"doesn't": "does not",
	"doesn’t": "does not",
	"didn't":  "did not",
	"didn’t":  "did not",

	"haven't": "have not",
	"haven’t": "have not",
	"hadn't":  "had not",
	"hadn’t":  "had not",

	"can't": "can not",
	"can’t": "can not",

	"won't":   "will not",
	"won’t":   "will not",
	"will've": "will have",
	"will’ve": "will have",

	"wouldn't": "would not",
	"wouldn’t": "would not",
	"would've": "would have",
	"would’ve": "would have",

	"couldn't": "could not",
	"couldn’t": "could not",
	"could've": "could have",
	"could’ve": "could have",

	"shouldn't": "should not",
	"shouldn’t": "should not",
	"should've": "should have",
	"should’ve": "should have",

	"mightn't": "might not",
	"mightn’t": "might not",
	"might've": "might have",
	"might’ve": "might have",

	"mustn't": "must not",
	"mustn’t": "must not",
	"must've": "must have",
	"must’ve": "must have",

	"gonna":  "going to",
	"gotta":  "got to",
	"wanna":  "want to",
	"gimme":  "give me",
	"cannot": "can not",
}
