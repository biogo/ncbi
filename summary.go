// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
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

func (i *Item) unmarshal(dec *xml.Decoder) error {
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

func (d *Doc) unmarshal(dec *xml.Decoder, st stack) error {
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
				err := i.unmarshal(dec)
				if !(reflect.DeepEqual(i, Item{})) {
					d.Items = append(d.Items, i)
				}
				if err != nil {
					return err
				}
				continue
			}
			st = st.push(t.Name.Local)
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
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
			st, err = st.pair(t.Name.Local)
			if err != nil {
				return err
			}
			if st.empty() {
				return nil
			}
		}
	}
	return nil
}

// A Summary holds the deserialised results of an ESummary request.
type Summary struct {
	Database string
	Docs     []Doc
	Err      error
}

// Unmarshal fills the fields of a Summary from an XML stream read from r.
func (s *Summary) Unmarshal(r io.Reader) error {
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
			if t.Name.Local == "DocSum" {
				var d Doc
				err := d.unmarshal(dec, st[len(st)-1:])
				if !(reflect.DeepEqual(d, Doc{})) {
					s.Docs = append(s.Docs, d)
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
			case "ERROR":
				s.Err = Error(string(t))
			case "eSummaryResult":
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
