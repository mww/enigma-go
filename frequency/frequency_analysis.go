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
package frequency

import "math"

//import "fmt"

var englishExpectedFrequency = map[rune]float64{
	'A': 8.167,
	'B': 1.492,
	'C': 2.782,
	'D': 4.253,
	'E': 12.702,
	'F': 2.228,
	'G': 2.015,
	'H': 6.094,
	'I': 6.966,
	'J': 0.153,
	'K': 0.772,
	'L': 4.025,
	'M': 2.406,
	'N': 6.749,
	'O': 7.507,
	'P': 1.929,
	'Q': 0.095,
	'R': 5.987,
	'S': 6.327,
	'T': 9.056,
	'U': 2.758,
	'V': 0.978,
	'W': 2.360,
	'X': 0.150,
	'Y': 1.974,
	'Z': 0.074,
}

type Analysis struct {
	characters map[rune]float64
	total      float64
}

func NewAnalysis() *Analysis {
	a := Analysis{nil, 0}

	a.characters = make(map[rune]float64)
	for k, _ := range englishExpectedFrequency {
		a.characters[k] = 0.0
	}

	return &a
}

func (a *Analysis) Add(c rune) {
	a.characters[c] = a.characters[c] + 1
	a.total += 1
	//fmt.Printf("%c %f: total: %f\n", c, a.characters[c], a.total)
}

func (a *Analysis) Diff() float64 {
	diff := 0.0
	for k, _ := range englishExpectedFrequency {
		count := a.characters[k]
		if count == 0 {
			continue
		}
		actual := (count / a.total) * 100
		expected := englishExpectedFrequency[k]
		diff += math.Abs(expected - actual)
		//fmt.Printf("%c actual: %f expected: %f diff: %f\n", k, actual, expected, diff)
	}
	//fmt.Printf("Returning diff of: %f\n", diff)
	return diff
}
