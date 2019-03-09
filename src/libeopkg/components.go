//
// Copyright © 2017-2019 Solus Project
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

package libeopkg

import (
	"encoding/xml"
	"os"
	"sort"
)

// A Component as seen through the eyes of XML
type Component struct {
	Name string // ID of this component, i.e. "system.base"

	// Translated short name
	LocalName []LocalisedField

	// Translated summary
	Summary []LocalisedField

	// Translated description
	Description []LocalisedField

	Group      string // Which group this component belongs to
	Maintainer struct {
		Name  string // Name of the component maintainer
		Email string // Contact e-mail address of component maintainer
	}
}

// Components is a simple helper wrapper for loading from components.xml files
type Components struct {
	Components []Component `xml:"Components>Component"`
}

// ComponentSort allows us to quickly sort our components by name
type ComponentSort []Component

func (g ComponentSort) Len() int {
	return len(g)
}

func (g ComponentSort) Less(a, b int) bool {
	return g[a].Name < g[b].Name
}

func (g ComponentSort) Swap(a, b int) {
	g[a], g[b] = g[b], g[a]
}

// NewComponents will load the Components data from the XML file
func NewComponents(xmlfile string) (cs *Components, err error) {
	cFile, err := os.Open(xmlfile)
	if err != nil {
		return
	}
	defer cFile.Close()
	cs = &Components{}
	dec := xml.NewDecoder(cFile)
	if err = dec.Decode(cs); err != nil {
		return
	}
	// Sort components by name
	sort.Sort(ComponentSort(cs.Components))

	// Ensure there are no empty Lang= fields
	for i := range cs.Components {
		comp := &cs.Components[i]
		FixMissingLocalLanguage(&comp.LocalName)
		FixMissingLocalLanguage(&comp.Summary)
		FixMissingLocalLanguage(&comp.Description)
	}
	return
}
