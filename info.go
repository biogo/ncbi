// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

// <!--
//                 This is the Current DTD for Entrez eInfo
// $Id: eInfo_020511.dtd 94706 2006-12-04 21:51:33Z olegh $
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
//
//
// <!ELEMENT	DbList		(DbName+)>
//
// <!ELEMENT	Field		(Name,
//                 FullName,
// 				Description,
// 				TermCount,
// 				IsDate,
// 				IsNumerical,
// 				SingleToken,
// 				Hierarchy,
//                 IsHidden)>
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
	Name        string
	FullName    string
	Description string
	TermCount   int
	IsDate      bool
	IsNumerical bool
	SingleToken bool
	Hierarchy   bool
	IsHidden    bool
}

func (f *Field) unmarshal(dec *xml.Decoder, st stack) error {
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
			st = st.push(t.Name.Local)
		case xml.CharData:
			switch name := st.peek(0); name {
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
			case "Field":
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.pair(t.Name.Local)
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

func (d *DbLink) unmarshal(dec *xml.Decoder, st stack) error {
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
			st = st.push(t.Name.Local)
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
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
			st, err = st.pair(t.Name.Local)
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

func (d *DbInfo) unmarshal(dec *xml.Decoder, st stack) error {
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
			st = st.push(t.Name.Local)
			switch t.Name.Local {
			case "FieldList":
				fdone = false
			case "LinkList":
				ldone = false
			case "Field":
				var f Field
				err := f.unmarshal(dec, st[len(st)-1:])
				if (f != Field{}) {
					d.FieldList = append(d.FieldList, f)
				}
				if err != nil {
					return err
				}
				st = st.drop()
				continue
			case "List":
				var l DbLink
				err := l.unmarshal(dec, st[len(st)-1:])
				if (l != DbLink{}) {
					d.LinkList = append(d.LinkList, l)
				}
				if err != nil {
					return err
				}
				st = st.drop()
				continue
			}
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
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
			st, err = st.pair(t.Name.Local)
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

// An Info holds the deserialised results of an EInfo request.
type Info struct {
	DbList []string
	DbInfo []DbInfo
	Err    error
}

// Unmarshal fills the fields of an Info from an XML stream read from r.
func (i *Info) Unmarshal(r io.Reader) error {
	dec := xml.NewDecoder(r)
	var st stack
	for {
		t, err := dec.Token()
		if err != nil {
			if err != io.EOF {
				return err
			}
			if !st.empty() {
				return io.ErrUnexpectedEOF
			}
			break
		}
		switch t := t.(type) {
		case xml.ProcInst:
		case xml.Directive:
		case xml.StartElement:
			st = st.push(t.Name.Local)
			if t.Name.Local == "DbInfo" {
				var d DbInfo
				err := d.unmarshal(dec, st[len(st)-1:])
				if !reflect.DeepEqual(d, DbInfo{}) {
					i.DbInfo = append(i.DbInfo, d)
				}
				if err != nil {
					return err
				}
				st.drop()
				continue
			}
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
			case "DbName":
				if st.peek(1) != "DbList" {
					return fmt.Errorf("entrez: unexpected tag: %q", name)
				}
				i.DbList = append(i.DbList, string(t))
			case "ERROR":
				i.Err = Error(string(t))
			case "eInfoResult", "DbList", "DbInfo":
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.pair(t.Name.Local)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
