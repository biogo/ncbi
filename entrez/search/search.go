// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"encoding/xml"
	"errors"
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

type TranslationStack []Node

// AST returns the root node of an abstract syntax tree of the query. Nodes in the
// TranslationStack are altered by this method.
func (ts *TranslationStack) AST() (Node, error) {
	if ts == nil || len(*ts) == 0 {
		return nil, nil
	}

	n, _ := (*ts)[len(*ts)-1].Consume((*ts)[:len(*ts)-1])
	if n == nil {
		return nil, errors.New("entrez: translation stack corrupted")
	}

	return n, nil
}

var _ xml.Unmarshaler = (*TranslationStack)(nil)

func (ts *TranslationStack) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	*ts = (*ts)[:0]
	var (
		field string
		tm    = &Term{}
	)
	for {
		tok, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch tok := tok.(type) {
		case xml.StartElement:
			field = tok.Name.Local
		case xml.CharData:
			if field == "" {
				continue
			}
			switch field {
			case "OP":
				o := &Op{Operation: string(tok)}
				*ts = append(*ts, o)
			case "Term":
				tm.Term = string(tok)
			case "Field":
				tm.Field = string(tok)
			case "Count":
				c, err := strconv.Atoi(string(tok))
				if err != nil {
					return err
				}
				tm.Count = c
			case "Explode":
				if len(tok) != 1 {
					return errors.New("entrez: bad boolean")
				}
				switch tok[0] {
				case 'Y', 'N':
					tm.Explode = tok[0] == 'Y'
				default:
					return errors.New("entrez: bad boolean")
				}
			}
		case xml.EndElement:
			if tok.Name.Local == "TermSet" {
				*ts, tm = append(*ts, tm), &Term{}
			}
			field = ""
		}
	}
}

// A Node is an element of the ESearch translation stack.
type Node interface {
	// Consume takes the contents of the translation stack and returns a Node
	// representing the root node of an AST referring to nodes within the stack,
	// and any remaining nodes.
	Consume([]Node) (Node, []Node)
}

type Op struct {
	Operation string
	Operands  []Node
}

func (o *Op) Consume(s []Node) (Node, []Node) {
	var n Node
	switch o.Operation {
	case "AND", "OR", "NOT":
		// Handle binary operators.
		if len(s) < 2 {
			return nil, s
		}
		o.Operands = make([]Node, 2)
		n, s = s[len(s)-1].Consume(s[:len(s)-1])
		if n == nil {
			return nil, s
		}
		o.Operands[1] = n
		if len(s) < 1 {
			return nil, s
		}
		n, s = s[len(s)-1].Consume(s[:len(s)-1])
		if n == nil {
			return nil, s
		}
		o.Operands[0] = n
		return o, s
	case "RANGE", "GROUP": // It's still not entirely clear that RANGE is unary, but this seems to be the case.
		// Handle unary operators.
		if len(s) < 1 {
			return nil, s
		}
		n, s = s[len(s)-1].Consume(s[:len(s)-1])
		if n == nil {
			return nil, s
		}
		o.Operands = []Node{n}
		return o, s
	}
	return nil, s
}

type Term struct {
	Term    string
	Field   string
	Count   int
	Explode bool
}

func (tm *Term) Consume(s []Node) (Node, []Node) { return tm, s }

type Translation struct {
	From string
	To   string
}

type NotFound struct {
	Phrase []string `xml:"PhraseNotFound"`
	Field  []string `xml:"FieldNotFound"`
}

type Warnings struct {
	Ignored  []string `xml:"PhraseIgnored"`
	NotFound []string `xml:"QuotedPhraseNotFound"`
	Message  []string `xml:"OutputMessage"`
}
