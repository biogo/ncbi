// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"code.google.com/p/biogo.entrez/stack"
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
)

// <!--
//                 This is the Current DTD for Entrez eSearch
// $Id: eSearch_020511.dtd 85163 2006-06-28 17:35:21Z olegh $
// -->
// <!-- ================================================================= -->
//
// <!ELEMENT       Count           (#PCDATA)>	<!-- \d+ -->
// <!ELEMENT       RetMax          (#PCDATA)>	<!-- \d+ -->
// <!ELEMENT       RetStart        (#PCDATA)>	<!-- \d+ -->
// <!ELEMENT       Id              (#PCDATA)>	<!-- \d+ -->
//
// <!ELEMENT       From            (#PCDATA)>	<!-- .+ -->
// <!ELEMENT       To              (#PCDATA)>	<!-- .+ -->
// <!ELEMENT       Term            (#PCDATA)>	<!-- .+ -->
//
// <!ELEMENT       Field           (#PCDATA)>	<!-- .+ -->
//
// <!ELEMENT        QueryKey       (#PCDATA)>	<!-- \d+ -->
// <!ELEMENT        WebEnv         (#PCDATA)>	<!-- \S+ -->
//
// <!ELEMENT       Explode         (#PCDATA)>	<!-- (Y|N) -->
// <!ELEMENT       OP              (#PCDATA)>	<!-- (AND|OR|NOT|RANGE|GROUP) -->
// <!ELEMENT       IdList          (Id*)>
//
// <!ELEMENT       Translation     (From, To)>
// <!ELEMENT       TranslationSet  (Translation*)>
//
// <!ELEMENT       TermSet (Term, Field, Count, Explode)>
// <!ELEMENT       TranslationStack        ((TermSet|OP)*)>
//
// <!-- Error message tags  -->
//
// <!ELEMENT        ERROR                  (#PCDATA)>	<!-- .+ -->
//
// <!ELEMENT        OutputMessage		    (#PCDATA)>	<!-- .+ -->
// <!ELEMENT        QuotedPhraseNotFound   (#PCDATA)>	<!-- .+ -->
// <!ELEMENT        PhraseIgnored          (#PCDATA)>	<!-- .+ -->
// <!ELEMENT        FieldNotFound          (#PCDATA)>	<!-- .+ -->
// <!ELEMENT        PhraseNotFound         (#PCDATA)>	<!-- .+ -->
// <!ELEMENT        QueryTranslation       (#PCDATA)>	<!-- .+ -->
//
// <!ELEMENT        ErrorList      (PhraseNotFound*,FieldNotFound*)>
// <!ELEMENT        WarningList   	(PhraseIgnored*,
// 				QuotedPhraseNotFound*,
// 				OutputMessage*)>
// <!-- Response tags -->
//
//
// <!ELEMENT       eSearchResult   (((
//                                 Count,
//                                     (RetMax,
//                                     RetStart,
//                                     QueryKey?,
//                                     WebEnv?,
//                                     IdList,
//                                     TranslationSet,
//                                     TranslationStack?,
//                                     QueryTranslation
//                                     )?
//                                 ) | ERROR),
// 				ErrorList?,
// 				WarningList?
// 				)>

type Op string

func (o Op) Consume(s []Node) Node {
	// TODO Flesh out Op to be a struct:
	//
	//  type Op struct {
	//  	Operation string
	//  	Operands  []Node
	//  }
	//
	// Then we can build an AST for the search. To do this we need to understand what
	// RANGE and GROUP actually do - this is not specified.
	switch o {
	case "AND", "OR", "NOT", "RANGE", "GROUP":
		return o
	}
	return nil
}

type Term struct {
	Term    string
	Field   string
	Count   int
	Explode bool
}

func (tm Term) Consume(_ []Node) Node { return tm }

func (tm *Term) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
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
			case "Term":
				tm.Term = string(t)
			case "Field":
				tm.Field = string(t)
			case "Count":
				c, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				tm.Count = c
			case "Explode":
				switch b := string(t); b {
				case "Y", "N":
					tm.Explode = b == "Y"
				default:
					return fmt.Errorf("entrez: bad boolean: %q", b)
				}
			case "TermSet":
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

type Translation struct {
	From string
	To   string
}

func (tr *Translation) Unmarshal(dec *xml.Decoder, st stack.Stack) error {
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
			case "From":
				tr.From = string(t)
			case "To":
				tr.To = string(t)
			case "TranslationSet", "Translation":
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

// A Node is an element of the ESearch translation stack.Stack.
type Node interface {
	Consume([]Node) Node
}
