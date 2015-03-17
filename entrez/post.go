// Copyright ©2013 The bíogo Authors. All rights reserved.
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
	InvalidIds []int `xml:"InvalidIdList>Id"`
	*History
	Err *string `xml:"ERROR"`
}
