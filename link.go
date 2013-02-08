// Copyright ©2013 The bíogo.entrez Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package entrez

import (
	"encoding/xml"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

// <!--
//                 This is the Current DTD for Entrez eLink
// $Id: eLink_020511.dtd 56256 2005-02-18 17:13:40Z olegh $
// -->
// <!-- ================================================================= -->
//
// <!ELEMENT       ERROR           (#PCDATA)>	<!-- .+ -->
// <!ELEMENT       Info            (#PCDATA)>	<!-- .+ -->
//
// <!ELEMENT	Id		(#PCDATA)>	<!-- \d+ -->
// <!ATTLIST	Id
// 			HasLinkOut  (Y|N)	#IMPLIED
// 			HasNeighbor (Y|N)	#IMPLIED
// 			>
//
// <!ELEMENT	Score		(#PCDATA)>	<!-- \d+ -->
// <!ELEMENT	DbFrom		(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	DbTo		(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	LinkName	(#PCDATA)>	<!-- \S+ -->
//
// <!ELEMENT	IdList		(Id*)>
//
// <!-- cmd=neighbor -->
// <!ELEMENT	Link		(Id, Score?)>
//
// <!ELEMENT	LinkSetDb	(DbTo, LinkName, (Link*|Info), ERROR?)>
//
// <!-- cmd=links -->
//
// <!ELEMENT	Url			(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT	IconUrl		(#PCDATA)>	<!-- \S+ -->
// <!ELEMENT    SubjectType	(#PCDATA)>	<!-- .+ -->
// <!ELEMENT    Attribute	(#PCDATA)>	<!-- .+ -->
// <!ELEMENT    Name		(#PCDATA)>	<!-- .+ -->
// <!ELEMENT    NameAbbr	(#PCDATA)>	<!-- \S+ -->
//
// <!ELEMENT	Provider (
// 				Name,
// 				NameAbbr,
// 				Id,
// 				Url,
// 				IconUrl?
// 			)>
//
//
// <!ELEMENT	ObjUrl	(
// 				Url,
// 				IconUrl?,
// 				LinkName?,
//              SubjectType*,
//              Attribute*,
//              Provider
// 			)>
//
// <!ELEMENT	IdUrlSet	(Id,(ObjUrl+|Info))>
//
// <!ELEMENT	IdUrlList	(IdUrlSet*,ERROR?)>
//
//
// <!-- cmd=ncheck & lcheck -->
// <!ELEMENT	IdCheckList	(Id*,ERROR?)>
//
//
// <!-- Common -->
// <!ELEMENT	LinkSet		(DbFrom,
// 				((IdList?, LinkSetDb*) | IdUrlList | IdCheckList | ERROR)
// 				)>
//
// <!ELEMENT	eLinkResult	(LinkSet*, ERROR?)>

// The following structure observable by: http://eutils.ncbi.nlm.nih.gov/entrez/eutils/elink.fcgi?dbfrom=protein&db=pubmed&id=15718680,157427902&cmd=acheck
// is not described by the above dtd:
// <?xml version="1.0"?>
// <!DOCTYPE eLinkResult PUBLIC "-//NLM//DTD eLinkResult, 23 November 2010//EN" "http://www.ncbi.nlm.nih.gov/entrez/query/DTD/eLink_101123.dtd">
// <eLinkResult>
// 	<LinkSet>
// 		<DbFrom>protein</DbFrom>
// 		<IdCheckList>
// 			<IdLinkSet>
// 				<Id>15718680</Id>
// 				<LinkInfo>
// 					<DbTo>pubmed</DbTo>
// 					<LinkName>protein_pubmed</LinkName>
// 					<MenuTag>PubMed Links</MenuTag>
// 					<HtmlTag>PubMed</HtmlTag>
// 					<Priority>128</Priority>
// 				</LinkInfo>
// 				<LinkInfo>
// 					<DbTo>pubmed</DbTo>
// 					<LinkName>protein_pubmed_refseq</LinkName>
// [snip]
// 					<HtmlTag>LinkOut</HtmlTag>
// 					<Priority>255</Priority>
// 				</LinkInfo>
// 			</IdLinkSet>
// 			<IdLinkSet>
// 				<Id>157427902</Id>
// 				<LinkInfo>
// 					<DbTo>pubmed</DbTo>
// 					<LinkName>protein_pubmed</LinkName>
// [snip]
// 					<Priority>128</Priority>
// 				</LinkInfo>
// 				<LinkInfo>
// 					<DbTo>LinkOut</DbTo>
// 					<LinkName>ExternalLink</LinkName>
// 					<MenuTag>LinkOut</MenuTag>
// 					<HtmlTag>LinkOut</HtmlTag>
// 					<Priority>255</Priority>
// 				</LinkInfo>
// 			</IdLinkSet>
// 		</IdCheckList>
// [snip]

// The Category field in ObjUrl is not mentioned, but certainly exists in the wild.
// And what does 'IMPLIED' mean on a boolean‽

type LinkId struct {
	Id          int
	HasLinkOut  bool
	HasNeighbor bool
	Score       int
}

func (li *LinkId) unmarshal(dec *xml.Decoder, st stack) error {
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
			if t.Name.Local == "Id" {
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "HasLinkOut":
						switch b := attr.Value; b {
						case "Y", "N":
							li.HasLinkOut = b == "Y"
						default:
							return fmt.Errorf("eutil: bad boolean: %q", b)
						}
					case "HasNeighbor":
						switch b := attr.Value; b {
						case "Y", "N":
							li.HasNeighbor = b == "Y"
						default:
							return fmt.Errorf("eutil: bad boolean: %q", b)
						}
					default:
						return fmt.Errorf("entrez: unknown attribute: %q", attr.Name.Local)
					}
				}
			}
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
			case "Id":
				id, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				li.Id = id
			case "Score":
				s, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				li.Score = s
			case "Link":
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

type LinkSetDb struct {
	DbTo     string
	LinkName string
	Link     []LinkId
}

func (ls *LinkSetDb) unmarshal(dec *xml.Decoder, st stack) error {
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
			if t.Name.Local == "Link" {
				var li LinkId
				err := li.unmarshal(dec, st[len(st)-1:])
				if (li != LinkId{}) {
					ls.Link = append(ls.Link, li)
				}
				if err != nil {
					return err
				}
				st = st.drop()
				continue
			}
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
			case "DbTo":
				ls.DbTo = string(t)
			case "LinkName":
				ls.LinkName = string(t)
			case "LinkSetDb":
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

type Provider struct {
	Name     string
	NameAbbr string
	Id       int
	Url      string
	IconUrl  *string
}

func (p *Provider) unmarshal(dec *xml.Decoder, st stack) error {
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
			case "Name":
				p.Name = string(t)
			case "NameAbbr":
				p.NameAbbr = string(t)
			case "Id":
				id, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				p.Id = id
			case "Url":
				p.Url = string(t)
			case "IconUrl":
				s := string(t)
				p.IconUrl = &s
			case "Provider":
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

type ObjUrl struct {
	Url         string
	IconUrl     *string
	LinkName    *string
	SubjectType []string
	Category    []string
	Attribute   []string
	Provider    Provider
}

func (ou *ObjUrl) unmarshal(dec *xml.Decoder, st stack) error {
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
			if t.Name.Local == "Provider" {
				err := ou.Provider.unmarshal(dec, st[len(st)-1:])
				if err != nil {
					return err
				}
				st = st.drop()
				continue
			}
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
			case "Url":
				ou.Url = string(t)
			case "IconUrl":
				s := string(t)
				ou.IconUrl = &s
			case "LinkName":
				s := string(t)
				ou.LinkName = &s
			case "SubjectType":
				ou.SubjectType = append(ou.SubjectType, string(t))
			case "Category":
				ou.Category = append(ou.Category, string(t))
			case "Attribute":
				ou.Attribute = append(ou.Attribute, string(t))
			case "ObjUrl":
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

type IdUrlSet struct {
	Id      int
	ObjUrls []ObjUrl
	Info    string
}

func (us *IdUrlSet) unmarshal(dec *xml.Decoder, st stack) error {
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
			if t.Name.Local == "ObjUrl" {
				var ou ObjUrl
				err := ou.unmarshal(dec, st[len(st)-1:])
				if !reflect.DeepEqual(ou, IdUrlSet{}) {
					us.ObjUrls = append(us.ObjUrls, ou)
				}
				if err != nil {
					return err
				}
				st = st.drop()
				continue
			}
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
			case "Id":
				id, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				us.Id = id
			case "Info":
				us.Info = string(t)
			case "IdUrlSet", "ObjUrl":
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

type IdUrlList struct {
	IdUrlSets []IdUrlSet
	Err       error
}

func (ul *IdUrlList) unmarshal(dec *xml.Decoder, st stack) error {
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
			if t.Name.Local == "IdUrlSet" {
				var us IdUrlSet
				err := us.unmarshal(dec, st[len(st)-1:])
				if !reflect.DeepEqual(us, IdUrlSet{}) {
					ul.IdUrlSets = append(ul.IdUrlSets, us)
				}
				if err != nil {
					return err
				}
				st = st.drop()
				continue
			}
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
			case "ERROR":
				ul.Err = Error(string(t))
			case "IdUrlList":
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

type Id struct {
	Id          int
	HasLinkOut  bool
	HasNeighbor bool
}

type IdChecks struct {
	Ids []Id
	Err error
}

func (ic *IdChecks) unmarshal(dec *xml.Decoder, st stack) error {
	var hasLinkOut, hasNeighbor bool
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
			if t.Name.Local == "Id" {
				for _, attr := range t.Attr {
					switch attr.Name.Local {
					case "HasLinkOut":
						switch b := attr.Value; b {
						case "Y", "N":
							hasLinkOut = b == "Y"
						default:
							return fmt.Errorf("eutil: bad boolean: %q", b)
						}
					case "HasNeighbor":
						switch b := attr.Value; b {
						case "Y", "N":
							hasNeighbor = b == "Y"
						default:
							return fmt.Errorf("eutil: bad boolean: %q", b)
						}
					default:
						return fmt.Errorf("entrez: unknown attribute: %q", attr.Name.Local)
					}
				}
			}
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
			case "ERROR":
				ic.Err = Error(string(t))
			case "Id":
				id, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				ic.Ids = append(ic.Ids, Id{Id: id, HasLinkOut: hasLinkOut, HasNeighbor: hasNeighbor})
				hasLinkOut = false
				hasNeighbor = false
			case "IdCheckList":
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

type LinkSet struct {
	DbFrom     string
	IdList     []int
	LinkSetDbs []LinkSetDb
	IdUrls     []IdUrlList
	IdChecks   *IdChecks
	Err        error
}

func (li *LinkSet) unmarshal(dec *xml.Decoder, st stack) error {
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
			switch t.Name.Local {
			case "IdUrlList":
				var ul IdUrlList
				err := ul.unmarshal(dec, st[len(st)-1:])
				if !reflect.DeepEqual(ul, IdUrlList{}) {
					li.IdUrls = append(li.IdUrls, ul)
				}
				if err != nil {
					return err
				}
				st = st.drop()
			case "LinkSetDb":
				var ls LinkSetDb
				err := ls.unmarshal(dec, st[len(st)-1:])
				if !reflect.DeepEqual(ls, LinkSetDb{}) {
					li.LinkSetDbs = append(li.LinkSetDbs, ls)
				}
				if err != nil {
					return err
				}
				st = st.drop()
			case "IdCheckList":
				var id IdChecks
				err := id.unmarshal(dec, st[len(st)-1:])
				if !reflect.DeepEqual(id, IdChecks{}) {
					li.IdChecks = &id
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
			case "DbFrom":
				li.DbFrom = string(t)
			case "Id":
				if st.peek(1) != "IdList" {
					return fmt.Errorf("entrez: unexpected tag: %q", name)
				}
				id, err := strconv.Atoi(string(t))
				if err != nil {
					return err
				}
				li.IdList = append(li.IdList, id)
			case "LinkSet", "IdList", "IdUrlList":
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

// A Link holds the deserialised results of an ELink request.
type Link struct {
	LinkSets []LinkSet
	Err      error
}

// Unmarshal fills the fields of a Link from an XML stream read from r.
func (l *Link) Unmarshal(r io.Reader) error {
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
			if t.Name.Local == "LinkSet" {
				var ls LinkSet
				err := ls.unmarshal(dec, st[len(st)-1:])
				if !(reflect.DeepEqual(ls, LinkSet{})) {
					l.LinkSets = append(l.LinkSets, ls)
				}
				if err != nil {
					return err
				}
				st = st.drop()
				continue
			}
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
			case "ERROR":
				l.Err = Error(string(t))
			case "eLinkResult":
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
