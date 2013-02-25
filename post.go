// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

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
	InvalidIds []int   `xml:"InvalidIdList>Id"`
	QueryKey   *int    `xml:"QueryKey"`
	WebEnv     *string `xml:"WebEnv"`
	Err        *string `xml:"ERROR"`
}

// History returns a History containing the Post's query key and web environment.
func (p *Post) History() *History {
	var h *History
	if p.QueryKey != nil {
		h = &History{QueryKey: *p.QueryKey}
	}
	if p.WebEnv != nil {
		if h == nil {
			h = &History{}
		}
		h.WebEnv = *p.WebEnv
	}
	return h
}
