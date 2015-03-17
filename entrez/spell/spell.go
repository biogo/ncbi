// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spell

import (
	"encoding/xml"
	"io"
)

// "github.com/biogo/entrez/xml"

// <!--
// This is the Current DTD for Entrez eSpell
// $Id:
// -->
// <!-- ================================================================= -->
//
// <!ELEMENT Original                     (#PCDATA)>           <!-- \d+ -->
// <!ELEMENT Replaced                     (#PCDATA)>           <!-- \d+ -->
//
// <!ELEMENT Database                     (#PCDATA)>           <!-- \d+ -->
// <!ELEMENT Query                        (#PCDATA)>           <!-- \d+ -->
// <!ELEMENT CorrectedQuery               (#PCDATA)>           <!-- \d+ -->
// <!ELEMENT SpelledQuery                 (Replaced|Original)*> <!-- \d+ -->
// <!ELEMENT ERROR                        (#PCDATA)>           <!-- \d+ -->
//
// <!ELEMENT eSpellResult    (Database, Query, CorrectedQuery, SpelledQuery, ERROR)>

// All terms listed for eSpell are NOT {\d+}. Interestingly, no blame.

// A Replacement is text fragment that indicates a change specified by ESpell.
type Replacement interface {
	String() string
	Type() string
}

// An Old string contains an original segment text of a query.
type Old string

func (o Old) String() string { return string(o) }
func (o Old) Type() string   { return "Original" }

// A New string contains a segment of replaced text of a query.
type New string

func (r New) String() string { return string(r) }
func (r New) Type() string   { return "Replaced" }

type Replacements []Replacement

var _ xml.Unmarshaler = (*Replacements)(nil)

func (r *Replacements) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	*r = (*r)[:0]
	var field string
	for {
		tok, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			field = tok.Name.Local
		case xml.CharData:
			switch field {
			case "Replaced":
				*r = append(*r, New(string(tok)))
			case "Original":
				*r = append(*r, Old(string(tok)))
			}
		case xml.EndElement:
			field = ""
		}
	}
}
