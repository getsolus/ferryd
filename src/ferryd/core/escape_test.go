package core_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/getsolus/ferryd/src/ferryd/core"
)

func TestEscapedWriter_Write(t *testing.T) {
	tests := map[string]struct {
		In, Out string
	}{
		"lato": {
			In:  `Lato is a sanserif type­face fam­ily designed in the Sum­mer 2010 by Warsaw-​​based designer Łukasz Dziedzic (“Lato” means “Sum­mer” in Pol­ish).`,
			Out: `Lato is a sanserif type&#xAD;face fam&#xAD;ily designed in the Sum&#xAD;mer 2010 by Warsaw-&#x200B;&#x200B;based designer Łukasz Dziedzic (“Lato” means “Sum&#xAD;mer” in Pol&#xAD;ish).`,
		},
		"translate-shell": {
			In:  ` py   master  /  repo  translate-shell  trans 'est-ce que ca marche?'`,
			Out: `&#xA0;py&#xA0;&#xE0B0;&#xA0;&#xE0A0;&#xA0;master&#xA0;&#xE0B0;&#xA0;/&#xA0;&#xE0B1;&#xA0;repo&#xA0;&#xE0B0;&#xA0;translate-shell&#xA0;&#xE0B0;&#xA0;trans 'est-ce que ca marche?'`,
		},
		"invalid": {
			In:  "Invalid: \xa0",
			Out: "Invalid: \ufffd",
		},
		"emoji": {
			In:  `If you want to “🙂” whenever you type “:-)”`,
			Out: `If you want to “&#x1F642;” whenever you type “:-)”`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			var out bytes.Buffer

			_, err := core.NewEscapedWriter(&out).Write([]byte(test.In))
			require.NoError(t, err)
			require.Equal(t, test.Out, out.String())
		})
	}
}
