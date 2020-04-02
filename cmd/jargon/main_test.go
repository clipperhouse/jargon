package main

import (
	"os"
	"testing"
)

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
			filein: "main.go",
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
			filein: "main.go",
			mode:   os.ModeAppend,

			err:     errTwoInput,
			pipedin: true,
			file:    true,
		},
	}

	for _, test := range tests {
		c := config{}
		err := setInput(&c, test.mode, test.filein)
		if err != test.err {
			t.Errorf("expected err %v, got %v", test.err, err)
		}
		if c.Pipedin != test.pipedin {
			t.Errorf("expected piped to be %t, got %t", test.pipedin, c.Pipedin)
		}
		if (c.Filein != nil) != test.file {
			t.Errorf("expected c.Filein to be %t", test.file)
		}
	}
}
