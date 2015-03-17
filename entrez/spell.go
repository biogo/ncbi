// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"github.com/biogo/ncbi/entrez/spell"
)

// <!--
// This is the Current DTD for Entrez eSpell
// $Id:
// -->
// <!-- ================================================================= -->
//
// <!ELEMENT Original		(#PCDATA)>            <!-- \d+ -->
// <!ELEMENT Replaced		(#PCDATA)>            <!-- \d+ -->
//
// <!ELEMENT Database		(#PCDATA)>            <!-- \d+ -->
// <!ELEMENT Query			(#PCDATA)>            <!-- \d+ -->
// <!ELEMENT CorrectedQuery	(#PCDATA)>            <!-- \d+ -->
// <!ELEMENT SpelledQuery	(Replaced|Original)*> <!-- \d+ -->
// <!ELEMENT ERROR			(#PCDATA)>            <!-- \d+ -->
//
// <!ELEMENT eSpellResult	(Database, Query, CorrectedQuery, SpelledQuery, ERROR)>

// All terms listed for eSpell are NOT {\d+}. Interestingly, no blame.

// A Spell holds the deserialised results of an ESpell request.
type Spell struct {
	Database  string             `xml:"Database"`
	Query     string             `xml:"Query"`
	Corrected string             `xml:"CorrectedQuery"`
	Replace   spell.Replacements `xml:"SpelledQuery"`
	Err       string             `xml:"ERROR"`
}
