//
// Copyright © 2017-2020 Solus Project
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package archive

// Archives represents a ist of archives belonging to specific Repo
type Archives []Archive

// Len returns the number of Archives in the list
func (as Archives) Len() int {
	return len(as)
}

// Less checks if one Archive is less than another for sorting
func (as Archives) Less(i, j int) bool {
	return as[i].Compare(as[j]) < 0
}

// Swap exchanges two entries in this list
func (as Archives) Swap(i, j int) {
	as[i], as[j] = as[j], as[i]
}

// Diff calculates the difference between two Archive lists
func (as Archives) Diff(others Archives) (diff Archives) {
	for _, a1 := range as {
		found := false
		for _, a2 := range others {
			if a1.ID == a2.ID {
				diff = append(diff, a1.Copy())
				found = true
				break
			}
		}
		if !found {
			a := a1.Copy()
			a.Status = StatusAdded
			diff = append(diff, a)
		}
	}
	for _, a2 := range others {
		found := false
		for _, a1 := range as {
			if a1.ID == a2.ID {
				found = true
				break
			}
		}
		if !found {
			a := a2.Copy()
			a.Status = StatusRemoved
			diff = append(diff, a)
		}
	}
	return
}
