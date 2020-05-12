package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

func main() {
	err := writeDictionary()

	if err != nil {
		panic(err)
	}
}

// Hard-coded exceptions that cause bad mappings
var ignore = map[string]bool{
	"drop-down": true,
	"datatable": true,
	"for":       true,
	"this":      true,
}

// run this in generator_test.go
func writeDictionary() error {
	pageSize := 100

	// Used in the template below
	mappings := make(map[string]string)
	seen := make(map[string]bool)

	for site, count := range sites {
		pages := count / pageSize
		for page := 1; page <= pages; page++ {
			wrapper, err := fetchTags(page, pageSize, site)
			if err != nil {
				return err
			}

			// Synonyms first
			for _, tag := range wrapper.Items {
				if tag.Moderator {
					// We don't want those
					continue
				}

				if len(tag.Synonyms) == 0 {
					continue
				}

				canonical := tag.Name // tag name

				var filtered []string
				for _, synonym := range tag.Synonyms {
					if ignore[synonym] {
						continue
					}

					// Split up the grams to allow calculation max gram length by the Synonyms constructor
					synonym = strings.ReplaceAll(synonym, "-", " ")

					filtered = append(filtered, synonym)
				}

				if len(filtered) == 0 {
					// Skip it
					continue
				}

				for _, synonym := range filtered {
					seen[synonym] = true
				}

				// Comma-separated string, a format Synonyms can handle
				synonyms := strings.Join(filtered, ", ")
				mappings[synonyms] = canonical
			}

			// Avoid duplicates; since we are querying multiple sites, duplication is common.
			for _, tag := range wrapper.Items {
				if tag.Moderator {
					// We don't want those
					continue
				}

				canonical := tag.Name // tag name

				// Split up the grams to allow calculation max gram length by the Synonyms constructor
				synonym := strings.ReplaceAll(tag.Name, "-", " ")

				if ignore[synonym] {
					continue
				}

				if seen[synonym] {
					// Ignore ones we've seen
					continue
				}

				// We want to identify tags as canonical even if they don't map to something different
				mappings[synonym] = canonical
			}

			if !wrapper.HasMore {
				break
			}

			if wrapper.Backoff > 10 {
				// That's too much for this run
				err := fmt.Errorf("abort: received a message to backoff %d seconds from api.stackexchange.com. That's too much, try again later. See http://api.stackexchange.com/docs/throttle", wrapper.Backoff)
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

	tmplErr := tmpl.Execute(&source, mappings)
	if tmplErr != nil {
		return tmplErr
	}

	// Break up some lines for readability
	split := strings.ReplaceAll(source.String(), `", "`, `",
"`)
	split = strings.ReplaceAll(split, `{"`, `{
"`)
	split = strings.ReplaceAll(split, `"}`, `",
}`)

	formatted, fmtErr := format.Source([]byte(split))
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
package stackoverflow

// This file is generated. Best not to modify it, as it will likely be overwritten.

var mappings = {{ printf "%#v" . }}
`))

// sites to query, with the number of tags to get, based on eyeballing how many of the top x are 'interesting'
var sites = map[string]int{
	"stackoverflow": 2000,
	// "serverfault":   600,
	// "gamedev":       300,
	// "datascience":   200,
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
