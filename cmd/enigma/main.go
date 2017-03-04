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
	"flag"
	"fmt"
	"os"

	enigma "github.com/mww/enigma-go"
	"github.com/mww/enigma-go/container"
	"github.com/mww/enigma-go/frequency"
)

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
		r = &enigmaResult{
			message: message,
			config:  config,
			diff:    diff,
		}
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
	rotors := []*enigma.Rotor{enigma.Rotor1(), enigma.Rotor2(), enigma.Rotor3()}
	s := make([]interface{}, len(rotors))
	for i, v := range rotors {
		s[i] = v
	}
	rotorPermutations := permutations(s, false)

	reflectors := []*enigma.Rotor{enigma.ReflectorA(), enigma.ReflectorB(), enigma.ReflectorC()}

	s = make([]interface{}, len(enigma.LETTERS))
	for i, v := range enigma.LETTERS {
		s[i] = v
	}
	startingPositions := permutations(s, true)

	writer := make(chan *enigmaResult, 250000)
	addCounter := 0
	removeCounter := 0
	for _, rotors := range rotorPermutations {
		r1, r2, r3 := rotors.a.(*enigma.Rotor), rotors.b.(*enigma.Rotor), rotors.c.(*enigma.Rotor)

		for _, reflector := range reflectors {
			for _, pos := range startingPositions {
				p1, p2, p3 := pos.a.(rune), pos.b.(rune), pos.c.(rune)

				m := enigma.NewMachine(r1, r2, r3, reflector, p1, p2, p3)
				go runMachine(m, encryptedMessage, writer)
				addCounter++
			}
		}
	}

	resultList := container.NewSortedFixedSizeList(numberOfResults)
	for removeCounter < addCounter {
		r := <-writer
		if !resultList.MaybeAdd(r) {
			freeResult(r)
		}
		removeCounter++
	}

	results := make([]*enigmaResult, numberOfResults)
	itr := resultList.Iterator()
	for i := 0; itr.HasNext(); i++ {
		r := itr.Next().(*enigmaResult)
		results[i] = r
	}
	return &results
}

// Calculate all of the N choose 3 permutations.
func permutations(items []interface{}, duplicates bool) []*triple {
	// TODO(mww): Calculate this dynamically
	var size int
	switch len(items) {
	case 3:
		size = 6
	case 4:
		size = 24
	case 5:
		size = 60
	default:
		size = 17576 // a large value for 26 choose 3
	}
	// Make a large buffer so we don't have to resize it much. 16000 will be big
	// enough for 26 choose 3
	result := make([]*triple, 0, size)
	for _, a := range items {
		for _, b := range items {
			if !duplicates && a == b {
				continue
			}

			for _, c := range items {
				if !duplicates && (a == b || a == c || b == c) {
					continue
				}

				result = append(result, &triple{a, b, c})
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
