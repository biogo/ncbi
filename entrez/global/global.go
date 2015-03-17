// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package global

// <!--
//         This is the Current DTD for Entrez eGSearch
//         $Id: egquery.dtd 39250 2004-05-03 16:19:48Z yasmax $
// -->
// <!-- ================================================================= -->
//
// <!ELEMENT       DbName          (#PCDATA)>      <!-- .+ -->
// <!ELEMENT       MenuName        (#PCDATA)>      <!-- .+ -->
// <!ELEMENT       Count           (#PCDATA)>      <!-- \d+ -->
// <!ELEMENT       Status          (#PCDATA)>      <!-- .+ -->
// <!ELEMENT       Term            (#PCDATA)>      <!-- .+ -->
//
// <!ELEMENT       ResultItem      (
//                                      DbName,
//                                      MenuName,
//                                      Count,
//                                      Status
//                                 )>
// <!ELEMENT       eGQueryResult  (ResultItem+)>
//
// <!ELEMENT       Result         (Term, eGQueryResult)>

type Result struct {
	Database string `xml:"DbName"`
	MenuName string `xml:"MenuName"`
	Count    int    `xml:"Count"`
	Status   string `xml:"Status"`
}
