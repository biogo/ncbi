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

const (
	phraseNotFound = "phrase not found"
	fieldNotFound  = "field not found"
)

type NotFoundError struct {
	Type  string
	Value string
}

func (e NotFoundError) Error() string { return fmt.Sprintf("entrez: %s: %q", e.Type, e.Value) }

const (
	phraseIgnored        = "phrase ignored"
	quotedPhraseNotFound = "quoted phrase not found"
	outputMessage        = "output message"
)

type Warning struct {
	Type  string
	Value string
}

func (w Warning) Error() string { return fmt.Sprintf("entrez: warning: %s: %q", w.Type, w.Value) }

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

func (tm *Term) unmarshal(dec *xml.Decoder, st stack) error {
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

type Translation struct {
	From string
	To   string
}

func (tr *Translation) unmarshal(dec *xml.Decoder, st stack) error {
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
			case "From":
				tr.From = string(t)
			case "To":
				tr.To = string(t)
			case "TranslationSet", "Translation":
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

type Node interface {
	Consume([]Node) Node
}

// A Search holds the deserialised results of an ESearch request.
type Search struct {
	Database         string
	Count            int
	RetMax           int
	RetStart         int
	History          *History
	IdList           []int
	Translations     []Translation
	TranslationStack []Node
	QueryTranslation *string
	Err              error
	Errors           []error
	Warnings         []error
}

// Unmarshal fills the fields of a Search from an XML stream read from r.
func (s *Search) Unmarshal(r io.Reader) error {
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
			switch t.Name.Local {
			case "Translation":
				var tr Translation
				err := tr.unmarshal(dec, st[len(st)-1:])
				if (tr != Translation{}) {
					s.Translations = append(s.Translations, tr)
				}
				if err != nil {
					return err
				}
				st = st.drop()
			case "TermSet":
				var tm Term
				err := tm.unmarshal(dec, st[len(st)-1:])
				if (tm != Term{}) {
					s.TranslationStack = append(s.TranslationStack, tm)
				}
				if err != nil {
					return err
				}
				st = st.drop()
			}
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
			case "Count":
				c, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				s.Count = c
			case "RetMax":
				m, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				s.RetMax = m
			case "RetStart":
				st, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				s.RetStart = st
			case "QueryKey":
				if s.History == nil {
					s.History = &History{}
				}
				k, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				s.History.QueryKey = k
			case "WebEnv":
				if s.History == nil {
					s.History = &History{}
				}
				s.History.WebEnv = string(t)
			case "Id":
				id, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				s.IdList = append(s.IdList, id)
			case "OP":
				o := Op(string(t))
				if o.Consume(s.TranslationStack) == nil {
					return fmt.Errorf("entrez: illegal operator: %q", o)
				}
				s.TranslationStack = append(s.TranslationStack, o)
			case "QueryTranslation":
				st := string(t)
				s.QueryTranslation = &st
			case "ERROR":
				s.Err = Error(string(t))
			case "PhraseNotFound", "FieldNotFound":
				if st.peek(1) != "ErrorList" {
					return fmt.Errorf("entrez: unexpected tag: %q", name)
				}
				switch name {
				case "PhraseNotFound":
					s.Errors = append(s.Errors, NotFoundError{Type: phraseNotFound, Value: string(t)})
				case "FieldNotFound":
					s.Errors = append(s.Errors, NotFoundError{Type: fieldNotFound, Value: string(t)})
				}
			case "PhraseIgnored", "QuotedPhraseNotFound", "OutputMessage":
				if st.peek(1) != "WarningList" {
					return fmt.Errorf("entrez: unexpected tag: %q", name)
				}
				switch name {
				case "PhraseIgnored":
					s.Warnings = append(s.Warnings, Warning{Type: phraseIgnored, Value: string(t)})
				case "QuotedPhraseNotFound":
					s.Warnings = append(s.Warnings, Warning{Type: quotedPhraseNotFound, Value: string(t)})
				case "OutputMessage":
					s.Warnings = append(s.Warnings, Warning{Type: outputMessage, Value: string(t)})
				}
			case "eSearchResult", "IdList", "TranslationSet", "TranslationStack", "ErrorList", "WarningList":
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
