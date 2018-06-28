// Package stackexchange implements a jargon.Dictionary of tags and synonyms for use with jargon lemmatizers
package stackexchange

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"
	"time"
)

var exists = struct{}{}

// Hyphenated trailing versions, like python-2.7. Things like html5 are considered a unique name, not a version per se.
var trailingVersion = regexp.MustCompile(`-[\d.]+$`)

// run this in generator_test.go
func writeDictionary() error {
	pageSize := 100

	// Used in the template below
	data := struct {
		Tags     map[string]string
		Synonyms map[string]string
	}{
		Tags:     make(map[string]string),
		Synonyms: make(map[string]string),
	}

	for site, count := range sites {
		pages := count / pageSize
		for page := 1; page <= pages; page++ {
			wrapper, err := fetchTags(page, pageSize, site)
			if err != nil {
				return err
			}

			// Avoid duplicates; since we are querying multiple sites, duplication is common.
			for _, item := range wrapper.Items {
				if item.Moderator {
					// We don't want those
					continue
				}

				// Trailing versions (like ruby-on-rails-4) are not interesting for our purposes
				canonical := trailingVersion.ReplaceAllString(item.Name, "")

				if isStopWord(canonical) {
					continue // skip it
				}

				key := normalize(canonical)
				// Only append tags if they haven't been added already
				if _, found := data.Tags[key]; !found {
					data.Tags[key] = canonical
				}

				// Only add synonyms if they haven't been added already
				for _, s := range item.Synonyms {
					synonym := trailingVersion.ReplaceAllString(s, "")

					if synonym == canonical {
						continue // skip it
					}

					if isStopWord(synonym) {
						continue // skip it
					}

					key := normalize(synonym)
					if _, found := data.Synonyms[key]; !found {
						data.Synonyms[key] = canonical
					}
				}
			}

			if !wrapper.HasMore {
				break
			}

			if wrapper.Backoff > 10 {
				// That's too much for this run
				err := fmt.Errorf("Abort: received a message to backoff %d seconds from api.stackexchange.com. That's too much, try again later. See http://api.stackexchange.com/docs/throttle", wrapper.Backoff)
				return err
			}

			// Try to avoid throttling
			if wrapper.Backoff > 0 {
				backoff := time.Duration(wrapper.Backoff) * time.Second
				log.Printf("Backing off %d seconds, per backoff message from api.stackexchange.com. See http://api.stackexchange.com/docs/throttle", wrapper.Backoff)
				time.Sleep(backoff)
			}
			// A little extra to be safe
			// Guideline in the documentation is 30 requests/sec: http://api.stackexchange.com/docs/throttle
			time.Sleep(50 * time.Millisecond)
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

var tmpl = template.Must(template.New("").Parse(`
package stackexchange

// This file is generated. Best not to modify it, as it will likely be overwritten.

// Dictionary is the main exported Dictionary of Stack Exchange tags and synonyms, from the following Stack Exchange sites: Stack Overflow,
// Server Fault, Game Dev and Data Science. It's indended to identify canonical tags (technologies),
// e.g. Ruby on Rails (3 words) will be replaced with ruby-on-rails (1 word).
// It includes the most popular {{ .Tags | len }} tags and {{ .Synonyms | len }} synonyms
var Dictionary = &dictionary{ 
	tags: tags, 
	synonyms: synonyms,
}

var tags = {{ printf "%#v" .Tags }}

var synonyms = {{ printf "%#v" .Synonyms }}
`))

// sites to query, with the number of tags to get, based on eyeballing how many of the top x are 'interesting'
var sites = map[string]int{
	"stackoverflow": 2000,
	"serverfault":   600,
	"gamedev":       300,
	"datascience":   200,
}
var tagsURL = "http://api.stackexchange.com/2.2/tags?page=%d&pagesize=%d&order=desc&sort=popular&site=%s&filter=!4-J-du8hXSkh2Is1a&key=%s"
var stackExchangeAPIKey = "*AbAX7kb)BKJTlmKgb*Tkw(("

var client = http.Client{
	Timeout: time.Second * 3, // Maximum of 2 secs
}
var empty = wrapper{}

func fetchTags(page, pageSize int, site string) (wrapper, error) {
	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 100
	}

	url := fmt.Sprintf(tagsURL, page, pageSize, site, stackExchangeAPIKey)
	r, httpErr := client.Get(url)
	if httpErr != nil {
		return empty, httpErr
	}

	defer r.Body.Close()

	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		return empty, readErr
	}

	wrapper := wrapper{}
	jsonErr := json.Unmarshal(body, &wrapper)
	if jsonErr != nil {
		return empty, jsonErr
	}

	if len(wrapper.ErrorName) > 0 || len(wrapper.ErrorMessage) > 0 {
		return wrapper, fmt.Errorf("%s: %s", wrapper.ErrorName, wrapper.ErrorMessage)
	}

	return wrapper, nil
}

type item struct {
	Name      string   `json:"name"`
	Synonyms  []string `json:"synonyms"`
	Moderator bool     `json:"is_moderator_only"`
}

type wrapper struct {
	Items        []item `json:"items"`
	HasMore      bool   `json:"has_more"`
	ErrorName    string `json:"error_name"`
	ErrorMessage string `json:"error_message"`
	// In seconds
	Backoff int `json:"backoff"`
}
