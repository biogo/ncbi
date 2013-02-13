// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package spell

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

// An Old string contains the original text of a replacement sequence.
type Old string

func (o Old) String() string { return string(o) }
func (o Old) Type() string   { return "Original" }

// A New string contains the replacement text of a replacement sequence.
type New string

func (r New) String() string { return string(r) }
func (r New) Type() string   { return "Replacement" }
