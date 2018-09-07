package numbers

import "testing"

type test struct {
	input    []string
	expected expected
}

type expected struct {
	result string
	ok     bool
}

func TestInts(t *testing.T) {
	tests := []test{
		{[]string{"three"}, expected{"3", true}},
		{[]string{"five"}, expected{"5", true}},
		{[]string{"thirtyfive"}, expected{"35", true}},
		{[]string{"thirty-five"}, expected{"35", true}},
		{[]string{"foo"}, expected{"", false}},
		{[]string{"three", "hundred"}, expected{"300", true}},
		{[]string{"3", "hundred"}, expected{"300", true}},
		{[]string{"+3", "hundred"}, expected{"300", true}},
		{[]string{"-5", "billion"}, expected{"-5000000000", true}},
		{[]string{"3", "hundred", "million"}, expected{"300000000", true}},

		{[]string{"4.58", "hundred"}, expected{"458", true}},
		{[]string{"4.581", "hundred"}, expected{"458.1", true}},
	}

	for _, test := range tests {
		result, ok := Dictionary.Lookup(test.input)

		if ok != test.expected.ok {
			t.Errorf("got ok %v, expected %v, for input %q", ok, test.expected.ok, test.input)
		}
		if result != test.expected.result {
			t.Errorf("got result %q, expected %q, for input %q", result, test.expected.result, test.input)
		}
	}
}
