// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"encoding/xml"
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

type Replacement interface {
	String() string
	Type() string
}

type Old string

func (o Old) String() string { return string(o) }
func (o Old) Type() string   { return "Original" }

type New string

func (r New) String() string { return string(r) }
func (r New) Type() string   { return "Replacement" }

// A Spell holds the deserialised results of an ESpell request.
type Spell struct {
	Database  string
	Query     string
	Corrected string
	Replace   []Replacement
	Err       error
}

// Unmarshal fills the fields of a Spell from an XML stream read from r.
func (s *Spell) Unmarshal(r io.Reader) error {
	dec := xml.NewDecoder(r)
	var (
		st  stack
		set bool
	)
	for {
		t, err := dec.Token()
		if err != nil {
			if err != io.EOF {
				return err
			}
			if !st.empty() {
				return io.ErrUnexpectedEOF
			}
			break
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.push(t.Name.Local)
			set = false
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
			case "Database":
				s.Database = string(t)
			case "Query":
				s.Query = string(t)
			case "CorrectedQuery":
				s.Corrected = string(t)
			case "Original", "Replaced":
				if st.peek(1) != "SpelledQuery" {
					return fmt.Errorf("entrez: unexpected tag: %q", name)
				}
				if name == "Original" {
					s.Replace = append(s.Replace, Old(string(t)))
				} else {
					s.Replace = append(s.Replace, New(string(t)))
				}
			case "ERROR":
				s.Err = Error(string(t))
			case "eSpellResult", "SpelledQuery":
			default:
				s.Err = Error(fmt.Sprintf("unknown name: %q", name))
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
			set = true
		case xml.EndElement:
			st, err = st.pair(t.Name.Local)
			if err != nil {
				return err
			}
			if !set {
				switch t.Name.Local {
				case "Original":
					s.Replace = append(s.Replace, Old(""))
				case "Replaced":
					s.Replace = append(s.Replace, New(""))
				}
			}
		}
	}
	return nil
}
