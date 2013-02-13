// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package info

import (
	"code.google.com/p/biogo.entrez/stack"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
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
	Name          string
	FullName      string
	Description   string
	TermCount     int
	IsDate        bool
	IsNumerical   bool
	SingleToken   bool
	Hierarchy     bool
	IsHidden      bool
	IsRangeable   bool
	IsTruncatable bool
}

func (f *Field) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
	for {
		t, err := dec.Token()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.Push(t.Name.Local)
		case xml.CharData:
			switch name := st.Peek(0); name {
			case "Name":
				f.Name = string(t)
			case "FullName":
				f.FullName = string(t)
			case "Description":
				f.Description = string(t)
			case "TermCount":
				c, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				f.TermCount = c
			case "IsDate":
				switch b := string(t); b {
				case "Y", "N":
					f.IsDate = b == "Y"
				default:
					return fmt.Errorf("eutil: bad boolean: %q", b)
				}
			case "IsNumerical":
				switch b := string(t); b {
				case "Y", "N":
					f.IsNumerical = b == "Y"
				default:
					return fmt.Errorf("eutil: bad boolean: %q", b)
				}
			case "SingleToken":
				switch b := string(t); b {
				case "Y", "N":
					f.SingleToken = b == "Y"
				default:
					return fmt.Errorf("eutil: bad boolean: %q", b)
				}
			case "Hierarchy":
				switch b := string(t); b {
				case "Y", "N":
					f.Hierarchy = b == "Y"
				default:
					return fmt.Errorf("eutil: bad boolean: %q", b)
				}
			case "IsHidden":
				switch b := string(t); b {
				case "Y", "N":
					f.IsHidden = b == "Y"
				default:
					return fmt.Errorf("eutil: bad boolean: %q", b)
				}
			case "IsRangable":
				switch b := string(t); b {
				case "Y", "N":
					f.IsRangeable = b == "Y"
				default:
					return fmt.Errorf("eutil: bad boolean: %q", b)
				}
			case "IsTruncatable":
				switch b := string(t); b {
				case "Y", "N":
					f.IsTruncatable = b == "Y"
				default:
					return fmt.Errorf("eutil: bad boolean: %q", b)
				}
			case "Field":
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
			if t.Name.Local == "Field" {
				return nil
			}
		}
	}
	return nil
}

type DbLink struct {
	Name        string
	FullName    string
	Description string
	DbTo        string
}

func (d *DbLink) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
	for {
		t, err := dec.Token()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.Push(t.Name.Local)
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "Name":
				d.Name = string(t)
			case "FullName":
				d.FullName = string(t)
			case "Description":
				d.Description = string(t)
			case "DbTo":
				d.DbTo = string(t)
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
			if t.Name.Local == "LinkList" {
				return nil
			}
		}
	}
	return nil
}

type DbInfo struct {
	DbName      string
	MenuName    string
	Description string
	Count       int
	LastUpdate  string
	FieldList   []Field
	LinkList    []DbLink
}

func (d *DbInfo) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
	var ldone, fdone bool
	for {
		t, err := dec.Token()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.Push(t.Name.Local)
			switch t.Name.Local {
			case "FieldList":
				fdone = false
			case "LinkList":
				ldone = false
			case "Field":
				var f Field
				err := f.Unmarshal(dec, st[len(st)-1:])
				if (f != Field{}) {
					d.FieldList = append(d.FieldList, f)
				}
				if err != nil {
					return err
				}
				st = st.Drop()
				continue
			case "List":
				var l DbLink
				err := l.Unmarshal(dec, st[len(st)-1:])
				if (l != DbLink{}) {
					d.LinkList = append(d.LinkList, l)
				}
				if err != nil {
					return err
				}
				st = st.Drop()
				continue
			}
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "DbName":
				d.DbName = string(t)
			case "MenuName":
				d.MenuName = string(t)
			case "Description":
				d.Description = string(t)
			case "Count":
				c, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				d.Count = c
			case "LastUpdate":
				d.LastUpdate = string(t)
			case "FieldList", "LinkList", "DbInfo":
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
			switch t.Name.Local {
			case "FieldList":
				fdone = true
			case "LinkList":
				ldone = true
			}
			if ldone && fdone {
				return nil
			}
		}
	}
	return nil
}
