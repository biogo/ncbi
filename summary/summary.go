// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package summary

import (
	"code.google.com/p/biogo.entrez/stack"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

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
	Value string
	Name  string
	Type  string
}

func (i *Item) Unmarshal(dec *xml.Decoder) error {
	for {
		t, err := dec.Token()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		switch t := t.(type) {
		case xml.StartElement:
			return errors.New("entrez: unexpected tag")
		case xml.EndElement:
			if t.Name.Local == "Item" {
				return nil
			}
		case xml.CharData:
			i.Value = string(t)
		}
	}
	return nil
}

type Doc struct {
	Id    int
	Items []Item
}

func (d *Doc) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
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
			if t.Name.Local == "Item" {
				var i Item
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "Name":
						i.Name = attr.Value
					case "Type":
						i.Type = attr.Value
					default:
						return fmt.Errorf("entrez: unknown attribute: %q", attr.Name.Local)
					}
				}
				err := i.Unmarshal(dec)
				if !(reflect.DeepEqual(i, Item{})) {
					d.Items = append(d.Items, i)
				}
				if err != nil {
					return err
				}
				continue
			}
			st = st.Push(t.Name.Local)
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "Id":
				c, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				d.Id = c
			case "DocSum":
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
	return nil
}
