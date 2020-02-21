// Package ascii folds Unicode characters to their ASCII equivalents where possible.
// Ported from Lucene org.apache.lucene.analysis.miscellaneous
package ascii

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

// Fold is a jargon.TokenFilter, which converts alphabetic, numeric, and symbolic Unicode characters
// which are not in the first 127 ASCII characters (the "Basic Latin" Unicode
// block) into their ASCII equivalents, if one exists.
// Ported from Lucene org.apache.lucene.analysis.miscellaneous
var Fold = &filter{}

type filter struct{}

func (f *filter) Lookup(s []string) (string, bool) {
	word := s[0]
	return fold(word)
}

func (f *filter) MaxGramLength() int {
	return 1
}