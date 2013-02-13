// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"code.google.com/p/biogo.entrez/info"
	"code.google.com/p/biogo.entrez/stack"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
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
	DbList []string
	DbInfo []info.DbInfo
	Err    error
}

// Unmarshal fills the fields of an Info from an XML stream read from r.
func (i *Info) Unmarshal(r io.Reader) error {
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
			if t.Name.Local == "DbInfo" {
				var d info.DbInfo
				err := d.Unmarshal(dec, st[len(st)-1:])
				if !reflect.DeepEqual(d, info.DbInfo{}) {
					i.DbInfo = append(i.DbInfo, d)
				}
				if err != nil {
					return err
				}
				st.Drop()
				continue
			}
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "DbName":
				if st.Peek(1) != "DbList" {
					return fmt.Errorf("entrez: unexpected tag: %q", name)
				}
				i.DbList = append(i.DbList, string(t))
			case "ERROR":
				i.Err = errors.New(string(t))
			case "eInfoResult", "DbList", "DbInfo":
			default:
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
