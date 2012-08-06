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

import (
	"fmt"
)

type Machine struct {
	// The 3 Rotors and a reflector
	r1, r2, r3, reflector *Rotor

	// The start positions of the 3 Rotors
	s1, s2, s3 rune

	// The current positions of the 3 Rotors, stored as integers in the
	// range 0-25.
	p1, p2, p3 int32
}

func NewMachine(r1, r2, r3, reflector *Rotor, s1, s2, s3 rune) *Machine {
	m := Machine{r1, r2, r3, reflector, s1, s2, s3, 0, 0, 0}
	m.p1 = s1 - 'A'
	m.p2 = s2 - 'A'
	m.p3 = s3 - 'A'
	return &m
}

/*
	Given an input letter, the letter that would be pressed on the keyboard,
	Step() will move the routers and output the resulting letter.
*/
func (m *Machine) Step(input rune) rune {
	m.moveRotors()
	x := getOutputIndex(m.r3, m.p3, input-'A', false)
	x = getOutputIndex(m.r2, m.p2, x, false)
	x = getOutputIndex(m.r1, m.p1, x, false)
	x = getOutputIndex(m.reflector, 0, x, false)
	x = getOutputIndex(m.r1, m.p1, x, true)
	x = getOutputIndex(m.r2, m.p2, x, true)
	x = getOutputIndex(m.r3, m.p3, x, true)
	return LETTERS[x]
}

func (m *Machine) moveRotors() {
	m.p3 = (m.p3 + 1) % 26
	if m.r3.Turnover(m.p3) {
		m.p2 = (m.p2 + 1) % 26
		if m.r2.Turnover(m.p2) {
			m.p1 = (m.p1 + 1) % 26
		}
	}

	// Handles double-stepping case
	if m.r3.Turnover(m.p3-1) && m.r2.Turnover(m.p2+1) {
		m.p2 = (m.p2 + 1) % 26
		if m.r2.Turnover(m.p2) {
			m.p1 = (m.p1 + 1) % 26
		}
	}
}

func (m *Machine) String() string {
	return fmt.Sprintf("KEY: %c%c%c\nROTORS: %s, %s, %s\nREFLECTOR: %s",
		m.s1, m.s2, m.s3, m.r1, m.r2, m.r3, m.reflector)
}

func getOutputIndex(r *Rotor, offset, inputIndex int32, reverse bool) int32 {
	value := (inputIndex + offset) % 26
	letter := LETTERS[value]
	letter = r.Get(letter, reverse)
	value = (letter - 'A') - offset
	if value < 0 {
		value += 26
	}
	return value
}
