// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"github.com/biogo/ncbi/entrez/summary"
)

// <!--
// This is the Current DTD for Entrez eSummary version 2
// $Id: eSummary_041029.dtd 49514 2004-10-29 15:52:04Z parantha $
// -->
// <!-- ================================================================= -->
//
// <!ELEMENT Id                (#PCDATA)>          <!-- \d+ -->
//
// <!ELEMENT Item              (#PCDATA|Item)*>   <!-- .+ -->
//
// <!ATTLIST Item
//     Name CDATA #REQUIRED
//     Type (Integer|Date|String|Structure|List|Flags|Qualifier|Enumerator|Unknown) #REQUIRED
// >
//
// <!ELEMENT ERROR             (#PCDATA)>  <!-- .+ -->
//
// <!ELEMENT DocSum            (Id, Item+)>
//
// <!ELEMENT eSummaryResult    (DocSum|ERROR)+>

// A Summary holds the deserialised results of an ESummary request.
type Summary struct {
	Database  string
	Documents []summary.Document `xml:"DocSum"`
	Err       []string           `xml:"ERROR"`
}
