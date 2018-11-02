package main

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	body := strings.NewReader(`
	We can lemmatize plain prose, perhaps a technical job listing: 
	
	We are looking for experienced Rails developers, with experience in HTML 5 and T-SQL.
	
	Experience with ObjC and cpp and vue and React  Native is a plus.
	`)

	req := httptest.NewRequest("POST", "/text", body)
	w := httptest.NewRecorder()

	jargonHandler(w, req)

	resp := w.Result()
	result, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		t.Error(err)
	}

	got := string(result)

	if !strings.Contains(got, "ruby-on-rails") {
		t.Errorf(`should have found "ruby-on-rails" in result, got %q`, got)
	}

	if !strings.Contains(got, `<span class="lemma">objective-c</span>`) {
		t.Errorf(`should have found <span class="lemma">objective-c</span> in result, got %q`, got)
	}
}
