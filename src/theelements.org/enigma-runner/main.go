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

package main

import (
	"bytes"
	"container/list"
	"flag"
	"fmt"
	"theelements.org/enigma"
	"theelements.org/frequency"
)

var ROTOR_1 = enigma.NewRotor("Rotor 1, 1930", "EKMFLGDQVZNTOWYHXUSPAIBRCJ", 'Q')
var ROTOR_2 = enigma.NewRotor("Rotor 2, 1930", "AJDKSIRUXBLHWTMCQGZNPYFVOE", 'E')
var ROTOR_3 = enigma.NewRotor("Rotor 3, 1930", "BDFHJLCPRTXVZNYEIWGAKMUSQO", 'V')
var ROTOR_4 = enigma.NewRotor("Rotor 4, 1938", "ESOVPZJAYQUIRHXLNFTGKDCMWB", 'J')
var ROTOR_5 = enigma.NewRotor("Rotor 5, 1938", "VZBRGITYUPSDNHLXAWMJQOFECK", 'Z')

var REFLECTOR_A = enigma.NewRotor("Reflector A", "EJMZALYXVBWFCRQUONTSPIKHGD", 'Z')
var REFLECTOR_B = enigma.NewRotor("Reflector B", "YRUHQSLDPXNGOKMIEBFZCWVJAT", 'Z')
var REFLECTOR_C = enigma.NewRotor("Reflector C", "FVPJIAOYEDRZXWGCTKUQSBNMHL", 'Z')

/*
	Simple container to hold configuration data.
*/
type triple struct {
	a, b, c interface{}
}

// XXX remove this function
func (t *triple) String() string {
	return fmt.Sprintf("I'm a triple %s,%s,%s", t.a, t.b, t.c)
}

type enigmaResult struct {
	message string
	diff    float64
}

var message = flag.String("message", "", "The encrypted message to crack.")

func main() {
	flag.Parse()

	rotors := []*enigma.Rotor{ROTOR_1, ROTOR_2, ROTOR_3}
	s := make([]interface{}, len(rotors))
	for i, v := range rotors {
		s[i] = v
	}
	rotorPermutations := permutations(s, false)

	reflectors := list.New()
	//reflectors.PushBack(REFLECTOR_A)
	reflectors.PushBack(REFLECTOR_B)
	//reflectors.PushBack(REFLECTOR_C)

	s = make([]interface{}, len(enigma.LETTERS))
	for i, v := range enigma.LETTERS {
		s[i] = v
	}
	startingPositions := permutations(s, true)

	writer := make(chan *enigmaResult)
	for e1 := rotorPermutations.Front(); e1 != nil; e1 = e1.Next() {
		rotors := e1.Value.(*triple)
		r1, r2, r3 := rotors.a.(*enigma.Rotor), rotors.b.(*enigma.Rotor), rotors.c.(*enigma.Rotor)
		for e2 := reflectors.Front(); e2 != nil; e2 = e2.Next() {
			reflector := e2.Value.(*enigma.Rotor)
			for e3 := startingPositions.Front(); e3 != nil; e3 = e3.Next() {
				pos := e3.Value.(*triple)
				p1, p2, p3 := pos.a.(rune), pos.b.(rune), pos.c.(rune)

				m := enigma.NewMachine(r1, r2, r3, reflector, p1, p2, p3)
				go run(m, message, writer)
			}
		}
	}

	for {
		r := <-writer
		fmt.Printf("%s %f\n", r.message, r.diff)
	}
}

func permutations(items []interface{}, duplicates bool) *list.List {
	result := list.New()
	for _, a := range items {
		for _, b := range items {
			if !duplicates && a == b {
				continue
			}

			for _, c := range items {
				if !duplicates && (a == b || a == c || b == c) {
					continue
				}

				result.PushBack(&triple{a, b, c})
			}
		}
	}
	return result
}

func run(m *enigma.Machine, message *string, writer chan *enigmaResult) {
	var buf bytes.Buffer
	analysis := frequency.NewAnalysis()
	for _, c := range *message {
		l := m.Step(c)
		buf.WriteRune(l)
		analysis.Add(l)
	}
	result := enigmaResult{buf.String(), analysis.Diff()}
	writer <- &result
}
