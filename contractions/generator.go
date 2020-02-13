package contractions

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"strings"
	"text/template"
)

func write() error {
	data := make(map[string]string)

	cases := []func(string) string{
		strings.ToLower,
		title,
		strings.ToUpper,
	}

	for contraction, expansion := range contractions {
		for _, apostrophed := range apostrophes(contraction) {
			for _, f := range cases {
				key := f(apostrophed)
				existing, exists := data[key]
				if exists {
					return fmt.Errorf("attempting to re-add key %q (previous value was %q)", key, existing)
				}
				data[key] = f(expansion)
			}
		}
	}

	var source bytes.Buffer

	tmplErr := tmpl.Execute(&source, data)
	if tmplErr != nil {
		return tmplErr
	}

	formatted, fmtErr := format.Source(source.Bytes())
	if fmtErr != nil {
		return fmtErr
	}

	f, createErr := os.Create("generated.go")
	if createErr != nil {
		return createErr
	}
	defer f.Close()

	_, writeErr := f.Write(formatted)
	if writeErr != nil {
		return writeErr
	}

	return nil
}

const a = "'"

func apostrophes(s string) []string {
	result := []string{s}
	if strings.Contains(s, a) {
		// smart quote variation
		result = append(result, strings.Replace(s, a, "’", -1))
	}
	return result
}

func title(s string) string {
	head := strings.ToUpper(s[:1])
	tail := s[1:]
	return head + tail
}

var tmpl = template.Must(template.New("").Parse(`
package contractions

// This file is generated. Best not to modify it, as it will likely be overwritten.

// maps do not guarantee order, so this will look random
var variations = {{ printf "%#v" . }}
`))

var contractions = map[string]string{
	// prefer a map of explicit lookups, vs some logic/loop to generalize
	// downside: to handle case consistently, it gets verbose
	"i'll":    "i will",
	"you'll":  "you will",
	"she'll":  "she will",
	"he'll":   "he will",
	"we'll":   "we will",
	"they'll": "they will",

	"i'm":     "i am",
	"you're":  "you are",
	"she's":   "she is", // arguably 'she has', would need parts-of-speech to determine
	"he's":    "he is",
	"we're":   "we are",
	"they're": "they are",

	"i've":    "i have",
	"you've":  "you have",
	"we've":   "we have",
	"they've": "they have",

	"i'd":    "i would", // arguably "i had"
	"you'd":  "you would",
	"she'd":  "she would",
	"he'd":   "he would",
	"we'd":   "we would",
	"they'd": "they would",

	"isn't":  "is not",
	"aren't": "are not",
	"wasn't": "was not",

	"don't":   "do not",
	"doesn't": "does not",
	"didn't":  "did not",

	"haven't": "have not",
	"hadn't":  "had not",

	"can't": "can not",

	"won't":   "will not",
	"will've": "will have",

	"wouldn't": "would not",
	"would've": "would have",

	"couldn't": "could not",
	"could've": "could have",

	"shouldn't": "should not",
	"should've": "should have",

	"mightn't": "might not",
	"might've": "might have",

	"mustn't": "must not",
	"must've": "must have",

	"gonna":  "going to",
	"gotta":  "got to",
	"wanna":  "want to",
	"gimme":  "give me",
	"cannot": "can not",
}
