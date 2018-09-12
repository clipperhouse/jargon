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
		{[]string{"three", "hundred"}, expected{"300", true}},
		{[]string{"3", "hundred"}, expected{"300", true}},
		{[]string{"+3", "hundred"}, expected{"300", true}},
		{[]string{"-5", "billion"}, expected{"-5000000000", true}},
		{[]string{"3", "hundred", "million"}, expected{"300000000", true}},

		{[]string{"4.58", "hundred"}, expected{"458", true}},
		{[]string{"4.581", "hundred"}, expected{"458.1", true}},

		{[]string{"foo"}, expected{"", false}},
		{[]string{"foo three"}, expected{"", false}},
		{[]string{"foo 3"}, expected{"", false}},
		{[]string{"hundred"}, expected{"", false}},
		{[]string{"hundred", "3"}, expected{"", false}},
		{[]string{"million", "seven"}, expected{"", false}},
		{[]string{"three", "foo"}, expected{"", false}},
		{[]string{"3", "foo"}, expected{"", false}},
		{[]string{"three", "hundred", "foo"}, expected{"", false}},
		{[]string{"a", "hundred"}, expected{"", false}},
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

func BenchmarkNormalize(b *testing.B) {
	strs := []string{
		"-25",
		"thirty-five",
		"1,000,000",
		"three",
		"seventyseven",
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, s := range strs {
			normalize(s)
		}
	}
}
