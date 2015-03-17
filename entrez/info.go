// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"github.com/biogo/ncbi/entrez/info"
)

// <!--
//                 This is the Current DTD for Entrez eInfo
// $Id: eInfo_020511.dtd 361872 2012-05-04 17:46:41Z fialkov $
// -->
// <!-- ================================================================= -->
//
// <!ELEMENT	DbName		(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	Name		(#PCDATA)>	<!-- .+ -->
// <!ELEMENT	FullName	(#PCDATA)>	<!-- .+ -->
// <!ELEMENT	Description	(#PCDATA)>	<!-- .+ -->
// <!ELEMENT	TermCount	(#PCDATA)>	<!-- \d+ -->
// <!ELEMENT	Menu		(#PCDATA)>	<!-- .+ -->
// <!ELEMENT	DbTo		(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	MenuName	(#PCDATA)>	<!-- .+ -->
// <!ELEMENT	Count		(#PCDATA)>	<!-- \d+ -->
// <!ELEMENT	LastUpdate	(#PCDATA)>	<!-- \d+ -->
//
// <!ELEMENT	ERROR		(#PCDATA)>	<!-- .+ -->
//
// <!ELEMENT	IsDate		(#PCDATA)>	<!-- (Y|N) -->
// <!ELEMENT	IsNumerical	(#PCDATA)>	<!-- (Y|N) -->
// <!ELEMENT	SingleToken	(#PCDATA)>	<!-- (Y|N) -->
// <!ELEMENT	Hierarchy	(#PCDATA)>	<!-- (Y|N) -->
// <!ELEMENT	IsHidden	(#PCDATA)>	<!-- (Y|N) -->
// <!ELEMENT    IsRangable      (#PCDATA)>      <!-- (Y|N) -->
// <!ELEMENT    IsTruncatable   (#PCDATA)>      <!-- (Y|N) -->
//
//
// <!ELEMENT	DbList		(DbName+)>
//
// <!ELEMENT	Field		(Name,
//              FullName,
// 				Description,
// 				TermCount,
// 				IsDate,
// 				IsNumerical,
// 				SingleToken,
// 				Hierarchy,
//              IsHidden,
//              IsRangable,
//              IsTruncatable)>
//
// <!ELEMENT	Link		(Name,Menu,Description,DbTo)>
//
//
// <!ELEMENT	LinkList	(Link*)>
// <!ELEMENT	FieldList	(Field*)>
//
//
// <!ELEMENT	DbInfo		(DbName,
// 				MenuName,
// 				Description,
// 				Count,
// 				LastUpdate,
// 				FieldList,
// 				LinkList?)>
//
// <!ELEMENT	eInfoResult	(DbList|DbInfo|ERROR)>

// An Info holds the deserialised results of an EInfo request.
type Info struct {
	DbList []string     `xml:"DbList>DbName"`
	DbInfo *info.DbInfo `xml:"DbInfo"`
	Err    string       `xml:"ERROR"`
}
