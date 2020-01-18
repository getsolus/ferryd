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

package jobs

// List is a list of Jobs
type List []*Job

// Len returns the length of the list
func (l List) Len() int {
	return len(l)
}

// Less compares tow jobs by creation time
func (l List) Less(i, j int) bool {
	return l[i].Created.Time.Before(l[j].Created.Time)
}

// Swap switches two Jobs for sorting
func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
