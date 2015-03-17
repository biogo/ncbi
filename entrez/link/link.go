// Copyright ©2013 The bíogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package link

import (
	"encoding/xml"
	"errors"
	"io"
	"strconv"
)

// <!--
//                 This is the Current DTD for Entrez eLink
// $Id: eLink_101123.dtd 349314 2012-01-09 23:26:00Z fialkov $
// -->
// <!-- ================================================================= -->
//
// <!ELEMENT	ERROR			(#PCDATA)>	<!-- .+ -->
// <!ELEMENT	Info			(#PCDATA)>	<!-- .+ -->
//
// <!ELEMENT	Id				(#PCDATA)>	<!-- \d+ -->
// <!ATTLIST	Id
// 			HasLinkOut  (Y|N)	#IMPLIED
// 			HasNeighbor (Y|N)	#IMPLIED
// 			>
//
// <!ELEMENT	Score			(#PCDATA)>	<!-- \d+ -->
// <!ELEMENT	DbFrom			(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	DbTo			(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	LinkName		(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	WebEnv			(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	MenuTag			(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	HtmlTag			(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	Priority		(#PCDATA)>	<!-- \S+ -->
//
// <!ELEMENT	IdList		(Id*)>
//
// <!-- cmd=neighbor -->
// <!ELEMENT	Link		(Id, Score?)>
// <!ELEMENT	QueryKey		(#PCDATA)>
//
// <!ELEMENT	LinkSetDb	(DbTo, LinkName, (Link*|Info), ERROR?)>
// <!ELEMENT	LinkSetDbHistory	(DbTo, LinkName, (QueryKey|Info), ERROR?)>
//
// <!-- cmd=llinks -->
//
// <!ELEMENT	Url			    (#PCDATA)>	<!-- \S+ -->
// <!ATTLIST	Url			LNG     (DA|DE|EN|EL|ES|FR|IT|IW|JA|NL|NO|RU|SV|ZH)     "EN">
//
// <!ELEMENT	IconUrl			(#PCDATA)>	<!-- \S+ -->
// <!ATTLIST	IconUrl		LNG     (DA|DE|EN|EL|ES|FR|IT|IW|JA|NL|NO|RU|SV|ZH)     "EN">
//
// <!ELEMENT	SubjectType		(#PCDATA)>	<!-- .+ -->
// <!ELEMENT	Category		(#PCDATA)>	<!-- .+ -->
// <!ELEMENT	Attribute		(#PCDATA)>	<!-- .+ -->
// <!--ELEMENT	LinkName		(#PCDATA)-->	<!--defined in neighbor section--><!-- \S+ -->
// <!ELEMENT	Name			(#PCDATA)>	<!-- .+ -->
// <!ELEMENT	NameAbbr		(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	SubProvider		(#PCDATA)>
//
// <!ELEMENT   FirstChar		(#PCDATA)>
//
// <!ELEMENT	Provider (
// 				Name,
// 				NameAbbr,
// 				Id,
// 				Url,
// 				IconUrl?
// 			)>
//
// <!ELEMENT	ObjUrl	(
// 				Url,
// 				IconUrl?,
// 				LinkName?,
//              SubjectType*,
// 				Category*,
//              Attribute*,
//              Provider,
//              SubProvider?
// 			)>
//
// <!ELEMENT	IdUrlSet	(Id,(ObjUrl+|Info))>
//
// <!ELEMENT	FirstChars  (FirstChar*)>
//
// <!ELEMENT	LinkInfo	(DbTo, LinkName, MenuTag?, HtmlTag?, Url?, Priority)>
// <!ELEMENT	IdLinkSet	(Id, LinkInfo*)>
// <!ELEMENT	IdUrlList	(IdUrlSet* | FirstChars*)>
//
// <!-- cmd=ncheck & lcheck & acheck -->
// <!ELEMENT	IdCheckList	((Id|IdLinkSet)*,ERROR?)>
//
// <!-- Common -->
// <!ELEMENT	LinkSet		(DbFrom,
// 				((IdList?, ((ERROR?, LinkSetDb)*  |  (LinkSetDbHistory*, WebEnv))) | IdUrlList | IdCheckList | ERROR), ERROR?
// 				)>
//
// <!ELEMENT	eLinkResult	(LinkSet*, ERROR?)>

type Link struct {
	Id    Id   `xml:"Id"`
	Score *int `xml:"Score"`
}

type LinkSetDb struct {
	DbTo     string  `xml:"DbTo"`
	LinkName string  `xml:"LinkName"`
	Link     []Link  `xml:"Link"`
	Info     *string `xml:"Info"`
	Err      error   `xml:"ERROR"`
}

type Url struct {
	Url  string `xml:",chardata"`
	Lang string `xml:"LNG,attr"`
}

type Provider struct {
	Name     string `xml:"Name"`
	NameAbbr string `xml:"NameAbbr"`
	Id       Id     `xml:"Id"`
	Url      Url    `xml:"Url"`
	IconUrl  *Url   `xml:"IconUrl"`
}

type ObjUrl struct {
	Url         Url      `xml:"Url"`
	IconUrl     *Url     `xml:"IconUrl"`
	LinkName    *string  `xml:"LinkName"`
	SubjectType []string `xml:"SubjectType"`
	Category    []string `xml:"Category"`
	Attribute   []string `xml:"Attribute"`
	Provider    Provider `xml:"Provider"`
	SubProvider *string  `xml:"SubProvider"`
}

type IdUrlSet struct {
	Id     Id       `xml:"Id"`
	ObjUrl []ObjUrl `xml:"ObjUrl"`
	Info   *string  `xml:"Info"`
}

type IdUrlList struct {
	IdUrlSets  []IdUrlSet `xml:"IdUrlSet"`
	FirstChars [][]string `xml:"FirstChars>FirstChar"`
}

type IdLinkSet struct {
	Id       Id         `xml:"Id"`
	LinkInfo []LinkInfo `xml:"LinkInfo"`
}

type LinkInfo struct {
	DbTo     string  `xml:"DbTo"`
	LinkName string  `xml:"LinkName"`
	MenuTag  *string `xml:"MenuTag"`
	HtmlTag  *string `xml:"HtmlTag"`
	Url      *Url    `xml:"Url"`
	Priority int     `xml:"Priority"`
}

type Id struct {
	Id          int   `xml:"Id"`
	HasLinkOut  *bool `xml:",attr"`
	HasNeighbor *bool `xml:",attr"`
}

var _ xml.Unmarshaler = (*Id)(nil)

func (id *Id) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	for _, attr := range start.Attr {
		var b *bool
		switch attr.Value {
		case "Y", "N":
			b = new(bool)
			*b = attr.Value == "Y"
		default:
			return errors.New("entrez: bad boolean")
		}
		switch attr.Name.Local {
		case "HasLinkOut":
			id.HasNeighbor = b
		case "HasNeighbor":
			id.HasNeighbor = b
		}
	}

	for {
		tok, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		switch tok := tok.(type) {
		case xml.CharData:
			i, err := strconv.Atoi(string(tok))
			if err != nil {
				return err
			}
			id.Id = i
		}
	}
}

type IdCheckList struct {
	Id        []Id        `xml:"Id"`
	IdLinkSet []IdLinkSet `xml:"IdLinkSet"`
	Err       error       `xml:"ERROR"`
}

type LinkSetDbHistory struct {
	DbTo     string  `xml:"DbTo"`
	LinkName string  `xml:"LinkName"`
	QueryKey *int    `xml:"QueryKey"`
	Info     *string `xml:"Info"`
	Err      error   `xml:"ERROR"`
}

type LinkSet struct {
	DbFrom           string             `xml:"DbFrom"`
	IdList           []Id               `xml:"IdList>Id"`
	Neighbor         []LinkSetDb        `xml:"LinkSetDb"`
	LinkSetDbHistory []LinkSetDbHistory `xml:"LinkSetDbHistory"`
	WebEnv           *string            `xml:"WebEnv"`
	IdUrlList        *IdUrlList         `xml:"IdUrlList"`
	IdCheckList      *IdCheckList       `xml:"IdCheckList"`
	Err              []string           `xml:"ERROR"`
}
