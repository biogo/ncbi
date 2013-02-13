// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"code.google.com/p/biogo.entrez/search"
	"code.google.com/p/biogo.entrez/stack"
	"encoding/xml"
	"errors"
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

type NotFound struct {
	Type  string
	Value string
}

func (e NotFound) Error() string { return fmt.Sprintf("entrez: %s: %q", e.Type, e.Value) }

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

// A Search holds the deserialised results of an ESearch request.
type Search struct {
	Database         string
	Count            int
	RetMax           int
	RetStart         int
	History          *History
	IdList           []int
	Translations     []search.Translation
	TranslationStack []search.Node
	QueryTranslation *string
	Err              error
	Errors           []error
	Warnings         []error
}

// Unmarshal fills the fields of a Search from an XML stream read from r.
func (s *Search) Unmarshal(r io.Reader) error {
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
			switch t.Name.Local {
			case "Translation":
				var tr search.Translation
				err := tr.Unmarshal(dec, st[len(st)-1:])
				if (tr != search.Translation{}) {
					s.Translations = append(s.Translations, tr)
				}
				if err != nil {
					return err
				}
				st = st.Drop()
			case "TermSet":
				var tm search.Term
				err := tm.Unmarshal(dec, st[len(st)-1:])
				if (tm != search.Term{}) {
					s.TranslationStack = append(s.TranslationStack, tm)
				}
				if err != nil {
					return err
				}
				st = st.Drop()
			}
		case xml.CharData:
			if st.Empty() {
				continue
			}
			switch name := st.Peek(0); name {
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
				o := search.Op(string(t))
				if o.Consume(s.TranslationStack) == nil {
					return fmt.Errorf("entrez: illegal operator: %q", o)
				}
				s.TranslationStack = append(s.TranslationStack, o)
			case "QueryTranslation":
				st := string(t)
				s.QueryTranslation = &st
			case "ERROR":
				s.Err = errors.New(string(t))
			case "PhraseNotFound", "FieldNotFound":
				if st.Peek(1) != "ErrorList" {
					return fmt.Errorf("entrez: unexpected tag: %q", name)
				}
				switch name {
				case "PhraseNotFound":
					s.Errors = append(s.Errors, NotFound{Type: phraseNotFound, Value: string(t)})
				case "FieldNotFound":
					s.Errors = append(s.Errors, NotFound{Type: fieldNotFound, Value: string(t)})
				}
			case "PhraseIgnored", "QuotedPhraseNotFound", "OutputMessage":
				if st.Peek(1) != "WarningList" {
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
			st, err = st.Pair(t.Name.Local)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
