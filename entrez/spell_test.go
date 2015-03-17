// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"encoding/xml"
	"strings"

	. "github.com/biogo/ncbi/entrez/spell"

	"gopkg.in/check.v1"
)

func (s *S) TestParseSpell(c *check.C) {
	for i, t := range []struct {
		retval string
		spell  Spell
	}{
		{`<?xml version="1.0"?>
<!DOCTYPE eSpellResult PUBLIC "-//NLM//DTD eSpellResult, 23 November 2004//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eSpell.dtd">
<eSpellResult>
	<Database>pubmed</Database>
	<Query>asthmaa OR alergies</Query>
	<CorrectedQuery>asthma or allergies</CorrectedQuery>
	<SpelledQuery><Original></Original><Replaced>asthma</Replaced><Original> OR </Original><Replaced>allergies</Replaced></SpelledQuery>
	<ERROR/>
</eSpellResult>
`,
			Spell{
				Database:  "pubmed",
				Query:     "asthmaa OR alergies",
				Corrected: "asthma or allergies",
				Replace: []Replacement{
					New("asthma"),
					Old(" OR "),
					New("allergies"),
				},
				Err: "",
			},
		},
		{
			`<?xml version="1.0"?>
<!DOCTYPE eSpellResult PUBLIC "-//NLM//DTD eSpellResult, 23 November 2004//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eSpell.dtd">
<eSpellResult>
	<Database>pubmed</Database>
	<Query>asthmaa OR alergies</Query>
	<CorrectedQuery>asthma or allergies</CorrectedQuery>
	<SpelledQuery>
		<Original></Original>
		<Replaced>asthma</Replaced>
		<Original> OR </Original>
		<Replaced>allergies</Replaced>
	</SpelledQuery>
	<ERROR/>
</eSpellResult>
`,
			Spell{
				Database:  "pubmed",
				Query:     "asthmaa OR alergies",
				Corrected: "asthma or allergies",
				Replace: []Replacement{
					New("asthma"),
					Old(" OR "),
					New("allergies"),
				},
				Err: "",
			},
		},
	} {
		var sp Spell
		err := xml.NewDecoder(strings.NewReader(t.retval)).Decode(&sp)
		c.Check(err, check.Equals, nil, check.Commentf("Test: %d", i))
		c.Check(sp, check.DeepEquals, t.spell, check.Commentf("Test: %d", i))
	}
}
