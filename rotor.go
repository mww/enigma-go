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

var LETTERS = []rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K',
	'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

type Rotor struct {
	description      string
	turnover         int32
	forward, reverse map[rune]rune
}

func NewRotor(description, mapping string, turnoverLetter rune) *Rotor {
	r := Rotor{description: description}
	r.turnover = (turnoverLetter - 'A') + 1
	r.forward = buildForward(mapping)
	r.reverse = buildReverse(r.forward)
	return &r
}

func Rotor1() *Rotor {
	return NewRotor("Rotor 1, 1930", "EKMFLGDQVZNTOWYHXUSPAIBRCJ", 'Q')
}

func Rotor2() *Rotor {
	return NewRotor("Rotor 2, 1930", "AJDKSIRUXBLHWTMCQGZNPYFVOE", 'E')
}

func Rotor3() *Rotor {
	return NewRotor("Rotor 3, 1930", "BDFHJLCPRTXVZNYEIWGAKMUSQO", 'V')
}

func Rotor4() *Rotor {
	return NewRotor("Rotor 4, 1938", "ESOVPZJAYQUIRHXLNFTGKDCMWB", 'J')
}

func Rotor5() *Rotor {
	return NewRotor("Rotor 5, 1938", "VZBRGITYUPSDNHLXAWMJQOFECK", 'Z')
}

func ReflectorA() *Rotor {
	return NewRotor("Reflector A", "EJMZALYXVBWFCRQUONTSPIKHGD", 'Z')
}

func ReflectorB() *Rotor {
	return NewRotor("Reflector B", "YRUHQSLDPXNGOKMIEBFZCWVJAT", 'Z')
}

func ReflectorC() *Rotor {
	return NewRotor("Reflector C", "FVPJIAOYEDRZXWGCTKUQSBNMHL", 'Z')
}

func buildForward(mapping string) map[rune]rune {
	m := make(map[rune]rune)
	for i, char := range mapping {
		m[LETTERS[i]] = char
	}

	return m
}

func buildReverse(forward map[rune]rune) map[rune]rune {
	m := make(map[rune]rune)
	for from, to := range forward {
		m[to] = from
	}
	return m
}

func (r *Rotor) Get(letter rune, reverse bool) rune {
	if reverse {
		return r.reverse[letter]
	}
	return r.forward[letter]
}

func (r *Rotor) Turnover(position int32) bool {
	return position == r.turnover
}

func (r *Rotor) String() string {
	return r.description
}
