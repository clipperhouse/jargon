package ascii

// ported from Lucene org.apache.lucene.analysis.miscellaneous;

/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import (
	"bytes"
	"testing"

	"github.com/clipperhouse/jargon"
)

func TestLatin1Accents(t *testing.T) {
	input := []byte(`Des mot clés À LA CHAÎNE À Á Â Ã Ä Å Æ Ç È É Ê Ë Ì Í Î Ï Ĳ Ð Ñ
 Ò Ó Ô Õ Ö Ø Œ Þ Ù Ú Û Ü Ý Ÿ à á â ã ä å æ ç è é ê ë ì í î ï ĳ
 ð ñ ò ó ô õ ö ø œ ß þ ù ú û ü ý ÿ ﬁ ﬂ`)

	expected := []byte(`Des mot cles AA LA CHAINE AA AA AA AA AA AA AE C E E E E I I I I IJ D N
 O O O O O O OE TH U U U U Y Y a a a a a a ae c e e e e i i i i ij
 d n o o o o o o oe ss th u u u u y y fi fl`)

	inputs := bytes.Fields(input)
	expecteds := bytes.Fields(expected)

	if len(inputs) != len(expecteds) {
		t.Error("input and expected should be the same length")
	}

	got, folded := FoldString(input)

	if !folded {
		t.Errorf("input should not have been folded")
	}

	gots := bytes.Fields(got)

	if len(gots) != len(expecteds) {
		t.Error("got and expected should be the same length")
	}

	for i := 0; i < len(inputs); i++ {
		_, folded := FoldString(inputs[i])

		shouldntFold := i < 2 || i == 4
		shouldFold := !shouldntFold

		if shouldntFold && folded {
			t.Errorf("input %s should not have folded", inputs[i])
		}
		if shouldFold && !folded {
			t.Errorf("input %s should have folded", inputs[i])
		}
	}
}

func BenchmarkFold(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tokens := jargon.TokenizeString("Let's go to the café mañana")
		_, err := Fold(tokens).Count()
		if err != nil {
			b.Error(err)
		}
	}
}
