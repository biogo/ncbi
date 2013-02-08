// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
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

type Result struct {
	Database string
	MenuName string
	Count    int
	Status   string
}

func (r *Result) unmarshal(dec *xml.Decoder, st stack) error {
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return io.ErrUnexpectedEOF
			}
			return err
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
			case "DbName":
				r.Database = string(t)
			case "MenuName":
				r.MenuName = string(t)
			case "Count":
				c, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				r.Count = c
			case "Status":
				r.Status = string(t)
			default:
				return fmt.Errorf("entrez: unknown name: %q", name)
			}
		case xml.EndElement:
			st, err = st.pair(t.Name.Local)
			if err != nil {
				return err
			}
			if st.empty() {
				return nil
			}
		}
	}
	panic("cannot reach")
}

// A Global holds the deserialised results of an EGQuery request.
type Global struct {
	Query   string
	Results []Result
}

// Unmarshal fills the fields of a Global from an XML stream read from r.
func (g *Global) Unmarshal(r io.Reader) error {
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
			if t.Name.Local == "ResultItem" {
				var res Result
				err := res.unmarshal(dec, st[len(st)-1:])
				g.Results = append(g.Results, res)
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
			case "Term":
				g.Query = string(t)
			case "eGQueryResult", "Result":
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
