// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"code.google.com/p/biogo.entrez/search"
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

// A Search holds the deserialised results of an ESearch request.
type Search struct {
	Database         string
	Count            int                     `xml:"Count"`
	RetMax           int                     `xml:"RetMax"`
	RetStart         int                     `xml:"RetStart"`
	QueryKey         *int                    `xml:"QueryKey"`
	WebEnv           *string                 `xml:"WebEnv"`
	IdList           []int                   `xml:"IdList>Id"`
	Translations     []search.Translation    `xml:"TranslationSet>Translation"`
	TranslationStack search.TranslationStack `xml:"TranslationStack"`
	QueryTranslation *string                 `xml:"QueryTranslation"`
	Err              *string                 `xml:"ERROR"`
	NotFound         *search.NotFound        `xml:"ErrorList"`
	Warnings         *search.Warnings        `xml:"WarningList"`
}

// History returns a History containing the Search's query key and web environment.
func (s *Search) History() *History {
	var h *History
	if s.QueryKey != nil {
		h = &History{QueryKey: *s.QueryKey}
	}
	if s.WebEnv != nil {
		if h == nil {
			h = &History{}
		}
		h.WebEnv = *s.WebEnv
	}
	return h
}
