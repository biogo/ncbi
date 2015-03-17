// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package summary

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

type Item struct {
	Value string `xml:",chardata"`
	Name  string `xml:",attr"`
	Type  string `xml:",attr"`
}

type Document struct {
	Id    int    `xml:"Id"`
	Items []Item `xml:"Item"`
}
