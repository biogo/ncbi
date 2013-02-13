// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"code.google.com/p/biogo.entrez/stack"
	"encoding/xml"
	"errors"
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
	History    *History
	Err        error
}

// Unmarshal fills the fields of a Post from an XML stream read from r.
func (p *Post) Unmarshal(r io.Reader) error {
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
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
			case "Id":
				if st.Peek(1) != "InvalidIdList" {
					return fmt.Errorf("entrez: unexpected tag: %q", name)
				}
				id, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				p.InvalidIds = append(p.InvalidIds, id)
			case "QueryKey":
				if p.History == nil {
					p.History = &History{}
				}
				k, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				p.History.QueryKey = k
			case "WebEnv":
				if p.History == nil {
					p.History = &History{}
				}
				p.History.WebEnv = string(t)
			case "ERROR":
				p.Err = errors.New(string(t))
			case "ePostResult", "InvalidIdList":
			default:
				p.Err = fmt.Errorf("unknown name: %q", name)
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
