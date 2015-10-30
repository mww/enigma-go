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
	assertEqualsInt32(t, 2, getOutputIndex(Rotor3(), 1, 'A'-'A', false))
	assertEqualsInt32(t, 3, getOutputIndex(Rotor2(), 0, 2, false))
	assertEqualsInt32(t, 5, getOutputIndex(Rotor1(), 0, 3, false))
	assertEqualsInt32(t, 18, getOutputIndex(ReflectorB(), 0, 5, false))
	assertEqualsInt32(t, 18, getOutputIndex(Rotor1(), 0, 18, true))
	assertEqualsInt32(t, 4, getOutputIndex(Rotor2(), 0, 18, true))
	assertEqualsInt32(t, 1, getOutputIndex(Rotor3(), 1, 4, true))
}

func TestMoveRotor(t *testing.T) {
	m := NewMachine(Rotor1(), Rotor2(), Rotor3(), ReflectorB(), 'A', 'A', 'A')

	assertRotorPositions(t, 0, 0, 0, m)
	m.moveRotors()
	assertRotorPositions(t, 0, 0, 1, m)
	m.moveRotors()
	assertRotorPositions(t, 0, 0, 2, m)
	m.moveRotors()
	assertRotorPositions(t, 0, 0, 3, m)

	m = NewMachine(Rotor1(), Rotor2(), Rotor3(), ReflectorB(), 'Q', 'D', 'U')
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
	m := NewMachine(Rotor1(), Rotor2(), Rotor3(), ReflectorB(), 'A', 'A', 'A')
	assertEqualsRune(t, 'B', m.Step('A'))
	assertEqualsRune(t, 'H', m.Step('P'))
	assertEqualsRune(t, 'S', m.Step('P'))
	assertEqualsRune(t, 'D', m.Step('L'))
	assertEqualsRune(t, 'R', m.Step('E'))

	m = NewMachine(Rotor1(), Rotor2(), Rotor3(), ReflectorB(), 'A', 'D', 'U')
	assertEqualsRune(t, 'W', m.Step('O'))
	assertEqualsRune(t, 'V', m.Step('R'))
	assertEqualsRune(t, 'I', m.Step('A'))
	assertEqualsRune(t, 'D', m.Step('N'))
	assertEqualsRune(t, 'E', m.Step('G'))
	assertEqualsRune(t, 'O', m.Step('E'))

	m = NewMachine(Rotor3(), Rotor1(), Rotor2(), ReflectorB(), 'V', 'P', 'C')
	assertEqualsRune(t, 'K', m.Step('F'))
	assertEqualsRune(t, 'V', m.Step('O'))
	assertEqualsRune(t, 'Z', m.Step('O'))
	assertEqualsRune(t, 'J', m.Step('T'))
	assertEqualsRune(t, 'I', m.Step('B'))
	assertEqualsRune(t, 'T', m.Step('A'))
	assertEqualsRune(t, 'Y', m.Step('L'))
	assertEqualsRune(t, 'W', m.Step('L'))
}
