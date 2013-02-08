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
//                 This is the Current DTD for Entrez ePost
// $Id: ePost_020511.dtd 161288 2009-05-26 18:34:21Z fialkov $
// -->
// <!-- ================================================================= -->
//
// <!ELEMENT	Id		(#PCDATA)>	<!-- \d+ -->
//
// <!ELEMENT	InvalidIdList	(Id+)>
// <!ELEMENT       QueryKey        (#PCDATA)>	<!-- \d+ -->
// <!ELEMENT       WebEnv          (#PCDATA)>	<!-- \S+ -->
// <!ELEMENT       ERROR           (#PCDATA)>	<!-- .+ -->
//
// <!ELEMENT     ePostResult       (InvalidIdList?,(QueryKey,WebEnv)?,ERROR?)>

// A Post holds the deserialised results of an EPost request.
type Post struct {
	InvalidIds []int
	QueryKey   *int
	WebEnv     *string
	Err        error
}

// Unmarshal fills the fields of a Post from an XML stream read from r.
func (p *Post) Unmarshal(r io.Reader) error {
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
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
			case "Id":
				if st.peek(1) != "InvalidIdList" {
					return fmt.Errorf("entrez: unexpected tag: %q", name)
				}
				id, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				p.InvalidIds = append(p.InvalidIds, id)
			case "QueryKey":
				k, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				p.QueryKey = &k
			case "WebEnv":
				s := string(t)
				p.WebEnv = &s
			case "ERROR":
				p.Err = Error(string(t))
			case "ePostResult", "InvalidIdList":
			default:
				p.Err = Error(fmt.Sprintf("unknown name: %q", name))
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
