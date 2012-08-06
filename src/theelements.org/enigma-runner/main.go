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
	"sort"
	"sync"
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

type result struct {
	message, config string
	diff            float64
}

type results []*result

func (r results) Len() int {
	return len(r)
}

func (r results) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r results) Less(i, j int) bool {
	return r[i].diff < r[j].diff
}

type AtomicInt struct {
	value int32
	lock  *sync.Mutex
}

func NewAtomicInt(value int32) *AtomicInt {
	i := AtomicInt{}
	i.value = value
	i.lock = new(sync.Mutex)
	return &i
}

func (i *AtomicInt) Inc() {
	i.lock.Lock()
	i.value += 1
	i.lock.Unlock()
}

func (i *AtomicInt) Dec() {
	i.lock.Lock()
	i.value -= 1
	i.lock.Unlock()
}

func (i *AtomicInt) Val() int32 {
	return i.value
}

var message = flag.String("message", "", "The encrypted message to crack.")
var numResults = flag.Int("results", 3, "The number of results to display")

func main() {
	flag.Parse()

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

	writer := make(chan *result)
	counter := NewAtomicInt(0)
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
				counter.Inc()
			}
		}
	}

	resultList := make(results, *numResults)
	size := 0
	max := 0.0
	for counter.Val() > 0 {
		r := <-writer
		if size < *numResults {
			resultList[size] = r
			if r.diff > max {
				max = r.diff
			}

			size += 1
			if size == *numResults {
				sort.Sort(resultList)
			}

		} else if r.diff < max {
			resultList[*numResults-1] = r
			sort.Sort(resultList)
		}
		counter.Dec()
	}

	for _, r := range resultList {
		fmt.Printf("%f %s\n%s\n", r.diff, r.message, r.config)
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

func run(m *enigma.Machine, message *string, writer chan *result) {
	var buf bytes.Buffer
	analysis := frequency.NewAnalysis()
	for _, c := range *message {
		l := m.Step(c)
		buf.WriteRune(l)
		analysis.Add(l)
	}
	writer <- &result{buf.String(), m.String(), analysis.Diff()}
}
