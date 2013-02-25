// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package search

import (
	"bytes"
	"code.google.com/p/biogo.entrez/xml"
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

func (ts *TranslationStack) UnmarshalXML(b []byte) error {
	*ts = (*ts)[:0]
	dec := xml.NewDecoder(bytes.NewReader(b))
	var (
		field string
		tm    Term
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
				*ts = append(*ts, Op(string(tok)))
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
				*ts, tm = append(*ts, tm), Term{}
			}
			field = ""
		}
	}

	panic("cannot reach")
}

// A Node is an element of the ESearch translation stack.
type Node interface {
	Consume([]Node) Node
}

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
