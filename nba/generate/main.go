package main

import (
	"bytes"
	"fmt"
	"go/format"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/PuerkitoBio/goquery"
	"github.com/clipperhouse/jargon/ascii"
)

func main() {
	names, err := fetchNames()
	check(err)

	mappings, err := getMappings(names)
	check(err)

	err = write(mappings)
	check(err)
}

func fetchNames() ([]string, error) {
	var names []string

	// Fetch Wikipedia page
	resp, err := http.Get("https://en.wikipedia.org/wiki/List_of_current_NBA_team_rosters")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	links := doc.Find("table.toccolours table.sortable tbody > tr > td:nth-of-type(3) > a")
	links.Each(func(i int, s *goquery.Selection) {
		names = append(names, s.Text())
	})

	return names, nil
}

func getMappings(names []string) (map[string]string, error) {
	mappings := map[string]string{}

	for _, name := range names {
		// Switch to first last
		split := strings.Split(name, ",")
		if len(split) != 2 {
			return nil, fmt.Errorf("expected two parts after splitting on name %q, got %d", name, len(split))
		}

		var synonyms []string
		canonical := strings.TrimSpace(split[1]) + " " + strings.TrimSpace(split[0])
		synonyms = append(synonyms, canonical)

		fold, folded := ascii.FoldString(canonical)
		if folded {
			synonyms = append(synonyms, fold)
		}

		key := strings.Join(synonyms, ", ")
		mappings[key] = canonical
	}

	return mappings, nil
}

func write(mappings map[string]string) error {
	var source bytes.Buffer

	err := tmpl.Execute(&source, mappings)
	if err != nil {
		return err
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
package nba

// This file is generated. Best not to modify it, as it will likely be overwritten.

var mappings = {{ printf "%#v" . }}
`))

func check(err error) {
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Stderr.WriteString("\n")
		os.Exit(1)
	}
}
