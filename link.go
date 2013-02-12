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
//                 SubjectType*,
// 				Category*,
//                 Attribute*,
//                 Provider,
//                 SubProvider?
// 			)>
//
// <!ELEMENT	IdUrlSet	(Id,(ObjUrl+|Info))>
//
// <!ELEMENT   FirstChars  (FirstChar*)>
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

// What does 'IMPLIED' mean on a boolean?

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

type Url struct {
	Url  string
	Lang string
}

func (u *Url) getLang(t xml.StartElement) error {
	for _, attr := range t.Attr {
		switch attr.Name.Local {
		case "LNG":
			switch attr.Value {
			case "DA", "DE", "EN", "EL", "ES", "FR", "IT", "IW", "JA", "NL", "NO", "RU", "SV", "ZH":
				u.Lang = attr.Value
			default:
				return fmt.Errorf("eutil: unknown language id: %q", attr.Value)
			}
		default:
			return fmt.Errorf("entrez: unknown attribute: %q", attr.Name.Local)
		}
	}
	return nil
}

type Provider struct {
	Name     string
	NameAbbr string
	Id       int
	Url      Url
	IconUrl  *Url
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
			switch t.Name.Local {
			case "Url":
				err = p.Url.getLang(t)
			case "IconUrl":
				p.IconUrl = &Url{}
				err = p.IconUrl.getLang(t)
			}
			if err != nil {
				return err
			}
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
				p.Url.Url = string(t)
			case "IconUrl":
				p.IconUrl.Url = string(t)
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
	Url         Url
	IconUrl     *Url
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
			switch t.Name.Local {
			case "Provider":
				err := ou.Provider.unmarshal(dec, st[len(st)-1:])
				if err != nil {
					return err
				}
				st = st.drop()
			case "Url":
				err = ou.Url.getLang(t)
			case "IconUrl":
				ou.IconUrl = &Url{}
				err = ou.IconUrl.getLang(t)
			}
			if err != nil {
				return err
			}
		case xml.CharData:
			if st.empty() {
				continue
			}
			switch name := st.peek(0); name {
			case "Url":
				ou.Url.Url = string(t)
			case "IconUrl":
				ou.IconUrl.Url = string(t)
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
