// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"code.google.com/p/biogo.entrez/spell"
	"code.google.com/p/biogo.entrez/stack"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
)

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

// A Spell holds the deserialised results of an ESpell request.
type Spell struct {
	Database  string
	Query     string
	Corrected string
	Replace   []spell.Replacement
	Err       error
}

// Unmarshal fills the fields of a Spell from an XML stream read from r.
func (s *Spell) Unmarshal(r io.Reader) error {
	dec := xml.NewDecoder(r)
	var st stack.Stack
	for {
		t, err := dec.Token()
		if err != nil {
			if err != io.EOF {
				return err
			}
			if !st.Empty() {
				return io.ErrUnexpectedEOF
			}
			break
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.Push(t.Name.Local)
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "Database":
				s.Database = string(t)
			case "Query":
				s.Query = string(t)
			case "CorrectedQuery":
				s.Corrected = string(t)
			case "Original", "Replaced":
				if st.Peek(1) != "SpelledQuery" {
					return fmt.Errorf("entrez: unexpected tag: %q", name)
				}
				if name == "Original" {
					s.Replace = append(s.Replace, spell.Old(string(t)))
				} else {
					s.Replace = append(s.Replace, spell.New(string(t)))
				}
			case "ERROR":
				s.Err = errors.New(string(t))
			case "eSpellResult", "SpelledQuery":
			default:
				s.Err = fmt.Errorf("unknown name: %q", name)
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
