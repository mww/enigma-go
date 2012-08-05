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

func createTestMachine(p1, p2, p3 rune) *Machine {
	r1 := NewRotor("Rotor 1, 1930", "EKMFLGDQVZNTOWYHXUSPAIBRCJ", 'Q')
	r2 := NewRotor("Rotor 2, 1930", "AJDKSIRUXBLHWTMCQGZNPYFVOE", 'E')
	r3 := NewRotor("Rotor 3, 1930", "BDFHJLCPRTXVZNYEIWGAKMUSQO", 'V')
	reflector := NewRotor("Reflector B", "YRUHQSLDPXNGOKMIEBFZCWVJAT", 'Z')
	return NewMachine(r1, r2, r3, reflector, p1, p2, p3)
}

func assertEqualsInt32(t *testing.T, expected, actual int32) {
	if expected != actual {
		t.Errorf("Expected %d, got %d\n", expected, actual)
	}
}

func assertEqualsRune(t *testing.T, expected, actual rune) {
	if expected != actual {
		t.Errorf("Expected %c, got %c\n", expected, actual)
	}
}

func assertRotorPositions(t *testing.T, e1, e2, e3 int32, m *Machine) {
	if e1 != m.p1 || e2 != m.p2 || e3 != m.p3 {
		t.Errorf("Expected %d,%d,%d, got %d,%d,%d", e1, e2, e3, m.p1, m.p2, m.p3)
	}
}

func TestGetOutputIndex(t *testing.T) {
	r1 := NewRotor("Rotor 1, 1930", "EKMFLGDQVZNTOWYHXUSPAIBRCJ", 'Q')
	r2 := NewRotor("Rotor 2, 1930", "AJDKSIRUXBLHWTMCQGZNPYFVOE", 'E')
	r3 := NewRotor("Rotor 3, 1930", "BDFHJLCPRTXVZNYEIWGAKMUSQO", 'V')
	reflector := NewRotor("Reflector B", "YRUHQSLDPXNGOKMIEBFZCWVJAT", 'Z')

	assertEqualsInt32(t, 2, getOutputIndex(r3, 1, 'A'-'A', false))
	assertEqualsInt32(t, 3, getOutputIndex(r2, 0, 2, false))
	assertEqualsInt32(t, 5, getOutputIndex(r1, 0, 3, false))
	assertEqualsInt32(t, 18, getOutputIndex(reflector, 0, 5, false))
	assertEqualsInt32(t, 18, getOutputIndex(r1, 0, 18, true))
	assertEqualsInt32(t, 4, getOutputIndex(r2, 0, 18, true))
	assertEqualsInt32(t, 1, getOutputIndex(r3, 1, 4, true))
}

func TestMoveRotor(t *testing.T) {
	m := createTestMachine('A', 'A', 'A')

	assertRotorPositions(t, 0, 0, 0, m)
	m.moveRotors()
	assertRotorPositions(t, 0, 0, 1, m)
	m.moveRotors()
	assertRotorPositions(t, 0, 0, 2, m)
	m.moveRotors()
	assertRotorPositions(t, 0, 0, 3, m)

	m = createTestMachine('Q', 'D', 'U')
	assertRotorPositions(t, 16, 3, 20, m)
	m.moveRotors()
	assertRotorPositions(t, 16, 3, 21, m)
	m.moveRotors()
	assertRotorPositions(t, 16, 4, 22, m)
	m.moveRotors()
	assertRotorPositions(t, 17, 5, 23, m)
	m.moveRotors()
	assertRotorPositions(t, 17, 5, 24, m)
	m.moveRotors()
	assertRotorPositions(t, 17, 5, 25, m)
	m.moveRotors()
	assertRotorPositions(t, 17, 5, 0, m)
}

func TestStep(t *testing.T) {
	m := createTestMachine('A', 'A', 'A')
	assertEqualsRune(t, 'B', m.Step('A'))
	assertEqualsRune(t, 'H', m.Step('P'))
	assertEqualsRune(t, 'S', m.Step('P'))
	assertEqualsRune(t, 'D', m.Step('L'))
	assertEqualsRune(t, 'R', m.Step('E'))

	m = createTestMachine('Q', 'D', 'U')
	assertEqualsRune(t, 'C', m.Step('O'))
	assertEqualsRune(t, 'X', m.Step('R'))
	assertEqualsRune(t, 'N', m.Step('A'))
	assertEqualsRune(t, 'E', m.Step('N'))
	assertEqualsRune(t, 'M', m.Step('G'))
	assertEqualsRune(t, 'W', m.Step('E'))
}
