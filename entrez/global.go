// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"github.com/biogo/ncbi/entrez/global"
)

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

// A Global holds the deserialised results of an EGQuery request.
type Global struct {
	Query   string          `xml:"Term"`
	Results []global.Result `xml:"eGQueryResult>ResultItem"`
}
