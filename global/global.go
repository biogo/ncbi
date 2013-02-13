// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package global

import (
	"code.google.com/p/biogo.entrez/stack"
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

func (r *Result) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
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
			st = st.Push(t.Name.Local)
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
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
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
			if st.Empty() {
				return nil
			}
		}
	}
	panic("cannot reach")
}
