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

import "testing"

func TestEnigmaRunner(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full test in short mode.")
	}

	encrypted := "ZTQBLVXKPBPGAVQBRYDYQEZNKRLMZTMRGBJSQKHDPHHNTNIDLYVFCOKZYYSMJFAHQBTEAVFKOXRPSQX"
	expected := "THISISASLIGHTLYLONGERTESTSOIHAVETOSEEIFICANKEEPWRITINGALONGERSTRINGTOUSEASINPUT"
	results := run(&encrypted, 3)
	if (*results)[0].message != expected {
		t.Errorf("Expected %s, got %s", expected, (*results)[0].message)
	}
}
