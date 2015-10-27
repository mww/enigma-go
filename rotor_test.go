/*
 	Copyright 2012 Mark Weaver

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/
package enigma

import "testing"

func TestBuildForward(t *testing.T) {
	expected := map[rune]rune{
		'A': 'E',
		'B': 'K',
		'C': 'M',
		'D': 'F',
		'E': 'L',
		'F': 'G',
		'G': 'D',
		'H': 'Q',
		'I': 'V',
		'J': 'Z',
		'K': 'N',
		'L': 'T',
		'M': 'O',
		'N': 'W',
		'O': 'Y',
		'P': 'H',
		'Q': 'X',
		'R': 'U',
		'S': 'S',
		'T': 'P',
		'U': 'A',
		'V': 'I',
		'W': 'B',
		'X': 'R',
		'Y': 'C',
		'Z': 'J',
	}

	forward := buildForward("EKMFLGDQVZNTOWYHXUSPAIBRCJ")
	for from, expectedResult := range expected {
		if forward[from] != expectedResult {
			t.Errorf("For letter %c, expected %c got %c\n",
				from, expectedResult, forward[from])
		}
	}
}

func TestBuildReverse(t *testing.T) {
	expected := map[rune]rune{
		'A': 'A',
		'J': 'B',
		'D': 'C',
		'K': 'D',
		'S': 'E',
		'I': 'F',
		'R': 'G',
		'U': 'H',
		'X': 'I',
		'B': 'J',
		'L': 'K',
		'H': 'L',
		'W': 'M',
		'T': 'N',
		'M': 'O',
		'C': 'P',
		'Q': 'Q',
		'G': 'R',
		'Z': 'S',
		'N': 'T',
		'P': 'U',
		'Y': 'V',
		'F': 'W',
		'V': 'X',
		'O': 'Y',
		'E': 'Z',
	}

	reverse := buildReverse(buildForward("AJDKSIRUXBLHWTMCQGZNPYFVOE"))
	for from, expectedResult := range expected {
		if reverse[from] != expectedResult {
			t.Errorf("For letter %c, expected %c got %c\n",
				from, expectedResult, reverse[from])
		}
	}
}

func TestGet(t *testing.T) {
	expectedForward := map[rune]rune{
		'A': 'E',
		'G': 'D',
		'M': 'O',
		'T': 'P',
	}

	expectedReverse := map[rune]rune{
		'E': 'A',
		'D': 'G',
		'O': 'M',
		'P': 'T',
	}

	r := NewRotor("short description", "EKMFLGDQVZNTOWYHXUSPAIBRCJ", 'Q')
	for from, expectedResult := range expectedForward {
		if r.Get(from, false) != expectedResult {
			t.Errorf("For letter %c, expected %c got %c\n",
				from, expectedResult, r.Get(from, false))
		}
	}

	for from, expectedResult := range expectedReverse {
		if r.Get(from, true) != expectedResult {
			t.Errorf("For letter %c, expected %c got %c\n",
				from, expectedResult, r.Get(from, true))
		}
	}
}

func TestTurnover(t *testing.T) {
	expected := map[int32]bool{
		0:  false, // A
		25: false, // Z
		16: false, // Q - turnover 1 past stated letter
		17: true,  // R
	}

	r := NewRotor("short description", "EKMFLGDQVZNTOWYHXUSPAIBRCJ", 'Q')
	for input, expectedResult := range expected {
		if r.Turnover(input) != expectedResult {
			t.Errorf("For input %d, expected %t got %t\n",
				input, expectedResult, r.Turnover(input))
		}
	}
}
