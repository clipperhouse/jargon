package main

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"github.com/clipperhouse/jargon"
	"github.com/clipperhouse/jargon/ascii"
	"github.com/clipperhouse/jargon/contractions"
	"github.com/clipperhouse/jargon/stackoverflow"
	"github.com/clipperhouse/jargon/stemmer"
	"github.com/spf13/afero"
)

var testfilein = "/tmp/in.txt"

func testConfig() (config, error) {
	c := config{}
	fs := afero.NewMemMapFs()

	file, err := fs.Create(testfilein)
	if err != nil {
		return c, err
	}

	_, err = file.WriteString("test input file")
	if err != nil {
		return c, err
	}

	c.Fs = fs
	return c, nil
}

func TestInput(t *testing.T) {
	type test struct {
		// input
		filein string
		mode   os.FileMode

		// expected
		err     error
		pipedin bool
		file    bool
	}

	tests := []test{
		{
			// File, not piped
			filein: testfilein,
			mode:   os.ModeCharDevice,

			err:     nil,
			pipedin: false,
			file:    true,
		},
		{
			// Piped, not file
			filein: "",
			mode:   os.ModeAppend,

			err:     nil,
			pipedin: true,
			file:    false,
		},
		{
			// Not piped, not file
			filein: "",
			mode:   os.ModeCharDevice,

			err:     errNoInput,
			pipedin: false,
			file:    false,
		},
		{
			// Both piped and file
			filein: testfilein,
			mode:   os.ModeAppend,

			err:     errTwoInput,
			pipedin: true,
			file:    true,
		},
		{
			// Both piped and file
			filein: "doesntexist",
			mode:   os.ModeAppend,

			err:     os.ErrNotExist,
			pipedin: false,
			file:    false,
		},
	}

	for _, test := range tests {
		c, err := testConfig()
		if err != nil {
			t.Error(err)
		}

		err = setInput(&c, test.mode, test.filein)
		if !errors.Is(err, test.err) {
			t.Errorf("expected err %v, got %v", test.err, err)
		}
		if c.Pipedin != test.pipedin {
			t.Errorf("expected piped to be %t, got %t", test.pipedin, c.Pipedin)
		}
		if (c.Filein != nil) != test.file {
			t.Errorf("expected c.Filein to be %t", test.file)
		}
		if c.Filein != nil {
			err := c.Filein.Close()
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestFilters(t *testing.T) {
	type test struct {
		// input
		args []string
		lang string

		// expected
		err     bool
		filters []jargon.Filter
	}

	tests := []test{
		{
			args: []string{"-stack", "-stem", "-ascii", "-contractions"},
			lang: "",

			err:     false,
			filters: []jargon.Filter{stackoverflow.Tags, stemmer.English, ascii.Fold, contractions.Expander},
		},
		{
			args: []string{"-stem"},
			lang: "spanish",

			err:     false,
			filters: []jargon.Filter{stemmer.Spanish},
		},
		{
			args: []string{"-stem"},
			lang: "foo",

			err:     true,
			filters: nil,
		},
	}
	for _, test := range tests {
		c, err := testConfig()
		if err != nil {
			t.Error()
		}

		err = setFilters(&c, test.args, test.lang)
		if (err != nil) != test.err {
			t.Errorf("expected err %v, got %v", test.err, err)
		}
		if !reflect.DeepEqual(c.Filters, test.filters) {
			t.Errorf("expected filters to match, args: %v, lang: %s", test.args, test.lang)
		}
	}
}

func TestOutput(t *testing.T) {
	type test struct {
		// input
		fileout string

		// expected
		err      error
		pipedout bool
		file     bool
	}

	testfileout := "testdata/out.txt"

	tests := []test{
		{
			// File, not piped
			fileout: testfileout,

			err:      nil,
			pipedout: false,
			file:     true,
		},
		{
			// Piped, not file
			fileout: "",

			err:      nil,
			pipedout: true,
			file:     false,
		},
	}

	for _, test := range tests {
		c, err := testConfig()
		if err != nil {
			t.Error(err)
		}

		err = setOutput(&c, test.fileout)
		if !errors.Is(err, test.err) {
			t.Errorf("expected err %v, got %v", test.err, err)
		}
		if c.Pipedout != test.pipedout {
			t.Errorf("expected piped to be %t, got %t", test.pipedout, c.Pipedout)
		}
		if (c.Fileout != nil) != test.file {
			t.Errorf("expected c.Fileout to be %t", test.file)
		}
		if c.Fileout != nil {
			err := c.Fileout.Close()
			if err != nil {
				t.Error(err)
			}
		}
	}
}
