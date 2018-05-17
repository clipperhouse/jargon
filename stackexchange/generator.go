// Package stackexchange implements a jargon.Dictionary of tags and synonyms for use with jargon lemmatizers
package stackexchange

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
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
	pages := 10
	pageSize := 100

	// Temporary, for detecting duplicates
	tagcheck := make(map[string]struct{})

	// Used in the template below
	data := struct {
		Tags     []string
		Synonyms map[string]string
	}{
		Tags:     make([]string, 0),
		Synonyms: make(map[string]string),
	}

	for _, site := range sites {
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
				name := trailingVersion.ReplaceAllString(item.Name, "")
				// Only append tags if they haven't been added already
				if _, found := tagcheck[name]; !found {
					tagcheck[name] = exists
					data.Tags = append(data.Tags, name)
				}

				// Only add synonyms if they haven't been added already
				for _, synonym := range item.Synonyms {
					key := trailingVersion.ReplaceAllString(synonym, "")
					_, found := data.Synonyms[key]
					if !found && key != name {
						data.Synonyms[synonym] = name
					}
					// What if the same synonym points to multiple canonical tags?
					// The above logic just uses the first one, maybe that's ok.
				}
			}

			if !wrapper.HasMore {
				break
			}
		}
	}

	t := template.Must(template.New("dict").Parse(tmpl))

	var source bytes.Buffer

	tmplErr := t.Execute(&source, data)
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

var tmpl = `
package stackexchange

// This file is generated. Best not to modify it, as it will likely be overwritten.

// Dictionary is the main exported Dictionary of Stack Exchange tags and synonyms, fetched via api.stackexchange.com
// It includes the most popular {{ .Tags | len }} tags and their {{ .Synonyms | len }} synonyms
var Dictionary = &dictionary{ 
	tags: tags, 
	synonyms: synonyms,
	// MaxGramLength is hard-coded in dictionary.go
}

var tags = {{ printf "%#v" .Tags }}

var synonyms = {{ printf "%#v" .Synonyms }}
`

var sites = []string{"stackoverflow", "serverfault"}
var tagsURL = "http://api.stackexchange.com/2.2/tags?page=%d&pagesize=%d&order=desc&sort=popular&site=%s&filter=!4-J-du8hXSkh2Is1a&page=%d"
var client = http.Client{
	Timeout: time.Second * 2, // Maximum of 2 secs
}
var empty = wrapper{}

func fetchTags(page, pageSize int, site string) (wrapper, error) {
	if page == 0 {
		page = 1
	}

	if pageSize == 0 {
		pageSize = 100
	}

	url := fmt.Sprintf(tagsURL, page, pageSize, site)
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

	return wrapper, nil
}

type item struct {
	Name      string   `json:"name"`
	Synonyms  []string `json:"synonyms"`
	Moderator bool     `json:"is_moderator_only"`
}

type wrapper struct {
	Items   []item `json:"items"`
	HasMore bool   `json:"has_more"`
}
