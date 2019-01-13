// Copyright © 2019 Makoto Ito
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package search

import (
	"bytes"
	"io/ioutil"
	"testing"
)

const (
	testNumVector = 4
	testVector    = `apple 1 1 1 1 1
banana 1 1 1 1 1
chocolate 0 0 0 0 0
dragon -1 -1 -1 -1 -1`
)

func TestParseAll(t *testing.T) {
	f := ioutil.NopCloser(bytes.NewReader([]byte(testVector)))
	defer f.Close()
	parser := NewParser(f)
	vectors := make(map[string][]float64)
	storeFunc := func(word string, vec []float64) {
		vectors[word] = vec
	}
	if err := parser.ParseAll(storeFunc); err != nil {
		t.Errorf("Failed to parse vector file: %s", err.Error())
	}

	if len(vectors) != testNumVector {
		t.Errorf("Expected vector len=%d, but got %d", testNumVector, len(vectors))
	}
}

func TestParse(t *testing.T) {
	f := ioutil.NopCloser(bytes.NewReader([]byte(testVector)))
	defer f.Close()
	parser := NewParser(f)
	word, vec, err := parser.parse("apple 1 1 1 1 1")
	if err != nil {
		t.Errorf("Failed to parse a vector str: %s", err.Error())
	}
	if word != "apple" {
		t.Errorf("Expected word=apple, but got %s", word)
	}
	if len(vec) != 5 {
		t.Errorf("Expected vector len=5, but got %d", len(vec))
	}
}
