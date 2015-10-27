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
	enigma "github.com/mww/enigma-go"
	"github.com/mww/enigma-go/container"
	"github.com/mww/enigma-go/frequency"
	"os"
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

type enigmaResult struct {
	message, config string
	diff            float64
}

var resultFreeList = make(chan *enigmaResult, 2000000)

func populateResultFreeList() {
	for i := 0; i < 1000000; i++ {
		resultFreeList <- new(enigmaResult)
	}
}

func newResult(message, config string, diff float64) *enigmaResult {
	var r *enigmaResult
	select {
	case r = <-resultFreeList:
		r.message, r.config, r.diff = message, config, diff
	default:
		r = new(enigmaResult)
		r.message, r.config, r.diff = message, config, diff
	}
	return r
}

func freeResult(result *enigmaResult) {
	resultFreeList <- result
}

func (r *enigmaResult) Less(other container.Comparer) bool {
	o, ok := other.(*enigmaResult)
	if !ok {
		return false
	}
	// The lower the score the better.
	return r.diff > o.diff
}

var message = flag.String("message", "", "The encrypted message to crack.")
var numResults = flag.Int("results", 3, "The number of results to display")

func main() {
	flag.Parse()

	if len(*message) < 1 {
		fmt.Println("usage: enigma --message=MESSAGE")
		os.Exit(-1)
	}

	results := run(message, *numResults)
	for _, r := range *results {
		fmt.Printf("%f %s\n%s\n", r.diff, r.message, r.config)
	}
}

func run(encryptedMessage *string, numberOfResults int) *[]*enigmaResult {
	populateResultFreeList()
	rotors := []*enigma.Rotor{ROTOR_1, ROTOR_2, ROTOR_3}
	s := make([]interface{}, len(rotors))
	for i, v := range rotors {
		s[i] = v
	}
	rotorPermutations := permutations(s, false)

	reflectors := list.New()
	reflectors.PushBack(REFLECTOR_A)
	reflectors.PushBack(REFLECTOR_B)
	reflectors.PushBack(REFLECTOR_C)

	s = make([]interface{}, len(enigma.LETTERS))
	for i, v := range enigma.LETTERS {
		s[i] = v
	}
	startingPositions := permutations(s, true)

	writer := make(chan *enigmaResult, 250000)
	add_counter := 0
	remove_counter := 0
	for e1 := rotorPermutations.Front(); e1 != nil; e1 = e1.Next() {
		rotors := e1.Value.(*triple)
		r1, r2, r3 := rotors.a.(*enigma.Rotor), rotors.b.(*enigma.Rotor), rotors.c.(*enigma.Rotor)

		for e2 := reflectors.Front(); e2 != nil; e2 = e2.Next() {
			reflector := e2.Value.(*enigma.Rotor)
			for e3 := startingPositions.Front(); e3 != nil; e3 = e3.Next() {
				pos := e3.Value.(*triple)
				p1, p2, p3 := pos.a.(rune), pos.b.(rune), pos.c.(rune)

				m := enigma.NewMachine(r1, r2, r3, reflector, p1, p2, p3)
				go runMachine(m, encryptedMessage, writer)
				add_counter++
			}
		}
	}

	resultList := container.NewSortedFixedSizeList(numberOfResults)
	for remove_counter < add_counter {
		r := <-writer
		if !resultList.MaybeAdd(r) {
			freeResult(r)
		}
		remove_counter++
	}

	results := make([]*enigmaResult, numberOfResults)
	itr := resultList.Iterator()
	for i := 0; itr.HasNext(); i++ {
		r := itr.Next().(*enigmaResult)
		results[i] = r
	}
	return &results
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

func runMachine(m *enigma.Machine, message *string, writer chan *enigmaResult) {
	var buf bytes.Buffer
	analysis := frequency.NewAnalysis()
	for _, c := range *message {
		l := m.Step(c)
		buf.WriteRune(l)
		analysis.Add(l)
	}
	writer <- newResult(buf.String(), m.String(), analysis.Diff())
	enigma.FreeMachine(m)
}
