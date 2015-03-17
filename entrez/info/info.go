// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package info

import (
	"encoding/xml"
	"errors"
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

type Field struct {
	Name          string `xml:"Name"`
	FullName      string `xml:"FullName"`
	Description   string `xml:"Description"`
	TermCount     int    `xml:"TermCount"`
	IsDate        Bool   `xml:"IsData"`
	IsNumerical   Bool   `xml:"IsNumerical"`
	SingleToken   Bool   `xml:"SingleToken"`
	Hierarchy     Bool   `xml:"Hierarchy"`
	IsHidden      Bool   `xml:"IsHidden"`
	IsRangeable   Bool   `xml:"IsRangable"`
	IsTruncatable Bool   `xml:"IsTruncatable"`
}

type DbLink struct {
	Name        string `xml:"Name"`
	FullName    string `xml:"FullName"`
	Description string `xml:"Description"`
	DbTo        string `xml:"DbTo"`
}

type DbInfo struct {
	DbName      string   `xml:"DbName"`
	MenuName    string   `xml:"MenuName"`
	Description string   `xml:"Description"`
	Count       int      `xml:"Count"`
	LastUpdate  string   `xml:"LastUpdate"`
	FieldList   []Field  `xml:"FieldList>Field"`
	LinkList    []DbLink `xml:"LinkList>Link"`
}

type Bool bool

var _ xml.Unmarshaler = (*Bool)(nil)

func (t *Bool) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	var c string
	err := dec.DecodeElement(&c, &start)
	if err != nil {
		return err
	}
	switch c {
	case "Y":
		*t = true
	case "N":
		*t = false
	default:
		return errors.New("entrez: bad boolean")
	}
	return nil
}
